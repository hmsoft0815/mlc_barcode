#!/usr/bin/env bash
# Network transports of the MCP server: Streamable HTTP (/mcp) and legacy SSE (/sse).
# Usage: http-check.sh <mcp-tester> <server-binary> <log-dir>
set -u
TESTER=$1 SERVER=$2 LOGS=$3

PORT=$(python3 -c 'import socket; s=socket.socket(); s.bind(("127.0.0.1",0)); print(s.getsockname()[1])')
"$SERVER" -addr "127.0.0.1:$PORT" > /dev/null 2>&1 &
PID=$!
trap 'kill $PID 2>/dev/null' EXIT
sleep 1

rc=0
"$TESTER" http-check -u "http://127.0.0.1:$PORT/mcp" || { echo "✗ http-check"; rc=1; }

"$TESTER" inspect -t streamable-http -u "http://127.0.0.1:$PORT/mcp" --min-score 100 > "$LOGS/mcptest-http.log" 2>&1 \
  || { echo "✗ inspect streamable-http"; rc=1; }
grep -E "Protocol|Score" "$LOGS/mcptest-http.log"

# SSE is the 2024-11-05 transport and negotiates 2025-11-25 at most;
# mcp-tester takes 10 points off per revision behind, hence 90.
"$TESTER" inspect -t sse -u "http://127.0.0.1:$PORT/sse" --min-score 90 > "$LOGS/mcptest-sse.log" 2>&1 \
  || { echo "✗ inspect sse"; rc=1; }
grep -E "Protocol|Score" "$LOGS/mcptest-sse.log"

exit $rc
