#!/bin/bash

set -euxo pipefail

# Runs "stable"-mode tests against a skycoin node configured with a pinned database
# "stable" mode tests assume the blockchain data is static, in order to check API responses more precisely

PORT="46420"
HOST="http://127.0.0.1:$PORT"
MODE="stable"
BINARY="skycoin-integration"

DATA_DIR=$(mktemp -d -t skycoin-data-dir.XXXXXX)
WALLET_DIR="${DATA_DIR}/wallets"

if [[ ! "$DATA_DIR" ]]; then
  echo "Could not create temp dir"
  exit 1
fi

# Compile the skycoin node
# We can't use "go run" because this creates two processes which doesn't allow us to kill it at the end
echo "compiling skycoin"
go build -o "$BINARY" cmd/skycoin/skycoin.go

# Run skycoin node with pinned blockchain database
echo "starting skycoin node in background with http listener on $HOST"

./skycoin-integration -disable-networking=true \
                      -web-interface-port=$PORT \
                      -download-peerlist=false \
                      -db-path=./src/gui/integration/blockchain-180.db \
                      -rpc-interface=false \
                      -launch-browser=false \
                      -data-dir="$DATA_DIR" \
                      -wallet-dir="$WALLET_DIR" &
SKYCOIN_PID=$!

echo "skycoin node pid=$SKYCOIN_PID"

echo "sleeping for startup"
sleep 3
echo "done sleeping"

SKYCOIN_INTEGRATION_TESTS=1 SKYCOIN_INTEGRATION_TEST_MODE=$MODE SKYCOIN_NODE_HOST=$HOST go test ./src/gui/integration/... -timeout=30s -v

echo "shutting down skycoin node"

# Shutdown skycoin node
kill -s SIGINT $SKYCOIN_PID
wait $SKYCOIN_PID

rm "$BINARY"
