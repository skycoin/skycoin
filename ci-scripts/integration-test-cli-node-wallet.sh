#!/bin/bash
# End to end test of the CLI spending from a wallet file that a running node
# holds, on a throwaway fiber coin chain. The node keeps its own copy of the
# wallet and writes it over the file, so the address the CLI makes for the
# change, or for a sweep, has to be added by the node. The test checks the node
# lists the address and counts its coins, and that the address is still in the
# file after the node saves the wallet. Needs go, curl and jq.

set -euo pipefail
cd "$(dirname "${BASH_SOURCE[0]}")/.."
source ci-scripts/joint-chain.sh
joint_chain_setup

node_addrs() { curl -s "$HOST/api/v1/wallet?id=$1" | jq -r '.entries[].address'; }
node_coins() { curl -s "$HOST/api/v1/wallet/balance?id=$1" | jq -r '.confirmed.coins'; }
file_addrs() { jq -r '.entries[].address' "$WALLET_DIR/$1"; }
# A label change makes the node write its copy of the wallet over the file
node_save() {
    curl -sf -X POST "$HOST/api/v1/wallet/update" -d "id=$1" -d "label=saved by the node" >/dev/null \
        || fail "the node did not save wallet $1"
}
# check_kept <wallet> <address> <droplets> <what>
check_kept() {
    node_addrs "$1" | grep -qx "$2" || fail "the node's wallet does not list the $4 address $2"
    [[ $(node_coins "$1") == "$3" ]] || fail "the node counts $(node_coins "$1") droplets in the wallet, want $3"
    node_save "$1"
    file_addrs "$1" | grep -qx "$2" || fail "the $4 address $2 left the wallet file when the node saved it"
    pass "the node lists the $4 address $2, counts $3 droplets, and kept it in the file"
}

echo "=== send from the customer wallet file in the node's wallet directory"
[[ $(node_addrs "$CUSTOMER" | wc -l) == 2 ]] || fail "the customer wallet should start with 2 addresses"
RAW=$("$SKY" cli createRawTransaction "$WALLET_DIR/$CUSTOMER" "$V1" 0.001)
CHANGE=$("$SKY" cli decodeRawTransaction "$RAW" | jq -r --arg to "$V1" '[.outputs[].dst | select(. != $to)][0]')
[[ -n "$CHANGE" && "$CHANGE" != null && "$CHANGE" != "$C0" && "$CHANGE" != "$C1" ]] || fail "the change went to $CHANGE, want a new address"
inject "$RAW" >/dev/null
wait_block 3
check_kept "$CUSTOMER" "$CHANGE" 1000 change

echo "=== sweep the vendor wallet file in the node's wallet directory"
OUT=$("$SKY" cli sweep --raw --json "$WALLET_DIR/$VENDOR")
DEST=$(jq -r .address <<<"$OUT")
[[ "$DEST" != "$V0" && "$DEST" != "$V1" ]] || fail "the sweep went to $DEST, want a new address"
inject "$(jq -r .rawtx <<<"$OUT")" >/dev/null
wait_block 4
check_kept "$VENDOR" "$DEST" 2000 sweep

echo "ALL CLI NODE WALLET TESTS PASSED"
