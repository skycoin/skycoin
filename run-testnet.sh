#!/usr/bin/env bash

# Runs a local Skycoin testnet

# Disclaimer : For the moment this script is a hack used for testing purposes
#              Proper implementation consists in launching/connecting to (sim)testnet fiber coin network

set -x

TEMP_DIR="/tmp/skytestnet.$1"
PEERS_FILE="${TEMP_DIR}/localhost-peers.txt"

echo "Creating temp dirs starting at $1"
echo "$1,$(expr $1 + 1),$(expr $1 + 2),$(expr $1 + 3),$(expr $1 + 4),$(expr $1 + 5),$(expr $1 + 6)" | tr , '\n' | xargs -I PORT mkdir -p "${TEMP_DIR}/PORT"
echo ""

echo "Creating local peers file"
echo "$1
$(expr $1 + 1)
$(expr $1 + 2)
$(expr $1 + 3)
$(expr $1 + 4)
$(expr $1 + 5)
$(expr $1 + 6)" | sed 's/^/127.0.0.1:/g' > ${PEERS_FILE}
cat ${PEERS_FILE}
echo ""

echo "Launching Skycoin nodes"
cut -d : -f 2 ${PEERS_FILE} | xargs -I SKYPORT screen -dm /bin/bash -c "PORT=SKYPORT; ./run-client.sh -localhost-only -custom-peers-file=$TEMP_DIR/localhost-peers.txt -download-peerlist=false -launch-browser=false -data-dir=$TEMP_DIR/\$PORT -web-interface-port=\$(expr \$PORT + 420) -port=\$PORT "


