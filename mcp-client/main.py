import os
import uuid
import json
from typing import Any, Dict, Optional

import httpx
from fastapi import FastAPI, WebSocket, WebSocketDisconnect
from fastapi.middleware.cors import CORSMiddleware
import uvicorn
from pydantic import BaseModel, Field
from dotenv import load_dotenv
from fastmcp import Client

# LangChain imports
from langchain_openai import ChatOpenAI
from langchain_google_genai import ChatGoogleGenerativeAI
from langchain.prompts import ChatPromptTemplate
from langchain_core.output_parsers import JsonOutputParser

load_dotenv()

MCP_SERVER_URL = os.getenv("MCP_SERVER_URL", "http://localhost:8084/mcp")
OPENAI_API_KEY = os.getenv("OPENAI_API_KEY", "")
OPENAI_MODEL = os.getenv("OPENAI_MODEL", "gpt-4o-mini")
GEMINI_API_KEY = os.getenv("GEMINI_API_KEY", "")
GEMINI_MODEL = os.getenv("GEMINI_MODEL", "gemini-1.5-flash")
LLM_PROVIDER = os.getenv("LLM_PROVIDER", "openai")  # "openai" or "gemini"
PORT = int(os.getenv("BACKEND_PORT", "9000"))
ALLOWED_ORIGINS = os.getenv("ALLOWED_ORIGINS", "http://localhost:5173").split(",")

app = FastAPI(title="GATARI MCP Chat Backend")

# ---- CORS ----
app.add_middleware(
    CORSMiddleware,
    allow_origins=[o.strip() for o in ALLOWED_ORIGINS],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# ---- MCP client ----
_mcp_client: Optional[Client] = None

async def get_mcp_client() -> Client:
    global _mcp_client
    if _mcp_client is None:
        _mcp_client = Client(MCP_SERVER_URL)
        await _mcp_client.__aenter__()  # establish connection
    return _mcp_client

async def mcp_list_tools():
    client = await get_mcp_client()
    tools = await client.list_tools()
    return [
        {
            "name": t.name,
            "description": getattr(t, "description", ""),
            "parameters": getattr(t, "parameters", {}),
        }
        for t in tools
    ]

async def mcp_call_tool(name: str, arguments: Dict[str, Any]):
    client = await get_mcp_client()
    return await client.call_tool(name, arguments)

# ---- LangChain LLM Factory ----
def get_llm(provider: str = LLM_PROVIDER):
    if provider == "openai":
        if not OPENAI_API_KEY:
            raise RuntimeError("OPENAI_API_KEY not set")
        return ChatOpenAI(
            model=OPENAI_MODEL,
            api_key=OPENAI_API_KEY,
            temperature=0.1,
        )
    elif provider == "gemini":
        if not GEMINI_API_KEY:
            raise RuntimeError("GEMINI_API_KEY not set")
        return ChatGoogleGenerativeAI(
            model=GEMINI_MODEL,
            google_api_key=GEMINI_API_KEY,
            temperature=0.1,
        )
    else:
        raise ValueError(f"Unsupported LLM provider: {provider}")

# ---- Prompt & Parser ----
INTENT_PROMPT = ChatPromptTemplate.from_messages([
    ("system",
     "You convert natural language requests into MCP tool invocations.\n"
     "Return ONLY a strict JSON object with fields: tool, args, summary, confirmation.\n"
     "If no suitable tool exists, set tool=null and args={{}}. Do not invent tools.\n"
     "Use the provided parameter schemas when forming args."),
    ("user", "{user_input}")
])
parser = JsonOutputParser()

async def llm_plan(user_message: str, tools: list[dict], provider: str = LLM_PROVIDER) -> dict:
    llm = get_llm(provider)

    tool_catalog = [
        {
            "name": t.get("name"),
            "description": t.get("description", ""),
            "parameters": t.get("parameters", {}),
        }
        for t in tools
    ]

    user_payload = json.dumps({
        "available_tools": tool_catalog,
        "user_request": user_message,
    }, ensure_ascii=False)

    chain = INTENT_PROMPT | llm | parser

    try:
        return await chain.ainvoke({"user_input": user_payload})
    except Exception as e:
        return {
            "tool": None,
            "args": {},
            "summary": f"Failed to parse plan: {e}",
            "confirmation": "I couldn't understand the request."
        }

# ---- WebSocket chat protocol ----
class WSIncoming(BaseModel):
    type: str = Field(..., description="Currently only 'user'")
    text: str

def ws_msg(kind: str, text: str, payload: Optional[dict] = None) -> str:
    out = {"type": kind, "text": text}
    if payload is not None:
        out["payload"] = payload
    return json.dumps(out, ensure_ascii=False)

@app.websocket("/ws")
async def ws_chat(ws: WebSocket):
    await ws.accept()
    pending_action: Optional[dict] = None
    tools = await mcp_list_tools()

    try:
        await ws.send_text(ws_msg("bot", "Hi! Tell me what to do. I'll propose an action and ask for confirmation."))
        while True:
            raw = await ws.receive_text()
            try:
                incoming = WSIncoming.model_validate_json(raw)
            except Exception:
                await ws.send_text(ws_msg("error", "Invalid message. Send JSON: { type: 'user', text: '...' }"))
                continue

            if incoming.type != "user":
                await ws.send_text(ws_msg("error", "Unsupported message type."))
                continue

            user_text = incoming.text.strip()

            # If awaiting confirmation
            if pending_action:
                low = user_text.lower()
                if low in ("y", "yes", "confirm", "ok", "okay", "do it"):
                    await ws.send_text(ws_msg("bot", "Executing…"))
                    try:
                        result = await mcp_call_tool(
                            pending_action["tool"], pending_action["args"]
                        )

                        # ✅ Convert to JSON-safe dict
                        if hasattr(result, "dict"):
                            result_payload = result.dict()
                        elif hasattr(result, "json"):
                            result_payload = json.loads(result.json())
                        else:
                            result_payload = str(result)

                        await ws.send_text(ws_msg("result", "✅ Done.", payload=result_payload))
                    except Exception as e:
                        await ws.send_text(ws_msg("error", f"❌ Execution failed: {e}"))
                    finally:
                        pending_action = None
                elif low in ("n", "no", "cancel", "stop"):
                    await ws.send_text(ws_msg("bot", "❌ Cancelled. What next?"))
                    pending_action = None
                else:
                    await ws.send_text(ws_msg("bot", "Please reply 'yes' or 'no'."))
                continue

            # Otherwise create a plan
            await ws.send_text(ws_msg("bot", "Thinking…"))
            try:
                plan = await llm_plan(user_text, tools, provider=LLM_PROVIDER)
            except Exception as e:
                await ws.send_text(ws_msg("error", f"LLM error: {e}"))
                continue

            print("Plan:", plan)

            tool = plan.get("tool")
            args = plan.get("args") or {}
            summary = plan.get("summary") or ""
            confirmation = plan.get("confirmation") or ""

            if not tool:
                await ws.send_text(ws_msg("bot", "I couldn't map that to a known action. Try rephrasing."))
                continue

            pending_action = {"tool": tool, "args": args}
            confirm_text = (
                f"{summary}\n\n{confirmation}\n\n"
                f"Tool: {tool}\nArgs: {json.dumps(args, ensure_ascii=False)}\n\n"
                "Type 'yes' to proceed or 'no' to cancel."
            )
            await ws.send_text(ws_msg("confirm", confirm_text, payload={"tool": tool, "args": args}))

    except WebSocketDisconnect:
        return

@app.get("/health")
async def health():
    return {"ok": True}

if __name__ == "__main__":
    uvicorn.run(
        "main:app",
        host="0.0.0.0",
        port=PORT,
        reload=True
    )