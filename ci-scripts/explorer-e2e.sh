#!/bin/bash
# Runs the explorer e2e suite against a skycoin node using a pinned blockchain
# database.
#
# The specs in explorer/test/specs assert exact values (block hashes, address
# balances, specific unconfirmed transactions), so they only pass against the
# pinned dataset in src/api/integration/testdata/blockchain-180.db. Pointing
# them at a live node makes the unconfirmed-transaction specs fail, because
# those transactions change.
#
# The explorer is served by the skycoin binary itself rather than by `ng serve`:
# `skycoin explorer` hosts the front-end and exposes the explorer-shaped API by
# proxying to a node. The committed bundle in explorer/dist is used deliberately
# so this exercises the artifact that actually ships. `make check-ui` separately
# guarantees that bundle matches the sources.

set -euo pipefail

# Absolute paths throughout: the wdio run below changes directory, and the exit
# trap must still be able to find the binary it has to delete.
REPO_ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
cd "$REPO_ROOT"

BINARY="${REPO_ROOT}/skycoin-explorer-e2e"
DB_PATH="${REPO_ROOT}/src/api/integration/testdata/blockchain-180.db"
EXPLORER_DIST="${REPO_ROOT}/explorer/dist"

NODE_PID=""
EXPLORER_PID=""
DATA_DIR=""

cleanup() {
  local status=$?
  set +e
  if [[ -n "$EXPLORER_PID" ]]; then
    kill "$EXPLORER_PID" 2>/dev/null
    wait "$EXPLORER_PID" 2>/dev/null
  fi
  if [[ -n "$NODE_PID" ]]; then
    kill -s SIGINT "$NODE_PID" 2>/dev/null
    wait "$NODE_PID" 2>/dev/null
  fi
  [[ -n "$DATA_DIR" ]] && rm -rf "$DATA_DIR"
  rm -f "$BINARY"
  exit $status
}
trap cleanup EXIT INT TERM

# Pick a free port starting from $1.
find_free_port() {
  local port=$1
  while (echo >"/dev/tcp/127.0.0.1/$port") >/dev/null 2>&1; do
    port=$((port + 1))
  done
  echo "$port"
}

# Poll a URL until it answers, rather than sleeping a fixed amount.
wait_for_http() {
  local url=$1 name=$2 tries=${3:-60}
  for _ in $(seq 1 "$tries"); do
    if curl -sf -o /dev/null --max-time 2 "$url"; then
      echo "$name is up: $url"
      return 0
    fi
    sleep 1
  done
  echo "ERROR: $name did not become ready at $url" >&2
  return 1
}

if [[ ! -f "$DB_PATH" ]]; then
  echo "ERROR: pinned blockchain database not found at $DB_PATH" >&2
  exit 1
fi

if [[ ! -f "$EXPLORER_DIST/index.html" ]]; then
  echo "ERROR: $EXPLORER_DIST/index.html not found. Run 'make build-ui' first." >&2
  exit 1
fi

NODE_PORT=$(find_free_port 8620)
EXPLORER_PORT=$(find_free_port $((NODE_PORT + 1)))
NODE_ADDR="http://127.0.0.1:${NODE_PORT}"
EXPLORER_ADDR="http://127.0.0.1:${EXPLORER_PORT}"

COMMIT=$(git rev-parse HEAD)
BRANCH=$(git rev-parse --abbrev-ref HEAD)
GOLDFLAGS="-X main.Commit=${COMMIT} -X main.Branch=${BRANCH}"

# Build once rather than using `go run`, which spawns a child process that the
# trap above could not reliably stop.
echo "compiling skycoin"
go build -o "$BINARY" -ldflags "${GOLDFLAGS}" .

DATA_DIR=$(mktemp -d -t skycoin-explorer-e2e.XXXXXX)

echo "starting skycoin node on ${NODE_ADDR} with the pinned database"
"$BINARY" daemon \
  --disable-networking=true \
  --web-interface-port="$NODE_PORT" \
  --download-peerlist=false \
  --db-path="$DB_PATH" \
  --db-read-only=true \
  --launch-browser=false \
  --data-dir="$DATA_DIR" \
  --enable-all-api-sets=true \
  >"${DATA_DIR}/node.log" 2>&1 &
NODE_PID=$!

wait_for_http "${NODE_ADDR}/api/v1/health" "skycoin node" || {
  tail -20 "${DATA_DIR}/node.log" >&2
  exit 1
}

echo "starting explorer on ${EXPLORER_ADDR}"
"$BINARY" explorer \
  --node-addr "$NODE_ADDR" \
  --server-host "127.0.0.1:${EXPLORER_PORT}" \
  --files-folder "$EXPLORER_DIST" \
  >"${DATA_DIR}/explorer.log" 2>&1 &
EXPLORER_PID=$!

wait_for_http "${EXPLORER_ADDR}/" "explorer" || {
  tail -20 "${DATA_DIR}/explorer.log" >&2
  exit 1
}

echo "running explorer e2e suite"
cd explorer
npx wdio run ./wdio-blockchain-180.conf.ts --baseUrl "$EXPLORER_ADDR"
