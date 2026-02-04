#!/bin/bash
# Script to cleanly restart the skyhw daemon

LOG_FILE="/tmp/skyhw.log"
DAEMON_BIN="/home/d0mo/go/bin/skyhw"
PORT=9510

echo "Stopping any existing skyhw processes..."
killall skyhw 2>/dev/null

# Wait and verify processes are dead
for i in {1..10}; do
    if ! pgrep -x skyhw > /dev/null; then
        echo "All skyhw processes stopped."
        break
    fi
    echo "Waiting for processes to stop... ($i/10)"
    sleep 1
done

# Force kill if still running
if pgrep -x skyhw > /dev/null; then
    echo "Force killing remaining processes..."
    killall -9 skyhw 2>/dev/null
    sleep 2
fi

# Verify port is free
if ss -tlnp | grep -q ":${PORT}"; then
    echo "ERROR: Port ${PORT} still in use!"
    ss -tlnp | grep ":${PORT}"
    exit 1
fi

echo "Starting daemon..."
"${DAEMON_BIN}" daemon > "${LOG_FILE}" 2>&1 &
DAEMON_PID=$!

echo "Daemon started with PID: ${DAEMON_PID}"

# Wait for daemon to be ready
for i in {1..15}; do
    if curl -s "http://127.0.0.1:${PORT}/api/v1/available" > /dev/null 2>&1; then
        echo "Daemon is ready!"
        break
    fi
    echo "Waiting for daemon to start... ($i/15)"
    sleep 1
done

# Test connection
echo ""
echo "Testing features endpoint..."
curl -s "http://127.0.0.1:${PORT}/api/v1/features" | head -20

echo ""
echo "Recent log entries:"
tail -10 "${LOG_FILE}"
