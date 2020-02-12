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
E2E_PROXY_CONFIG="./src/gui/static/proxy.config.js"
rm "$E2E_PROXY_CONFIG"

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

# Create a dummy wallet with an address existing in the blockchain-180.db dataset
mkdir "$WALLET_DIR"
cat >"${WALLET_DIR}/test_wallet.wlt" <<EOL
{
    "meta": {
        "coin": "skycoin",
        "cryptoType": "scrypt-chacha20poly1305",
        "encrypted": "true",
        "filename": "test_wallet.wlt",
        "label": "Test wallet",
        "lastSeed": "",
        "secrets": "dgB7Im4iOjEwNDg1NzYsInIiOjgsInAiOjEsImtleUxlbiI6MzIsInNhbHQiOiIvelgxOFdPQUlzK1FQOXZZWi9aVXlDVktmZWMzY29UdjNzU2h6cENmWDNvPSIsIm5vbmNlIjoid0Qxb0U5VldycW9RTmJKVyJ9qFmBxQnP42SKJsQavIW/8chLo3alLx/KZI/lFFU96iZhTeSAfLNtPajX+4bcAdsdsPPhoBLNRBBuy1O2NImjZOVEc3YPCpXQO2Zj6/AZKu6zRldSSRbyk2blLngHr9Iv2oS4CcofCUdQF6tfc8soU/Vef9pZAHEUn0Soi1i9iprK3trkq0CfgP3LR3faltBfTkJCkOOjNGbHgDfZrGL6TZpllxjEAlO2jzYqMvmucowq3MDlTplFMJoE5Fvw47gjSuOpdRQ0yK4EgTabXKZJbbjvWZzE9pCYuUE=",
        "seed": "",
        "tm": "1529948542",
        "type": "deterministic",
        "version": "0.2"
    },
    "entries": [
        {
            "address": "R6aHqKWSQfvpdo2fGSrq4F1RYXkBWR9HHJ",
            "public_key": "03cef9d4635c6f075a479415805134daa1b5fda6e0f6a82b154e04b26db6afa770",
            "secret_key": ""
        }
    ]
}
EOL

# Compile the skycoin node
# We can't use "go run" because that creates two processes which doesn't allow us to kill it at the end
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
                      -enable-all-api-sets=true \
                      -enable-api-sets=INSECURE_WALLET_SEED \
                      -wallet-dir="$WALLET_DIR" \
                      &
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
E2E_PROXY_CONFIG=$E2E_PROXY_CONFIG npm --prefix="./src/gui/static" run e2e

RESULT=$?

echo "shutting down skycoin node"

# Shutdown skycoin node
kill -s SIGINT $SKYCOIN_PID
wait $SKYCOIN_PID

rm "$BINARY"

if [[ $RESULT -ne 0 ]]; then
  exit $RESULT
else
  exit 0
fi
