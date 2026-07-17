# Droprcam

Droprcam is a standalone, lightweight daemon designed to replace the proprietary Nest `/usr/bin/connect` binary on legacy Dropcam HD (and eventually Dropcam Pro) cameras. 

Instead of relying on the Nest cloud, Droprcam boots the hardware from a raw state, initializes the Ambarella DSP and video pipelines, routes the ALSA audio hardware, and exposes everything locally over standard network protocols. 

This completely frees the camera, turning it into a standard local IP camera perfect for integration into **Home Assistant**, **Scrypted**, or **HomeKit**.



## Plugins (Using this camera as a device on Home Assistant/other)
- [**Home Assistant**](https://github.com/weegeeday/Droprcam-HACS)

## Features

- **Zero Cloud Dependency:** Runs entirely on the camera's local area network.
- **Unified Hardware Control:** Initializes the image sensor (OV9710) and configures the Ambarella encoder (`mediaserver`) to expose standard RTSP video locally.
- **Hardware HTTP API:** Control the physical IR Cut Filter (Night Vision) and Status LEDs directly via REST endpoints.
- **Audio Routing:** Unmutes and configures the internal AK4642 stereo codec for immediate microphone and speaker access.
- **Live Microphone Feed:** Streams raw PCM audio straight from the microphone over HTTP for easy ingestion into ffmpeg.
- **Two-Way Audio:** Pipe raw audio back to the camera's speaker via the HTTP Intercom API.
- **Tiny Footprint:** Written in Go and compiled specifically for the Ambarella A5s ARM1136J-S (Linux 2.6.38) constraint (Sub-5MB).

## Directory Structure

* `/wiki/` - Documentation, reverse-engineering notes, and API references.
* `main.go` - The core daemon entrypoint.
* `hardware.go` - Controls hardware initialization (ALSA, `init.sh`, `mediaserver`) and GPIO mapping.
* `api.go` - Exposes the hardware over port 8080.
* `stream.go` - Prepares the network stream endpoints.
* `setup.sh` - Automated cross-compilation environment script.

## Getting Started

Because the Dropcam uses an aging Linux 2.6.38 kernel, standard Go binaries (1.21+) will crash due to unsupported syscalls, and UPX compression stubs often fail. We must compile this using **Go 1.20** with software floating point enabled.

### 1. Build the Daemon
Simply run the setup script. It will automatically download the correct Go 1.20 toolchain and compile the binary for you:
```bash
./setup.sh
```

### 2. Deploy to Camera
Once built, transfer the `droprcam` binary to your camera (which must be in a rooted state). Due to size constraints, place it on the `/root/` persistent filesystem.

If you are using the [dropcam-rooter](https://github.com/weegeeday/Dropcam-Rooter) exploit, please read the **[Installation HOWTO](wiki/HOWTO.md)** for exact instructions on how to hook the daemon into the boot sequence!

```bash
# On the camera telnet shell
cd /root
# Download binary...
chmod +x /root/droprcam

# Start the daemon manually (if not using the boot wrapper)
/root/droprcam
```

### 3. Usage & Integration
*   **Video:** The daemon starts the stock Ambarella encoder. Access the video stream at: `rtsp://<camera_ip>/stream1`
*   **Audio/Control:** Read the [API Wiki](wiki/api.md) for how to pull the microphone feed and control Night Vision.
