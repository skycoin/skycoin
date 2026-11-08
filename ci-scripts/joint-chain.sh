# Sourced by the joint transaction tests. Starts a throwaway fiber coin chain
# whose genesis coins are distributed, then funds a customer carrier output of
# 0.002 coins and a vendor float output of 0.001 coins, each held by a
# deterministic wallet made from a known seed so a browser can load it too.
# Needs go, curl and jq. Call joint_chain_setup, then use what it sets:
#   HOST SKY DATA_DIR CUSTOMER VENDOR CUSTOMER_SEED VENDOR_SEED
#   C0 C1 V0 V1 (first and second addresses), CUX CH VUX VH (outputs, hours)

COIN="${COIN:-skycoin}"
CUSTOMER_SEED="${CUSTOMER_SEED:-abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about}"
VENDOR_SEED="${VENDOR_SEED:-zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo wrong}"

free_port() {
    local port=$1
    while (exec 3<>"/dev/tcp/127.0.0.1/$port") 2>/dev/null; do port=$((port+1)); done
    echo "$port"
}

PORT=$(free_port "${PORT:-6431}")
HOST="http://127.0.0.1:$PORT"
export RPC_ADDR="$HOST"

DATA_DIR=$(mktemp -d -t joint-txn-test.XXXXXX)
WALLET_DIR="$DATA_DIR/wallets"
FIBER_TOML="$DATA_DIR/fiber.toml"
NODE="$DATA_DIR/node.test"
SKY="$DATA_DIR/skycoin"
EXTRA_PIDS=()

joint_chain_cleanup() {
    local pid
    for pid in "${EXTRA_PIDS[@]}"; do
        kill "$pid" 2>/dev/null || true
        wait "$pid" 2>/dev/null || true
    done
    if [[ -n "${NODE_PID:-}" ]]; then
        kill -s SIGINT "$NODE_PID" 2>/dev/null || true
        wait "$NODE_PID" 2>/dev/null || true
    fi
    rm -rf "$DATA_DIR"
}
trap joint_chain_cleanup EXIT

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

make_wallet() {
    local label=$1 seed=$2 id
    id=$(curl -s -X POST "$HOST/api/v1/wallet/create" --data-urlencode "seed=$seed" \
        -d "label=$label" -d "type=deterministic" -d "scan=2" | jq -r '.meta.filename')
    curl -s -X POST "$HOST/api/v1/wallet/newAddress" -d "id=$id" -d "num=1" >/dev/null
    echo "$id"
}
addr() { "$SKY" cli listAddresses "$1" | jq -r ".addresses[$2]"; }
uxout() { "$SKY" cli addressOutputs "$1" | jq -r ".outputs.head_outputs[0].$2"; }

joint_chain_setup() {
    echo "=== build"
    go build -o "$SKY" .
    go test -c -tags testrunmain -o "$NODE" ./cmd/${COIN}/

    echo "=== genesis and distribution addresses"
    local genesis="$DATA_DIR/genesis.json" gaddr gpub gsec
    "$SKY" cli addressGen -n 1 -m json -s "joint txn test $(date +%s%N)" > "$genesis"
    gaddr=$(jq -r '.entries[0].address' "$genesis")
    gpub=$(jq -r '.entries[0].public_key' "$genesis")
    gsec=$(jq -r '.entries[0].secret_key' "$genesis")

    cat > "$FIBER_TOML" <<EOF
[node]
genesis_signature_str = ""
genesis_address_str = "${gaddr}"
blockchain_pubkey_str = "${gpub}"
blockchain_seckey_str = "${gsec}"
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

    local net=(--disable-outgoing --disable-incoming)
    echo "=== start node on $HOST"
    FIBER_TOML="$FIBER_TOML" GENESIS="$genesis" "$NODE" \
        "${net[@]}" --web-interface-port="$PORT" \
        --download-peerlist=false --launch-browser=false --data-dir="$DATA_DIR" \
        --enable-all-api-sets=true --wallet-dir="$WALLET_DIR" --disable-csrf \
        --block-publisher --test.run "^TestRunMain$" > "$DATA_DIR/node.log" 2>&1 &
    NODE_PID=$!
    for _ in $(seq 1 60); do
        curl -s "$HOST/api/v1/health" >/dev/null 2>&1 && break
        sleep 1
    done
    curl -s "$HOST/api/v1/health" >/dev/null || fail "node did not start"

    "$SKY" cli distributeGenesis "$gsec" >/dev/null
    wait_block 1
    pass "genesis distributed"

    local d1 d2 raw
    seed_row() { sed -n "${1}p" "$DATA_DIR/seeds.csv" | cut -d'"' -f4; }
    d1=$(make_wallet dist1 "$(seed_row 1)")
    d2=$(make_wallet dist2 "$(seed_row 2)")
    CUSTOMER=$(make_wallet customer "$CUSTOMER_SEED")
    VENDOR=$(make_wallet vendor "$VENDOR_SEED")
    C0=$(addr "$CUSTOMER" 0); C1=$(addr "$CUSTOMER" 1)
    V0=$(addr "$VENDOR" 0); V1=$(addr "$VENDOR" 1)
    pass "wallets customer=$CUSTOMER ($C0) vendor=$VENDOR ($V0)"

    echo "=== fund a carrier output and a float output"
    sleep 5
    raw=$("$SKY" cli createRawTransactionV2 "$d1" "$C0" 0.002)
    inject "$raw" >/dev/null
    raw=$("$SKY" cli createRawTransactionV2 "$d2" "$V0" 0.001)
    inject "$raw" >/dev/null
    wait_block 2

    CUX=$(uxout "$C0" hash); CH=$(uxout "$C0" calculated_hours)
    VUX=$(uxout "$V0" hash); VH=$(uxout "$V0" calculated_hours)
    [[ -n "$CUX" && "$CUX" != null && -n "$VUX" && "$VUX" != null ]] || fail "carrier outputs not found"
    pass "carrier $CUX ($CH h), float $VUX ($VH h)"
}

