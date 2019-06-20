#!/bin/bash

# Runs "live"-mode tests against a skycoin node that is already running
# "live" mode tests assume the blockchain data is active and may change at any time
# Data is checked for the appearance of correctness but the values themselves are not verified
# The skycoin node must be run with the wallet API enabled.

#Set Script Name variable
SCRIPT=`basename ${BASH_SOURCE[0]}`

COIN=${COIN:-skycoin}
PORT="6420"
RPC_PORT="$PORT"
HOST="http://127.0.0.1:$PORT"
RPC_ADDR="http://127.0.0.1:$RPC_PORT"
MODE="live"
TEST=""
UPDATE=""
TIMEOUT="60m"
# run go test with -v flag
VERBOSE=""
# run go test with -run flag
RUN_TESTS=""
# run wallet tests
TEST_LIVE_WALLET=""
FAILFAST=""
USE_CSRF=""
HEADER_CHECK="1"
DISABLE_NETWORKING=""

usage () {
  echo "Usage: $SCRIPT"
  echo "Optional command line arguments"
  echo "-t <string>  -- Test to run, api or cli; empty runs both tests"
  echo "-r <string>  -- Run test with -run flag"
  echo "-u <boolean> -- Update stable testdata"
  echo "-v <boolean> -- Run test with -v flag"
  echo "-w <boolean> -- Run wallet tests"
  echo "-f <boolean> -- Run test with -failfast flag"
  echo "-c <boolean> -- Pass this argument if the node has CSRF enabled"
  echo "-x <boolean> -- Pass this argument if the node has header check disabled"
  echo "-k <boolean> -- Run the tests that require networking disabled (live node must be run with -disable-networking)"
  exit 1
}

while getopts "h?t:r:uvwfck" args; do
case $args in
    h|\?)
        usage;
        exit;;
    t ) TEST=${OPTARG};;
    r ) RUN_TESTS="-run ${OPTARG}";;
    u ) UPDATE="--update";;
    v ) VERBOSE="-v";;
    w ) TEST_LIVE_WALLET="--test-live-wallet";;
    f ) FAILFAST="-failfast";;
    c ) USE_CSRF="1";;
    x ) HEADER_CHECK="";;
	k ) DISABLE_NETWORKING="true"
  esac
done

set -euxo pipefail

echo "checking if $COIN node is running"

HEALTH="$HOST/api/v1/health"

http_proxy="" https_proxy="" wget -O- $HEALTH 2>&1 >/dev/null

if [ ! $? -eq 0 ]; then
    echo "$COIN node is not running on $HOST"
    exit 1
fi

echo "checking if integration tests compile"
go test ./src/api/integration/...
go test ./src/cli/integration/...

if [[ -z $TEST || $TEST = "api" ]]; then

SKYCOIN_INTEGRATION_TESTS=1 SKYCOIN_INTEGRATION_TEST_MODE=$MODE SKYCOIN_NODE_HOST=$HOST \
	LIVE_DISABLE_NETWORKING=$DISABLE_NETWORKING \
    go test -count=1 ./src/api/integration/... $FAILFAST $UPDATE -timeout=$TIMEOUT $VERBOSE $RUN_TESTS $TEST_LIVE_WALLET

fi

if [[ -z $TEST || $TEST = "cli" ]]; then

SKYCOIN_INTEGRATION_TESTS=1 SKYCOIN_INTEGRATION_TEST_MODE=$MODE RPC_ADDR=$RPC_ADDR \
	SKYCOIN_NODE_HOST=$HOST USE_CSRF=$USE_CSRF HEADER_CHECK=$HEADER_CHECK LIVE_DISABLE_NETWORKING=$DISABLE_NETWORKING \
    go test -count=1 ./src/cli/integration/... $FAILFAST $UPDATE -timeout=$TIMEOUT $VERBOSE $RUN_TESTS $TEST_LIVE_WALLET

fi
