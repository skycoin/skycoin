#!/bin/bash
# Runs "disable-seed-api"-mode tests against a skycoin node configured with -enable-seed-api=false
# and /api/v1/wallet/seed api endpoint should return 403 forbidden error.

# Set Script Name variable
SCRIPT=`basename ${BASH_SOURCE[0]}`

# Find unused port
PORT="1024"
while $(lsof -Pi :$PORT -sTCP:LISTEN -t >/dev/null) ; do
    PORT=$((PORT+1))
done

RPC_PORT="$PORT"
HOST="http://127.0.0.1:$PORT"
RPC_ADDR="http://127.0.0.1:$RPC_PORT"
MODE="disable-seed-api"
BINARY="skycoin-integration"
TEST=""
RUN_TESTS=""
# run go test with -v flag
VERBOSE=""

usage () {
  echo "Usage: $SCRIPT"
  echo "Optional command line arguments"
  echo "-t <string>  -- Test to run, api or cli; empty runs both tests"
  echo "-v <boolean> -- Run test with -v flag"
  exit 1
}

while getopts "h?t:r:v" args; do
  case $args in
    h|\?)
        usage;
        exit;;
    t ) TEST=${OPTARG};;
    v ) VERBOSE="-v";;
    r ) RUN_TESTS="-run ${OPTARG}";;
  esac
done

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
go build -o "$BINARY" cmd/skycoin/skycoin.go

# Run skycoin node with pinned blockchain database
echo "starting skycoin node in background with http listener on $HOST"

./skycoin-integration -disable-networking=true \
                      -web-interface-port=$PORT \
                      -download-peerlist=false \
                      -db-path=./src/api/integration/testdata/blockchain-180.db \
                      -db-read-only=true \
                      -rpc-interface=true \
                      -launch-browser=false \
                      -data-dir="$DATA_DIR" \
                      -wallet-dir="$WALLET_DIR" \
                      -enable-wallet-api=true \
                      -enable-seed-api=false &
SKYCOIN_PID=$!

echo "skycoin node pid=$SKYCOIN_PID"

echo "sleeping for startup"
sleep 3
echo "done sleeping"

set +e

if [[ -z $TEST || $TEST = "api" ]]; then

SKYCOIN_INTEGRATION_TESTS=1 SKYCOIN_INTEGRATION_TEST_MODE=$MODE SKYCOIN_NODE_HOST=$HOST WALLET_DIR=$WALLET_DIR \
    go test ./src/api/integration/... -timeout=30s $VERBOSE $RUN_TESTS

API_FAIL=$?

fi

if [[ -z $TEST  || $TEST = "cli" ]]; then

# SKYCOIN_INTEGRATION_TESTS=1 SKYCOIN_INTEGRATION_TEST_MODE=$MODE RPC_ADDR=$RPC_ADDR \
#     go test ./src/cli/integration/... -timeout=30s $VERBOSE $RUN_TESTS

CLI_FAIL=$?

fi


echo "shutting down skycoin node"

# Shutdown skycoin node
kill -s SIGINT $SKYCOIN_PID
wait $SKYCOIN_PID

rm "$BINARY"


if [[ (-z $TEST || $TEST = "api") && $API_FAIL -ne 0 ]]; then
  exit $API_FAIL
elif [[ (-z $TEST || $TEST = "cli") && $CLI_FAIL -ne 0 ]]; then
  exit $CLI_FAIL
else
  exit 0
fi
# exit $FAIL
