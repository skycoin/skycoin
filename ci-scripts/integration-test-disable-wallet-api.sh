#!/bin/bash
# Runs "disable-wallet-api"-mode tests against a skycoin node configured with -disable-wallet-api option
# "disable-wallet-api"-mode confirms that no wallet related apis work, that the main index.html page
# does not load, and that a new wallet file is not created.

#Set Script Name variable
SCRIPT=`basename ${BASH_SOURCE[0]}`
PORT="46421"
RPC_PORT="46431"
HOST="http://127.0.0.1:$PORT"
RPC_ADDR="127.0.0.1:$RPC_PORT"
MODE="disable-wallet-api"
BINARY="skycoin-integration"
TEST=""
UPDATE=""
# run go test with -v flag
VERBOSE=""
# run go test with -run flag
RUN_TESTS=""
# run wallet tests
TEST_WALLET=""

usage () {
  echo "Usage: $SCRIPT"
  echo "Optional command line arguments"
  echo "-t <string>  -- Test to run, gui or cli; empty runs both tests"
  echo "-r <string>  -- Run test with -run flag"
  echo "-u <boolean> -- Update stable testdata"
  echo "-v <boolean> -- Run test with -v flag"
  echo "-w <boolean> -- Run wallet tests"
  exit 1
}

while getopts "h?t:r:uvw" args; do
  case $args in
    h|\?)
        usage;
        exit;;
    t ) TEST=${OPTARG};;
    r ) RUN_TESTS="-run ${OPTARG}";;
    u ) UPDATE="--update";;
    v ) VERBOSE="-v";;
    w ) TEST_WALLET="--test-wallet"
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
                      -db-path=./src/gui/integration/test-fixtures/blockchain-180.db \
                      -db-read-only=true \
                      -rpc-interface=true \
                      -rpc-interface-port=$RPC_PORT \
                      -launch-browser=false \
                      -data-dir="$DATA_DIR" \
                      -wallet-dir="$WALLET_DIR" \
                      -disable-wallet-api=true &
SKYCOIN_PID=$!

echo "skycoin node pid=$SKYCOIN_PID"

echo "sleeping for startup"
sleep 3
echo "done sleeping"

set +e

if [[ -z $TEST || $TEST = "gui" ]]; then

SKYCOIN_INTEGRATION_TESTS=1 SKYCOIN_INTEGRATION_TEST_MODE=$MODE SKYCOIN_NODE_HOST=$HOST WALLET_DIR=$WALLET_DIR go test ./src/gui/integration/... $UPDATE -timeout=30s $VERBOSE $RUN_TESTS $TEST_WALLET

GUI_FAIL=$?

fi

if [[ -z $TEST  || $TEST = "cli" ]]; then

SKYCOIN_INTEGRATION_TESTS=1 SKYCOIN_INTEGRATION_TEST_MODE=$MODE RPC_ADDR=$RPC_ADDR go test ./src/api/cli/integration/... $UPDATE -timeout=30s $VERBOSE $RUN_TESTS $TEST_WALLET

CLI_FAIL=$?

fi


echo "shutting down skycoin node"

# Shutdown skycoin node
kill -s SIGINT $SKYCOIN_PID
wait $SKYCOIN_PID

rm "$BINARY"


if [[ (-z $TEST || $TEST = "gui") && $GUI_FAIL -ne 0 ]]; then
  exit $GUI_FAIL
elif [[ (-z $TEST || $TEST = "cli") && $CLI_FAIL -ne 0 ]]; then
  exit $CLI_FAIL
else
  exit 0
fi
# exit $FAIL