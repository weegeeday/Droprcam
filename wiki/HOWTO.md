# How To Install Droprcam (with dropcam-rooter)

This guide assumes your camera has already been successfully rooted using the `dropcam-rooter` exploit, meaning you have telnet access and the `ie_auto.sh` script is running on boot to keep the camera alive on your Wi-Fi network.

By default, `dropcam-rooter` intercepts the stock Nest `connect` binary and replaces it with a dummy script. We are going to modify that script to launch our custom `droprcam` daemon instead, taking full control of the hardware.

## Automated Installation (Recommended)

We have provided a fully automated network installer that handles transferring the binary and patching the boot sequence for you.

1. On your host PC, navigate to the `Droprcam` directory.
2. Ensure you have built the binary first by running `./setup.sh`.
3. Run the installation script:
   ```bash
   ./install_to_camera.sh
   ```
4. Enter your camera's IP address when prompted.

The script will automatically:
* Spin up a temporary local HTTP server.
* Connect to your camera via Telnet.
* Download the compiled `droprcam` binary into the persistent `/root/` directory.
* Delete any existing `/mnt/dropcam/dummy_connect.sh` that might hijack the boot flow.
* Patch `/mnt/dropcam/ie_auto.sh` to permanently launch `droprcam` on every boot.
* Reboot the camera.

## Manual Installation

If you prefer to install it manually:

1. **Host PC:** Start a temporary HTTP server in the `Droprcam` folder:
   ```bash
   python3 -m http.server 8000
   ```

2. **Camera:** Telnet in, download the binary, and make it executable:
   ```bash
   wget http://<YOUR_PC_IP>:8000/droprcam -O /root/droprcam
   chmod +x /root/droprcam
   ```

3. **Camera:** Remove the `dummy_connect.sh` supervisor if it exists:
   ```bash
   rm -f /mnt/dropcam/dummy_connect.sh
   ```

4. **Camera:** Edit `/mnt/dropcam/ie_auto.sh` (e.g. `vi /mnt/dropcam/ie_auto.sh`) and replace the `sleep 3600` line inside the `dummy_connect` block with `/root/droprcam` and a `sleep 5` fallback:
   ```bash
   cat << 'EOF' > /tmp/dummy_connect
   #!/bin/sh
   while true; do
     /root/droprcam
     sleep 5
   done
   EOF
   ```

5. **Camera:** Save, flush, and reboot!
   ```bash
   sync
   reboot
   ```

## Verification
Once the camera reboots, you should immediately hear the camera's mechanical click (IR filter resetting) and see the blue status LED turn on. You are now ready to stream video via `rtsp://<camera_ip>/stream1` and control the hardware via the [HTTP API](api.md)!
