#!/bin/bash

set -euxo pipefail

# Runs "live"-mode tests against a skycoin node that is already running
# "live" mode tests assume the blockchain data is active and may change at any time
# Data is checked for the appearance of correctness but the values themselves are not verified

PORT="6420"
HOST="http://127.0.0.1:$PORT"
MODE="live"

echo "checking if skycoin node is running"

http_proxy="" https_proxy="" wget $HOST 2>&1 >/dev/null

if [ ! $? -eq 0 ]; then
    echo "Skycoin node is not running on $HOST"
    exit 1
fi

SKYCOIN_INTEGRATION_TESTS=1 SKYCOIN_INTEGRATION_TEST_MODE=$MODE SKYCOIN_NODE_HOST=$HOST go test ./src/gui/integration/... -timeout=30s -v
