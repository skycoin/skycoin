#!/bin/bash
# Runs UI e2e tests against a skycoin node configured with a pinned database

# Set Script Name variable
SCRIPT=`basename ${BASH_SOURCE[0]}`

# Find unused port
PORT="1024"
while $(lsof -Pi :$PORT -sTCP:LISTEN -t >/dev/null) ; do
    PORT=$((PORT+1))
done

RPC_ADDR="127.0.0.1:$PORT"
HOST="http://127.0.0.1:$PORT"
BINARY="skycoin-integration"
E2E_PROXY_CONFIG=$(mktemp -t e2e-proxy.config.XXXXXX.js)

COMMIT=$(git rev-parse HEAD)
BRANCH=$(git rev-parse --abbrev-ref HEAD)
GOLDFLAGS="-X main.Commit=${COMMIT} -X main.Branch=${BRANCH}"

set -euxo pipefail

DATA_DIR=$(mktemp -d -t skycoin-data-dir.XXXXXX)
WALLET_DIR="${DATA_DIR}/wallets"

if [[ ! "$DATA_DIR" ]]; then
  echo "Could not create temp dir"
  exit 1
fi

# Compile the skycoin node
# We can't use "go run" because this creates two processes which doesn't allow us to kill it at the end
echo "compiling skycoin"
go build -o "$BINARY" -ldflags "${GOLDFLAGS}" cmd/skycoin/skycoin.go

# Run skycoin node with pinned blockchain database
echo "starting skycoin node in background with http listener on $HOST"

./skycoin-integration -disable-networking=true \
                      -web-interface-port=$PORT \
                      -download-peerlist=false \
                      -db-path=./src/api/integration/testdata/blockchain-180.db \
                      -db-read-only=true \
                      -launch-browser=false \
                      -data-dir="$DATA_DIR" \
                      -enable-wallet-api=true \
                      -wallet-dir="$WALLET_DIR" \
                      -enable-seed-api=true &
SKYCOIN_PID=$!

echo "skycoin node pid=$SKYCOIN_PID"

echo "sleeping for startup"
sleep 3
echo "done sleeping"

set +e


cat >$E2E_PROXY_CONFIG <<EOL
const PROXY_CONFIG = {
  "/api/*": {
    "target": "$HOST",
    "secure": false,
    "logLevel": "debug",
    "bypass": function (req) {
      req.headers["host"] = '$RPC_ADDR';
      req.headers["referer"] = '$HOST';
      req.headers["origin"] = '$HOST';
    }
  }
};
module.exports = PROXY_CONFIG;
EOL

# Run e2e tests
E2E_PROXY_CONFIG=$E2E_PROXY_CONFIG npm --prefix="./src/gui/static" run e2e-choose-config

RESULT=$?

echo "shutting down skycoin node"

# Shutdown skycoin node
kill -s SIGINT $SKYCOIN_PID
wait $SKYCOIN_PID

rm "$BINARY"
rm "$E2E_PROXY_CONFIG"

if [[ $RESULT -ne 0 ]]; then
  exit $RESULT
else
  exit 0
fi
