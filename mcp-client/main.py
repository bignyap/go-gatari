import os
import uuid
import json
from typing import Any, Dict, Optional

import httpx
from fastapi import FastAPI, WebSocket, WebSocketDisconnect
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel, Field
from dotenv import load_dotenv

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

# ---- JSON-RPC helpers to talk to MCP server ----

def _jsonrpc_payload(method: str, params: Optional[Dict[str, Any]] = None):
    return {
        "jsonrpc": "2.0",
        "id": str(uuid.uuid4()),
        "method": method,
        "params": params or {},
    }

async def mcp_rpc(method: str, params: Optional[Dict[str, Any]] = None):
    async with httpx.AsyncClient(timeout=60) as client:
        r = await client.post(MCP_SERVER_URL, json=_jsonrpc_payload(method, params))
        r.raise_for_status()
        data = r.json()
        if "error" in data:
            raise RuntimeError(f"MCP error: {data['error']}")
        return data["result"]

async def mcp_list_tools():
    res = await mcp_rpc("tools/list")
    # Expected shape: { "tools": [ { "name": str, "description": str, "parameters": {...} }, ... ] }
    return res.get("tools", [])

async def mcp_call_tool(name: str, arguments: Dict[str, Any]):
    return await mcp_rpc("tools/call", {"name": name, "arguments": arguments})

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
    """
    Ask the LLM to map the user request to a tool + args.
    Returns a JSON dict matching INTENT_SCHEMA.
    """
    # We give the LLM the list of tool names and their JSON Schemas (if available).
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

    # Use JSON mode to force valid JSON
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

    # Final sanity checks on keys
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
    # type: 'user' messages from the frontend
    type: str = Field(..., description="Currently only 'user'")
    text: str

# Outgoing message envelope
def ws_msg(kind: str, text: str, payload: Optional[dict] = None) -> str:
    """
    kind: 'bot', 'confirm', 'result', 'error'
    """
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

            # If we are waiting for a yes/no on a pending action:
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

            # Otherwise: parse new intent with LLM
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

            # Show confirmation card
            pending_action = {"tool": tool, "args": args}
            confirm_text = f"{summary}\n\n{confirmation}\n\nTool: {tool}\nArgs: {json.dumps(args, ensure_ascii=False)}\n\nType 'yes' to proceed or 'no' to cancel."
            await ws.send_text(ws_msg("confirm", confirm_text, payload={"tool": tool, "args": args}))

    except WebSocketDisconnect:
        # Client left
        return

# --- Optional: simple health endpoint ---
@app.get("/health")
async def health():
    return {"ok": True}
