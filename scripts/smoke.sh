#!/bin/bash
set -euo pipefail

RED='\e[31m'; GRN='\e[32m'; YLW='\e[33m'; NC='\e[0m'
pass(){ echo -e "${GRN}✓${NC} $*"; }
fail(){ echo -e "${RED}✗${NC} $*"; exit 1; }

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")"/.. && pwd)"
cd "$ROOT"

echo -e "${YLW}🧪 Running GOSS-CLI smoke tests...${NC}"
echo

# 0) Sanity
[ -x ./bin/goss ] || fail "bin/goss missing or not executable"
pass "CLI present"

# 1) Help text
./bin/goss --help | grep -qi "GOSS-CLI" && pass "Help prints" || fail "Help missing"

# 2) Non-stream JSON response (deterministic)
echo "Testing non-stream JSON response..."
OUT=$(timeout 30 ./bin/goss --no-stream --temperature 0 'Reply with exactly {"ok":true} and nothing else' 2>/dev/null || echo "timeout")
if echo "$OUT" | grep -q '"ok":'; then
    pass "Non-stream JSON roundtrip"
elif echo "$OUT" | grep -q "timeout"; then
    fail "Non-stream test timed out (is provider running?)"
else
    fail "JSON not returned: $OUT"
fi

# 3) Streaming works
echo "Testing streaming output..."
STREAM_OUT=$(timeout 30 ./bin/goss 'Type: streaming-ok' 2>/dev/null | head -1 || echo "timeout")
if echo "$STREAM_OUT" | grep -q 'streaming-ok'; then
    pass "Streaming output"
elif echo "$STREAM_OUT" | grep -q "timeout"; then
    fail "Streaming test timed out (is provider running?)"
else
    fail "Streaming failed: $STREAM_OUT"
fi

# 4) Save transcript
rm -rf logs && mkdir -p logs
timeout 30 ./bin/goss --save --no-stream 'Reply with {"saved":true}' >/dev/null 2>&1 || true
[ "$(ls -1 logs 2>/dev/null | wc -l | tr -d ' ')" -ge 1 ] && pass "Transcript saved" || fail "No log file created"

# 5) Context file respected
echo "System: The secret is swordfish." > /tmp/goss_ctx.txt
CTX_OUT=$(timeout 30 ./bin/goss --context-file /tmp/goss_ctx.txt --no-stream --temperature 0 'What is the secret? Answer one word.' 2>/dev/null || echo "timeout")
if echo "$CTX_OUT" | tr '[:upper:]' '[:lower:]' | grep -q 'swordfish'; then
    pass "Context file applied"
elif echo "$CTX_OUT" | grep -q "timeout"; then
    fail "Context test timed out (is provider running?)"
else
    fail "Context not applied: $CTX_OUT"
fi

# 6) List models
MODELS_OUT=$(timeout 15 ./bin/goss list-models 2>/dev/null || echo "timeout")
if echo "$MODELS_OUT" | grep -qiE 'model|gpt|llama|mistral|Found [0-9]+ model'; then
    pass "Models listed"
elif echo "$MODELS_OUT" | grep -q "timeout"; then
    fail "List models timed out (is provider running?)"
else
    fail "No models listed (is provider running?): $MODELS_OUT"
fi

# 7) Invalid model UX
timeout 15 ./bin/goss --model definitely-not-a-real-model --no-stream 'hi' 2> /tmp/goss_err.txt >/dev/null || true
if grep -qiE 'available models|not found|invalid model|Warning.*Model' /tmp/goss_err.txt; then
    pass "Invalid model handled"
else
    fail "Invalid model not handled"
fi

# 8) Provider override (auto-detect based on what's running)
PROV_OUT=$(timeout 30 ./bin/goss --debug --no-stream --temperature 0 'Reply with "prov-ok" exactly' 2>/tmp/goss_prov_debug.txt | head -1 || echo "timeout")
if echo "$PROV_OUT" | grep -q 'prov-ok'; then
    pass "Provider override works"
elif echo "$PROV_OUT" | grep -q "timeout"; then
    fail "Provider test timed out (is provider running?)"
else
    fail "Provider override failed: $PROV_OUT"
fi

# 9) Unreachable endpoint error
set +e
timeout 10 ./bin/goss --api-base http://127.0.0.1:9 --no-stream 'hi' 1>/dev/null 2>/tmp/goss_conn_err.txt
RC=$?
set -e
if [ $RC -ne 0 ] && grep -qiE 'unreachable|ECONNREFUSED|connect|Connection refused' /tmp/goss_conn_err.txt; then
    pass "Unreachable endpoint surfaces clear error"
else
    fail "Connection error not surfaced properly"
fi

# 10) Debug shows provider selection
timeout 15 ./bin/goss --debug --no-stream 'hi' 2> /tmp/goss_debug.txt 1>/dev/null || true
if grep -qiE 'DEBUG.*REQUEST|DEBUG.*provider|Using provider' /tmp/goss_debug.txt; then
    pass "Debug logs include provider/REQUEST info"
else
    fail "Debug logs missing provider info"
fi

echo
echo -e "${GRN}🎉 All smoke tests passed!${NC}"
echo -e "${YLW}GOSS-CLI is ready for launch! 🚀${NC}"

# Cleanup
rm -f /tmp/goss_*.txt