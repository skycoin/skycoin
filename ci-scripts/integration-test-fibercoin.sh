#!/bin/bash
# Integration test for fibercoin creation lifecycle:
#   1. Generate genesis wallet
#   2. Create test fiber.toml
#   3. Generate distribution addresses (updates fiber.toml)
#   4. Start daemon with FIBER_TOML and GENESIS env vars
#   5. Verify genesis block was created
#   6. Distribute genesis coins
#   7. Verify distribution
#   8. Test BIP44 chain selection
#   9. Cleanup

set -euxo pipefail

SCRIPT=$(basename "${BASH_SOURCE[0]}")
COIN="${COIN:-skycoin}"

# Find unused port
PORT="1024"
while $(lsof -Pi :$PORT -sTCP:LISTEN -t >/dev/null 2>&1) ; do
    PORT=$((PORT+1))
done

HOST="http://127.0.0.1:$PORT"
RPC_ADDR="http://127.0.0.1:$PORT"

# Create temp directory for all test artifacts
DATA_DIR=$(mktemp -d -t fibercoin-test.XXXXXX)
WALLET_DIR="${DATA_DIR}/wallets"
GENESIS_WALLET="${DATA_DIR}/genesis-wallet.json"
FIBER_TOML="${DATA_DIR}/fiber.toml"
DIST_ADDRS="${DATA_DIR}/dist-addrs.txt"
DIST_SEEDS="${DATA_DIR}/dist-seeds.csv"

if [[ ! "$DATA_DIR" ]]; then
    echo "Could not create temp dir"
    exit 1
fi

echo "Test artifacts directory: $DATA_DIR"

BINARY="${COIN}-fibercoin-test.test"

cleanup() {
    echo "Cleaning up..."
    if [[ -n "${SKYCOIN_PID:-}" ]]; then
        kill -s SIGINT "$SKYCOIN_PID" 2>/dev/null || true
        wait "$SKYCOIN_PID" 2>/dev/null || true
    fi
    rm -f "$BINARY"
    rm -rf "$DATA_DIR"
    echo "Cleanup complete"
}
trap cleanup EXIT

# ─── Step 1: Generate genesis wallet ───────────────────────────────────────────
echo "=== Step 1: Generate genesis wallet ==="
GENESIS_SEED="test seed for fibercoin ci $(date +%s)"
go run . cli addressGen -n 1 -m json -s "$GENESIS_SEED" > "$GENESIS_WALLET"

# Extract genesis credentials from wallet JSON
GENESIS_ADDRESS=$(go run . cli addressGen -n 1 -m addrs -s "$GENESIS_SEED")
GENESIS_PUBKEY=$(python3 -c "import json,sys; w=json.load(open('$GENESIS_WALLET')); print(w['entries'][0]['public_key'])")
GENESIS_SECKEY=$(python3 -c "import json,sys; w=json.load(open('$GENESIS_WALLET')); print(w['entries'][0]['secret_key'])")

echo "Genesis address: $GENESIS_ADDRESS"
echo "Genesis pubkey:  $GENESIS_PUBKEY"

# ─── Step 2: Create initial fiber.toml ─────────────────────────────────────────
echo "=== Step 2: Create test fiber.toml ==="
GENESIS_TIMESTAMP=$(date +%s)

cat > "$FIBER_TOML" <<EOF
[node]
genesis_signature_str = ""
genesis_address_str = "${GENESIS_ADDRESS}"
blockchain_pubkey_str = "${GENESIS_PUBKEY}"
blockchain_seckey_str = "${GENESIS_SECKEY}"
genesis_timestamp = ${GENESIS_TIMESTAMP}
default_connections = []
display_name = "Fibercoin"
ticker = "FIB"
coin_hours_display_name = "Fiber Hours"
coin_hours_display_name_singular = "Fiber Hour"
coin_hours_ticker = "FH"

[params]
max_coin_supply = 100000000
initial_unlocked_count = 25
unlock_address_rate = 5
unlock_time_interval = 31536000
EOF

echo "Created fiber.toml at $FIBER_TOML"

# ─── Step 3: Generate distribution addresses ───────────────────────────────────
echo "=== Step 3: Generate distribution addresses ==="
FIBER_TOML="$FIBER_TOML" go run . cli fiberAddressGen \
    --num 100 \
    --overwrite \
    --addrs-file "$DIST_ADDRS" \
    --seeds-file "$DIST_SEEDS"

echo "Generated 100 distribution addresses"

# ─── Step 4: Compile and start daemon ──────────────────────────────────────────
echo "=== Step 4: Compile and start fibercoin daemon ==="

COMMIT=$(git rev-parse HEAD)
BRANCH=$(git rev-parse --abbrev-ref HEAD)
CMDPKG=$(go list ./cmd/${COIN})
GOLDFLAGS="-X ${CMDPKG}/commands.Commit=${COMMIT} -X ${CMDPKG}/commands.Branch=${BRANCH}"

go test -c -ldflags "${GOLDFLAGS}" -tags testrunmain -o "$BINARY" ./cmd/${COIN}/

FIBER_TOML="$FIBER_TOML" GENESIS="$GENESIS_WALLET" \
    ./"$BINARY" --disable-outgoing \
        --disable-incoming \
        --web-interface-port=$PORT \
        --download-peerlist=false \
        --launch-browser=false \
        --data-dir="$DATA_DIR" \
        --enable-all-api-sets=true \
        --wallet-dir="$WALLET_DIR" \
        --disable-csrf \
        --block-publisher \
        --test.run "^TestRunMain$" \
        &

SKYCOIN_PID=$!
echo "Fibercoin daemon pid=$SKYCOIN_PID"

echo "Waiting for daemon startup..."
sleep 5

# Wait for API to be available (up to 30 seconds)
for i in $(seq 1 30); do
    if curl -s "$HOST/api/v1/health" > /dev/null 2>&1; then
        echo "Daemon is ready"
        break
    fi
    if [[ $i -eq 30 ]]; then
        echo "FAIL: Daemon did not start within 30 seconds"
        exit 1
    fi
    sleep 1
done

# ─── Step 5: Verify genesis block ─────────────────────────────────────────────
echo "=== Step 5: Verify genesis block ==="

BLOCK0=$(curl -s "$HOST/api/v1/block?seq=0")
echo "Genesis block: $BLOCK0"

# Verify the genesis block exists and has correct address
BLOCK0_ADDR=$(echo "$BLOCK0" | python3 -c "import json,sys; b=json.load(sys.stdin); print(b['body']['txns'][0]['outputs'][0]['dst'])")
if [[ "$BLOCK0_ADDR" != "$GENESIS_ADDRESS" ]]; then
    echo "FAIL: Genesis block address mismatch. Expected $GENESIS_ADDRESS, got $BLOCK0_ADDR"
    exit 1
fi
echo "PASS: Genesis block address matches: $BLOCK0_ADDR"

# Verify genesis signature was written back to fiber.toml
UPDATED_SIG=$(python3 -c "
try:
    import tomllib
except ImportError:
    import tomli as tomllib
with open('$FIBER_TOML', 'rb') as f:
    cfg = tomllib.load(f)
print(cfg['node']['genesis_signature_str'])
")
if [[ -z "$UPDATED_SIG" ]]; then
    echo "FAIL: Genesis signature was not written back to fiber.toml"
    exit 1
fi
echo "PASS: Genesis signature written to fiber.toml: ${UPDATED_SIG:0:32}..."

# ─── Step 6: Distribute genesis coins ─────────────────────────────────────────
echo "=== Step 6: Distribute genesis coins ==="

RPC_ADDR="$RPC_ADDR" go run . cli distributeGenesis "$GENESIS_SECKEY"

echo "Distribution transaction submitted"

# Wait for block publisher to create block 1 (block creation interval is 10s)
echo "Waiting for block 1 to be created..."
for i in $(seq 1 30); do
    BLOCK1_CHECK=$(curl -s -o /dev/null -w "%{http_code}" "$HOST/api/v1/block?seq=1")
    if [[ "$BLOCK1_CHECK" == "200" ]]; then
        echo "Block 1 created after ~${i}s"
        break
    fi
    if [[ $i -eq 30 ]]; then
        echo "FAIL: Block 1 was not created within 30 seconds"
        exit 1
    fi
    sleep 1
done

# ─── Step 7: Verify distribution ──────────────────────────────────────────────
echo "=== Step 7: Verify distribution ==="

# Check that block 1 exists
BLOCK1=$(curl -s "$HOST/api/v1/block?seq=1")
echo "Distribution block: $BLOCK1"

BLOCK1_SEQ=$(echo "$BLOCK1" | python3 -c "import json,sys; b=json.load(sys.stdin); print(b['header']['seq'])")
if [[ "$BLOCK1_SEQ" != "1" ]]; then
    echo "FAIL: Distribution block not found (expected seq=1)"
    exit 1
fi
echo "PASS: Distribution block exists at seq=1"

# Verify the distribution block has the right number of outputs (100 distribution addresses)
BLOCK1_OUTPUTS=$(echo "$BLOCK1" | python3 -c "import json,sys; b=json.load(sys.stdin); print(len(b['body']['txns'][0]['outputs']))")
if [[ "$BLOCK1_OUTPUTS" != "100" ]]; then
    echo "FAIL: Distribution block has $BLOCK1_OUTPUTS outputs, expected 100"
    exit 1
fi
echo "PASS: Distribution block has 100 outputs"

# Check blockchain metadata shows 2 blocks
BM=$(curl -s "$HOST/api/v1/blockchain/metadata")
HEAD_SEQ=$(echo "$BM" | python3 -c "import json,sys; b=json.load(sys.stdin); print(b['head']['seq'])")
if [[ "$HEAD_SEQ" != "1" ]]; then
    echo "FAIL: Blockchain head should be seq=1, got $HEAD_SEQ"
    exit 1
fi
echo "PASS: Blockchain head is at seq=1"

# Check that genesis address balance is now 0
GENESIS_BAL=$(curl -s "$HOST/api/v1/balance?addrs=$GENESIS_ADDRESS")
GENESIS_COINS=$(echo "$GENESIS_BAL" | python3 -c "import json,sys; b=json.load(sys.stdin); print(b['confirmed']['coins'])")
if [[ "$GENESIS_COINS" != "0" ]]; then
    echo "FAIL: Genesis address should have 0 coins after distribution, got $GENESIS_COINS"
    exit 1
fi
echo "PASS: Genesis address balance is 0 after distribution"

# Check coin supply
COIN_SUPPLY=$(curl -s "$HOST/api/v1/coinSupply")
TOTAL_SUPPLY=$(echo "$COIN_SUPPLY" | python3 -c "import json,sys; b=json.load(sys.stdin); print(b['total_supply'])")
echo "Total supply: $TOTAL_SUPPLY"

# ─── Step 8: Test BIP44 chain selection ─────────────────────────────────────
echo "=== Step 8: Test BIP44 chain selection ==="

BIP44_SEED="abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"

# Create a BIP44 wallet
BIP44_CREATE=$(curl -s -X POST "$HOST/api/v1/wallet/create" \
    -d "type=bip44" \
    -d "seed=$BIP44_SEED" \
    -d "label=test-bip44")
echo "BIP44 wallet create response: $BIP44_CREATE"

BIP44_WALLET_ID=$(echo "$BIP44_CREATE" | python3 -c "import json,sys; w=json.load(sys.stdin); print(w['meta']['filename'])")
if [[ -z "$BIP44_WALLET_ID" ]]; then
    echo "FAIL: Could not create BIP44 wallet"
    exit 1
fi
echo "PASS: Created BIP44 wallet: $BIP44_WALLET_ID"

# Generate an address on the external chain
EXT_RESP=$(curl -s -X POST "$HOST/api/v1/wallet/newAddress" \
    -d "id=$BIP44_WALLET_ID" \
    -d "num=1" \
    -d "chain=external")
echo "External chain address response: $EXT_RESP"

EXT_ADDR=$(echo "$EXT_RESP" | python3 -c "import json,sys; r=json.load(sys.stdin); print(r['addresses'][0])")
if [[ -z "$EXT_ADDR" ]]; then
    echo "FAIL: Could not generate external chain address"
    exit 1
fi
echo "PASS: Generated external chain address: $EXT_ADDR"

# Generate an address on the change chain
CHG_RESP=$(curl -s -X POST "$HOST/api/v1/wallet/newAddress" \
    -d "id=$BIP44_WALLET_ID" \
    -d "num=1" \
    -d "chain=change")
echo "Change chain address response: $CHG_RESP"

CHG_ADDR=$(echo "$CHG_RESP" | python3 -c "import json,sys; r=json.load(sys.stdin); print(r['addresses'][0])")
if [[ -z "$CHG_ADDR" ]]; then
    echo "FAIL: Could not generate change chain address"
    exit 1
fi
echo "PASS: Generated change chain address: $CHG_ADDR"

# Get wallet details and verify chain entries
BIP44_WALLET=$(curl -s "$HOST/api/v1/wallet?id=$BIP44_WALLET_ID")
echo "BIP44 wallet details: $BIP44_WALLET"

# Verify external_entries exist in account 0
EXT_ENTRY_COUNT=$(echo "$BIP44_WALLET" | python3 -c "import json,sys; w=json.load(sys.stdin); print(len(w['accounts'][0]['external_entries']))")
if [[ "$EXT_ENTRY_COUNT" -lt 1 ]]; then
    echo "FAIL: BIP44 wallet account 0 has no external_entries"
    exit 1
fi
echo "PASS: BIP44 wallet has $EXT_ENTRY_COUNT external entry/entries"

# Verify change_entries exist in account 0
CHG_ENTRY_COUNT=$(echo "$BIP44_WALLET" | python3 -c "import json,sys; w=json.load(sys.stdin); print(len(w['accounts'][0]['change_entries']))")
if [[ "$CHG_ENTRY_COUNT" -lt 1 ]]; then
    echo "FAIL: BIP44 wallet account 0 has no change_entries"
    exit 1
fi
echo "PASS: BIP44 wallet has $CHG_ENTRY_COUNT change entry/entries"

# Verify external and change addresses are different
if [[ "$EXT_ADDR" == "$CHG_ADDR" ]]; then
    echo "FAIL: External address and change address are the same: $EXT_ADDR"
    exit 1
fi
echo "PASS: External ($EXT_ADDR) and change ($CHG_ADDR) addresses are different"

echo ""
echo "============================================"
echo "  ALL FIBERCOIN INTEGRATION TESTS PASSED"
echo "============================================"
