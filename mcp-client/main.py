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
from fastmcp import Client  # ✅ no more transport import

# OpenAI (official python client)
from openai import OpenAI

load_dotenv()

MCP_SERVER_URL = os.getenv("MCP_SERVER_URL", "http://localhost:8084/mcp")
OPENAI_API_KEY = os.getenv("OPENAI_API_KEY", "")
OPENAI_MODEL = os.getenv("OPENAI_MODEL", "gpt-4o-mini")
PORT = int(os.getenv("BACKEND_PORT", "9000"))
ALLOWED_ORIGINS = os.getenv("ALLOWED_ORIGINS", "http://localhost:5173").split(",")

app = FastAPI(title="GATARI MCP Chat Backend")

# CORS (so the React app can talk to us)
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
    return await client.list_tools()

async def mcp_call_tool(name: str, arguments: Dict[str, Any]):
    client = await get_mcp_client()
    return await client.call_tool(name, arguments)

# ---- OpenAI intent parsing ----

if not OPENAI_API_KEY:
    raise RuntimeError("OPENAI_API_KEY not set. Put it in backend/.env")

oai = OpenAI(api_key=OPENAI_API_KEY)

INTENT_SCHEMA = {
    "type": "object",
    "properties": {
        "tool": {"type": ["string", "null"]},
        "args": {"type": "object"},
        "summary": {"type": "string"},
        "confirmation": {"type": "string"},
    },
    "required": ["tool", "args", "summary", "confirmation"],
    "additionalProperties": False,
}

async def llm_plan(user_message: str, tools: list[dict]) -> dict:
    tool_catalog = []
    for t in tools:
        tool_catalog.append({
            "name": t.get("name"),
            "description": t.get("description", ""),
            "parameters": t.get("parameters", {}),
        })

    system = (
        "You convert natural language requests into MCP tool invocations.\n"
        "Return ONLY a strict JSON object with fields: tool, args, summary, confirmation.\n"
        "If no suitable tool exists, set tool=null and args={}. Do not invent tools.\n"
        "Use the provided parameter schemas when forming args."
    )

    user = json.dumps({
        "available_tools": tool_catalog,
        "user_request": user_message,
    }, ensure_ascii=False)

    resp = oai.chat.completions.create(
        model=OPENAI_MODEL,
        messages=[
            {"role": "system", "content": system},
            {"role": "user", "content": user},
        ],
        response_format={"type": "json_object"},
        temperature=0.1,
    )

    content = resp.choices[0].message.content or "{}"
    try:
        parsed = json.loads(content)
    except Exception:
        parsed = {"tool": None, "args": {}, "summary": "Failed to parse plan.", "confirmation": "I couldn't understand the request."}

    for k in ["tool", "args", "summary", "confirmation"]:
        if k not in parsed:
            if k == "tool":
                parsed[k] = None
            elif k == "args":
                parsed[k] = {}
            else:
                parsed[k] = ""
    return parsed

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
        await ws.send_text(ws_msg("bot", "Hi! Tell me what to do. I’ll propose an action and ask for confirmation."))
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

            if pending_action:
                low = user_text.lower()
                if low in ("y", "yes", "confirm", "ok", "okay", "do it"):
                    await ws.send_text(ws_msg("bot", "Executing…"))
                    try:
                        result = await mcp_call_tool(
                            pending_action["tool"], pending_action["args"]
                        )
                        await ws.send_text(ws_msg("result", "✅ Done.", payload=result))
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

            await ws.send_text(ws_msg("bot", "Thinking…"))
            try:
                plan = await llm_plan(user_text, tools)
            except Exception as e:
                await ws.send_text(ws_msg("error", f"LLM error: {e}"))
                continue

            tool = plan.get("tool")
            args = plan.get("args") or {}
            summary = plan.get("summary") or ""
            confirmation = plan.get("confirmation") or ""

            if not tool:
                await ws.send_text(ws_msg("bot", "I couldn't map that to a known action. Try rephrasing."))
                continue

            pending_action = {"tool": tool, "args": args}
            confirm_text = f"{summary}\n\n{confirmation}\n\nTool: {tool}\nArgs: {json.dumps(args, ensure_ascii=False)}\n\nType 'yes' to proceed or 'no' to cancel."
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