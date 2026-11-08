#!/bin/bash
# End to end test of a joint transaction on a throwaway fiber coin chain.
# Two wallets each fund a small output, then one transaction spends both: the
# customer gets the coins back and the vendor gets coin hours. Each wallet signs
# only its own input, and the result must verify, broadcast and confirm.
# Needs go, curl and jq.

set -euo pipefail

COIN="${COIN:-skycoin}"
PORT="${PORT:-6431}"
while (exec 3<>"/dev/tcp/127.0.0.1/$PORT") 2>/dev/null; do PORT=$((PORT+1)); done
HOST="http://127.0.0.1:$PORT"
export RPC_ADDR="$HOST"

DATA_DIR=$(mktemp -d -t joint-txn-test.XXXXXX)
WALLET_DIR="$DATA_DIR/wallets"
FIBER_TOML="$DATA_DIR/fiber.toml"
NODE="$DATA_DIR/node.test"
SKY="$DATA_DIR/skycoin"

cleanup() {
    if [[ -n "${NODE_PID:-}" ]]; then
        kill -s SIGINT "$NODE_PID" 2>/dev/null || true
        wait "$NODE_PID" 2>/dev/null || true
    fi
    rm -rf "$DATA_DIR"
}
trap cleanup EXIT

fail() { echo "FAIL: $*" >&2; tail -n 20 "$DATA_DIR/node.log" >&2; exit 1; }
pass() { echo "PASS: $*"; }
# The node has no peers, so inject without broadcasting
inject() {
    curl -sf -X POST "$HOST/api/v1/injectTransaction" -H "Content-Type: application/json" \
        -d "$(jq -n --arg tx "$1" '{rawtx: $tx, no_broadcast: true}')" | jq -r .
}

head_seq() { curl -s "$HOST/api/v1/blockchain/metadata" | jq -r '.head.seq'; }
wait_block() {
    local want=$1
    for _ in $(seq 1 60); do
        [[ "$(head_seq)" -ge "$want" ]] && return 0
        sleep 1
    done
    fail "block $want was not created"
}

echo "=== build"
go build -o "$SKY" .
go test -c -tags testrunmain -o "$NODE" ./cmd/${COIN}/

echo "=== genesis and distribution addresses"
GENESIS_WALLET="$DATA_DIR/genesis.json"
"$SKY" cli addressGen -n 1 -m json -s "joint txn test $(date +%s%N)" > "$GENESIS_WALLET"
GENESIS_ADDRESS=$(jq -r '.entries[0].address' "$GENESIS_WALLET")
GENESIS_PUBKEY=$(jq -r '.entries[0].public_key' "$GENESIS_WALLET")
GENESIS_SECKEY=$(jq -r '.entries[0].secret_key' "$GENESIS_WALLET")

cat > "$FIBER_TOML" <<EOF
[node]
genesis_signature_str = ""
genesis_address_str = "${GENESIS_ADDRESS}"
blockchain_pubkey_str = "${GENESIS_PUBKEY}"
blockchain_seckey_str = "${GENESIS_SECKEY}"
genesis_timestamp = $(date +%s)
default_connections = []

[params]
max_coin_supply = 100000000
initial_unlocked_count = 25
unlock_address_rate = 5
unlock_time_interval = 31536000
EOF

FIBER_TOML="$FIBER_TOML" "$SKY" cli fiberAddressGen --num 100 --overwrite \
    --addrs-file "$DATA_DIR/addrs.txt" --seeds-file "$DATA_DIR/seeds.csv" >/dev/null

echo "=== start node on $HOST"
FIBER_TOML="$FIBER_TOML" GENESIS="$GENESIS_WALLET" "$NODE" \
    --disable-outgoing --disable-incoming --web-interface-port="$PORT" \
    --download-peerlist=false --launch-browser=false --data-dir="$DATA_DIR" \
    --enable-all-api-sets=true --wallet-dir="$WALLET_DIR" --disable-csrf \
    --block-publisher --test.run "^TestRunMain$" > "$DATA_DIR/node.log" 2>&1 &
NODE_PID=$!
for _ in $(seq 1 60); do
    curl -s "$HOST/api/v1/health" >/dev/null 2>&1 && break
    sleep 1
done
curl -s "$HOST/api/v1/health" >/dev/null || fail "node did not start"

"$SKY" cli distributeGenesis "$GENESIS_SECKEY" >/dev/null
wait_block 1
pass "genesis distributed"

seed_row() { sed -n "${1}p" "$DATA_DIR/seeds.csv" | cut -d'"' -f4; }
make_wallet() {
    local label=$1 seed=$2 id
    id=$(curl -s -X POST "$HOST/api/v1/wallet/create" --data-urlencode "seed=$seed" \
        -d "label=$label" -d "type=deterministic" -d "scan=2" | jq -r '.meta.filename')
    curl -s -X POST "$HOST/api/v1/wallet/newAddress" -d "id=$id" -d "num=1" >/dev/null
    echo "$id"
}
CUSTOMER=$(make_wallet customer "$(seed_row 1)")
VENDOR=$(make_wallet vendor "$(seed_row 2)")
addr() { "$SKY" cli listAddresses "$1" | jq -r ".addresses[$2]"; }
C0=$(addr "$CUSTOMER" 0); C1=$(addr "$CUSTOMER" 1)
V0=$(addr "$VENDOR" 0); V1=$(addr "$VENDOR" 1)
[[ "$C0" == "$(sed -n 1p "$DATA_DIR/addrs.txt" | cut -d'"' -f2)" ]] || fail "customer wallet does not own distribution address 1"
pass "wallets customer=$CUSTOMER vendor=$VENDOR"

echo "=== fund a carrier output and a float output"
sleep 5
RAW=$("$SKY" cli createRawTransactionV2 "$CUSTOMER" "$C1" 0.002)
inject "$RAW" >/dev/null
RAW=$("$SKY" cli createRawTransactionV2 "$VENDOR" "$V1" 0.001)
inject "$RAW" >/dev/null
wait_block 2

uxout() { "$SKY" cli addressOutputs "$1" | jq -r ".outputs.head_outputs[0].$2"; }
CUX=$(uxout "$C1" hash); CH=$(uxout "$C1" calculated_hours)
VUX=$(uxout "$V1" hash); VH=$(uxout "$V1" calculated_hours)
[[ -n "$CUX" && "$CUX" != null && -n "$VUX" && "$VUX" != null ]] || fail "carrier outputs not found"
pass "carrier $CUX ($CH h), float $VUX ($VH h)"

echo "=== create the joint transaction"
PAY=$(( CH / 2 ))
RAW=$("$SKY" cli createJointTransaction --in "$CUX" --in "$VUX" \
    --out "$C0,0.002,0" --out "$V0,0.001,$PAY")
"$SKY" cli decodeRawTransaction --wallet "$CUSTOMER" "$RAW"
[[ $(jq -r '.inputs | length' <<<"$("$SKY" cli decodeRawTransaction "$RAW")") == 2 ]] || fail "want 2 inputs"

if "$SKY" cli createJointTransaction --in "$CUX" --in "$VUX" --out "$C0,0.001,0" >/dev/null 2>&1; then
    fail "coins left over without a change address were accepted"
fi
pass "leftover coins need --change-address"

echo "=== customer signs"
OUT=$("$SKY" cli signTransaction "$CUSTOMER" "$RAW" 2>"$DATA_DIR/sign.err")
grep -q "signed 1 of 2 inputs" "$DATA_DIR/sign.err" || fail "customer: $(cat "$DATA_DIR/sign.err")"
RAW=$(jq -r '.encoded_transaction' <<<"$OUT")
if "$SKY" cli verifyTransaction "$RAW" >/dev/null 2>&1; then
    fail "a half signed transaction verified as spendable"
fi
# The web wallet reviews a half signed transaction this way
UNSIGNED_OK=$(curl -s -X POST "$HOST/api/v2/transaction/verify" -H "Content-Type: application/json" \
    -d "$(jq -n --arg tx "$RAW" '{encoded_transaction: $tx, unsigned: true}')" | jq -r '.error == null and .data.unsigned')
[[ "$UNSIGNED_OK" == true ]] || fail "verify as unsigned rejected a half signed transaction"
if "$SKY" cli signTransaction "$VENDOR" "$RAW" --sign-indexes 0 >/dev/null 2>&1; then
    fail "vendor signed the customer's input"
fi
pass "customer signed 1 of 2, not yet spendable"

echo "=== vendor signs"
OUT=$("$SKY" cli signTransaction "$VENDOR" "$RAW" 2>"$DATA_DIR/sign.err")
grep -q "signed 2 of 2 inputs" "$DATA_DIR/sign.err" || fail "vendor: $(cat "$DATA_DIR/sign.err")"
RAW=$(jq -r '.encoded_transaction' <<<"$OUT")
"$SKY" cli decodeRawTransaction --wallet "$VENDOR" "$RAW" | tee "$DATA_DIR/review.txt"
grep -q "ready to broadcast" "$DATA_DIR/review.txt" || fail "review does not say fully signed"
"$SKY" cli verifyTransaction "$RAW" | grep -q "spendable" || fail "verify"
pass "vendor signed 2 of 2, verified"

echo "=== inject"
SEQ=$(head_seq)
TXID=$(inject "$RAW")
wait_block $((SEQ + 1))
STATUS=$("$SKY" cli transaction "$TXID" | jq -r '.transaction.status.confirmed')
[[ "$STATUS" == true ]] || fail "transaction $TXID not confirmed"
GOT=$("$SKY" cli addressOutputs "$V0" | jq -r --arg tx "$TXID" '[.outputs.head_outputs[] | select(.src_tx == $tx)][0].hours')
[[ "$GOT" == "$PAY" ]] || fail "vendor output has $GOT hours, want $PAY"
pass "joint transaction $TXID confirmed, vendor got $PAY hours"

echo "ALL JOINT TRANSACTION TESTS PASSED"
