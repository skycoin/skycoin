#!/bin/bash

# Runs "live"-mode tests against a skycoin node that is already running
# "live" mode tests assume the blockchain data is active and may change at any time
# Data is checked for the appearance of correctness but the values themselves are not verified

#Set Script Name variable
SCRIPT=`basename ${BASH_SOURCE[0]}`
PORT="6420"
RPC_PORT="6430"
HOST="http://127.0.0.1:$PORT"
RPC_ADDR="127.0.0.1:$RPC_PORT"
MODE="live"
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

echo "checking if skycoin node is running"

http_proxy="" https_proxy="" wget -O- $HOST 2>&1 >/dev/null

if [ ! $? -eq 0 ]; then
    echo "Skycoin node is not running on $HOST"
    exit 1
fi

if [[ -z $TEST || $TEST = "gui" ]]; then

SKYCOIN_INTEGRATION_TESTS=1 SKYCOIN_INTEGRATION_TEST_MODE=$MODE SKYCOIN_NODE_HOST=$HOST go test ./src/gui/integration/... $UPDATE -timeout=3m $VERBOSE $RUN_TESTS $TEST_WALLET

fi

if [[ -z $TEST || $TEST = "cli" ]]; then

SKYCOIN_INTEGRATION_TESTS=1 SKYCOIN_INTEGRATION_TEST_MODE=$MODE RPC_ADDR=$RPC_ADDR SKYCOIN_NODE_HOST=$HOST go test ./src/api/cli/integration/... $UPDATE -timeout=3m $VERBOSE $RUN_TESTS $TEST_WALLET

fi
