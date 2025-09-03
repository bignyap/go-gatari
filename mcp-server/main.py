import httpx
import os
import yaml
from pathlib import Path
from prance import ResolvingParser
from fastmcp import FastMCP
from fastmcp.server.openapi import RouteMap, MCPType
from dotenv import load_dotenv

import asyncio
import sys

# Load environment variables from .env file
load_dotenv()

if sys.platform.startswith("win"):
    asyncio.set_event_loop_policy(asyncio.WindowsSelectorEventLoopPolicy())

local_spec_path_str = os.getenv("OPENAPI_SPEC_PATH", "/app/_apidoc/go-admin/swagger.yaml")
LOCAL_SPEC_PATH = Path(local_spec_path_str)
BASE_URL = os.getenv("BASE_URL", "http://localhost:8080")
PORT = int(os.getenv("PORT", 8000))

def _resolve_case_insensitive(base_dir: Path, rel: str) -> Path:
    """
    Resolve rel against base_dir, case-insensitively.
    If rel has N parts, we search base_dir with '**' up to that depth.
    """
    target = base_dir / rel
    print("Target:", target)
    if target.exists():
        return target

    parts_len = len(Path(rel).parts)
    # Build a glob like '*/*/...'
    pattern = "*/" * (parts_len - 1) + "*"
    for f in base_dir.glob(pattern):
        if f.is_file() and f.name.lower() == Path(rel).name.lower():
            return f

    raise FileNotFoundError(f"Could not resolve {rel} in {base_dir}")

def _rewrite_local_refs_to_absolute(obj, base_dir: Path):
    """
    Recursively rewrite any local $ref like './paths/file.yaml#/X'
    to an absolute file path so prance resolves everything locally.
    """
    if isinstance(obj, dict):
        for k, v in list(obj.items()):
            if k == "$ref" and isinstance(v, str) and v.startswith("./"):
                if "#" in v:
                    rel, frag = v.split("#", 1)
                    abs_path = _resolve_case_insensitive(base_dir, rel).resolve()
                    obj[k] = f"{abs_path.as_posix()}#{frag}"
                else:
                    abs_path = _resolve_case_insensitive(base_dir, v).resolve()
                    obj[k] = abs_path.as_posix()
            else:
                _rewrite_local_refs_to_absolute(v, base_dir)
    elif isinstance(obj, list):
        for item in obj:
            _rewrite_local_refs_to_absolute(item, base_dir)


def _clean_invalid_schema_fields(obj):
    """
    Fix common spec issues that trigger validation errors:
      - remove 'schema: null' entries (seen under parameters).
    """
    if isinstance(obj, dict):
        if "schema" in obj and obj["schema"] is None:
            del obj["schema"]
        for _, v in list(obj.items()):
            _clean_invalid_schema_fields(v)
    elif isinstance(obj, list):
        for item in obj:
            _clean_invalid_schema_fields(item)


def load_and_resolve_openapi() -> dict:
    """
    Load the local swagger.yaml, rewrite local $refs to absolute file paths,
    clean known invalid fields, and resolve with prance (validation disabled).
    """
    spec_path = LOCAL_SPEC_PATH.resolve()
    if not spec_path.exists():
        raise FileNotFoundError(f"Swagger file not found: {spec_path}")

    raw = spec_path.read_text(encoding="utf-8")
    spec = yaml.safe_load(raw)

    _rewrite_local_refs_to_absolute(spec, spec_path.parent)
    _clean_invalid_schema_fields(spec)

    cleaned_path = spec_path.parent / "__swagger_cleaned.yaml"
    cleaned_path.write_text(yaml.safe_dump(spec, sort_keys=False), encoding="utf-8")

    parser = ResolvingParser(str(cleaned_path), validate=False)
    return parser.specification

# Expose the 'mcp' object directly for the 'fastmcp run' command to find.
openapi_spec = load_and_resolve_openapi()
client = httpx.AsyncClient(base_url=BASE_URL)

mcp = FastMCP.from_openapi(
    openapi_spec=openapi_spec,
    name="GATARI API Server",
    client=client,
    route_maps=[
        RouteMap(
            methods=["POST"],
            pattern=r".*",
            mcp_type=MCPType.TOOL,
            mcp_tags={"write-operation", "api-mutation"},
        ),
        RouteMap(
            methods=["GET"],
            pattern=r".*\{.*\}.*",
            mcp_type=MCPType.RESOURCE_TEMPLATE,
            mcp_tags={"detail-view", "parameterized"},
        ),
        RouteMap(
            methods=["GET"],
            pattern=r".*",
            mcp_type=MCPType.RESOURCE,
            mcp_tags={"list-data", "collection"},
        ),
    ],
)

if __name__ == "__main__":
    mcp.run(transport="http", port=PORT, host="0.0.0.0")