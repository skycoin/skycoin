#!/bin/bash
# Browser test of the web wallet's joint transaction page on a throwaway fiber
# coin chain. The chain funds a customer carrier and a vendor float output
# (joint-chain.sh), then the spec loads both seeds as browser wallets, creates
# the transaction as the customer, and signs it in the browser as each party.
# The page broadcasts through a second node peered with the block publisher, and
# the script checks the confirmed outputs. Needs go, curl, jq, and npm ci in
# explorer, whose WebdriverIO install runs the spec. CHROME_BIN picks the browser.

set -euo pipefail
REPO_ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
cd "$REPO_ROOT"
JOINT_PEERED=1
source ci-scripts/joint-chain.sh
joint_chain_setup
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

PAY=$(( CH / 2 ))
export JOINT_CUSTOMER_SEED="$CUSTOMER_SEED" JOINT_VENDOR_SEED="$VENDOR_SEED"
export JOINT_CUSTOMER_ADDR="$C0" JOINT_VENDOR_ADDR="$V0"
export JOINT_CUSTOMER_UX="$CUX" JOINT_VENDOR_UX="$VUX"
export JOINT_CUSTOMER_BACK="$C1" JOINT_VENDOR_PAID="$V1" JOINT_PAY="$PAY"
export JOINT_OUT_FILE="$DATA_DIR/signed.hex"

echo "=== run the browser spec"
(cd explorer && npx wdio run "$REPO_ROOT/src/skycoin-web/e2e/wdio-joint.conf.ts" --baseUrl "$WEB")
RAW=$(cat "$JOINT_OUT_FILE")
[[ -n "$RAW" ]] || fail "the spec wrote no transaction"
TXID=$("$SKY" cli decodeRawTransaction "$RAW" | jq -r .txid)

echo "=== wait for the page's broadcast to confirm"
for _ in $(seq 1 60); do
    [[ $("$SKY" cli transaction "$TXID" 2>/dev/null | jq -r '.transaction.status.confirmed') == true ]] && break
    sleep 1
done
[[ $("$SKY" cli transaction "$TXID" | jq -r '.transaction.status.confirmed') == true ]] || fail "transaction $TXID not confirmed"
out_field() {
    "$SKY" cli addressOutputs "$1" | jq -r --arg tx "$TXID" "[.outputs.head_outputs[] | select(.src_tx == \$tx)][0].$2"
}
GOT=$(out_field "$V1" hours)
[[ "$GOT" == "$PAY" ]] || fail "vendor output has $GOT hours, want $PAY"
BACK=$(out_field "$C1" coins)
[[ "$BACK" == "0.002000" ]] || fail "customer output has $BACK coins, want 0.002000"
pass "browser signed joint transaction $TXID confirmed, vendor got $PAY hours, customer got $BACK back"

echo "ALL WEB WALLET JOINT TRANSACTION TESTS PASSED"
