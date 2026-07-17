#!/bin/bash
set -e

echo "=========================================="
echo "    Droprcam Network Installation Script  "
echo "=========================================="

if [ ! -f "droprcam" ]; then
    echo "[-] Error: 'droprcam' binary not found!"
    echo "    Please run ./setup.sh to build it first."
    exit 1
fi

read -p "Enter your camera's IP address (e.g. 192.168.1.100): " CAM_IP

# Get the host's IP address (used for the python HTTP server)
HOST_IP=$(ip -o route get to $CAM_IP | sed -n 's/.*src \([0-9.]\+\).*/\1/p')

if [ -z "$HOST_IP" ]; then
    echo "[-] Error: Could not determine host IP address."
    exit 1
fi

echo "[+] Starting temporary HTTP server on $HOST_IP:8000..."
python3 -m http.server 8000 > /dev/null 2>&1 &
HTTP_PID=$!

# Ensure the HTTP server is killed when the script exits
trap "kill $HTTP_PID 2>/dev/null" EXIT

echo "[+] Connecting to camera at $CAM_IP via telnet..."
echo "[+] Transferring binary and updating boot script. Please wait..."

# We send a single compound command so the camera's shell executes it synchronously
# without relying on arbitrary sleep timers.
CMD="wget http://$HOST_IP:8000/droprcam -O /tmp/droprcam && chmod +x /tmp/droprcam && (killall droprcam 2>/dev/null; sleep 1); mv /tmp/droprcam /root/droprcam && grep -q droprcam /mnt/dropcam/ie_auto.sh || echo '(while true; do /root/droprcam; sleep 5; done) &' >> /mnt/dropcam/ie_auto.sh && echo 'Success' && sync && reboot"

(
  sleep 1
  echo "$CMD"
  sleep 15
) | telnet $CAM_IP || true

echo "[+] Waiting for camera to reboot and start Droprcam..."
sleep 20

# Curl loop to wait for the HTTP API to become available
until curl -s --max-time 1 "http://$CAM_IP:8080" > /dev/null 2>&1; do
    echo "  ... waiting for $CAM_IP:8080"
    sleep 2
done

echo ""
echo "[✓] Installation script finished!"
echo "[✓] Droprcam is successfully running and the API is reachable!"
