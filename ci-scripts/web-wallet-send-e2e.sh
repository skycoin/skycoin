#!/bin/bash
# Browser test of an ordinary send from a deterministic browser wallet on a
# throwaway fiber coin chain (joint-chain.sh). The spec loads a funded seed with
# one address, sends once signed in the browser and once as an unsigned
# transaction signed in the offline dialog, and checks each time that the wallet
# lists the address its change went to and still shows all its coins.
# Needs go, curl, jq, and npm ci in explorer. CHROME_BIN picks the browser.

set -euo pipefail
REPO_ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
cd "$REPO_ROOT"
JOINT_PEERED=1
source ci-scripts/joint-chain.sh
joint_chain_setup

SENDER_SEED="${SENDER_SEED:-legal winner thank year wave sausage worth useful legal winner thank yellow}"
echo "=== fund the sending wallet"
DIST=$(make_wallet dist3 "$(seed_row 3)")
SENDER=$(make_wallet sender "$SENDER_SEED")
curl -s -X POST "$HOST/api/v1/wallet/newAddress" -d "id=$SENDER" -d "num=1" >/dev/null
S0=$(addr "$SENDER" 0); S1=$(addr "$SENDER" 1); S2=$(addr "$SENDER" 2)
inject "$("$SKY" cli createRawTransactionV2 "$DIST" "$S0" 10)" >/dev/null
wait_block 3
[[ $(uxout "$S0" coins) == "10.000000" ]] || fail "the sending wallet was not funded"
pass "sender $S0 holds 10 coins"

# The wallet refuses a node without peers, and the peer relays its broadcast
joint_chain_add_peer

WEB_PORT=$(free_port $((PORT + 1)))
WEB="http://127.0.0.1:$WEB_PORT"
echo "=== start the web wallet on $WEB"
"$SKY" web --node-url "$PEER_HOST" --port "$WEB_PORT" --btc-electrum-url none > "$DATA_DIR/web.log" 2>&1 &
EXTRA_PIDS+=($!)
for _ in $(seq 1 60); do
    curl -sf -o /dev/null "$WEB/" && break
    sleep 1
done
curl -sf -o /dev/null "$WEB/" || { tail -20 "$DATA_DIR/web.log" >&2; fail "web wallet did not start"; }

export SEND_SEED="$SENDER_SEED" SEND_ADDRS="$S0,$S1,$S2" SEND_TO="$V1" SEND_NODE="$PEER_HOST"

echo "=== run the browser spec"
(cd explorer && npx wdio run "$REPO_ROOT/src/skycoin-web/e2e/wdio-send.conf.ts" --baseUrl "$WEB")

GOT=$(uxout "$S2" coins)
[[ "$GOT" == "5.000000" ]] || fail "the last change address holds $GOT coins, want 5.000000"
pass "both sends confirmed, the change of the second is $GOT coins on $S2"

echo "ALL WEB WALLET SEND TESTS PASSED"
