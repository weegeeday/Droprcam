# Dropcam Offline Audio Reference Wiki

This wiki documents the audio subsystem of the Dropcam HD (AK4642 codec) and Dropcam Pro (WM8974 codec), including ALSA mixer configurations, local playback/recording, real-time microphone streaming, and two-way intercom audio.

---

## 1. Audio Hardware Overview

*   **Dropcam HD (HWVER 1 & 2):** Uses the **Asahi Kasei AK4642** stereo audio codec connected to the Ambarella SoC via I2S.
*   **Dropcam Pro (HWVER 3, 4, & 5):** Uses the **Wolfson Microelectronics WM8974** mono audio codec connected via I2S.

The audio subsystem registers in the Linux kernel under the ALSA driver as:
*   **Card 0:** `DROPCAM`
*   **Device 0:** `AK4642-STREAM` or `WM8974-STREAM` (used for both playback and capture).

---

## 2. Mixer Initialization (Unmuting the Speaker & Mic)

By default at boot, the audio channels are muted. You must configure the mixer paths before recording or playing audio.

### 2.1 Dropcam HD (AK4642)
Run these commands to route and unmute the microphone and speaker:
```bash
# Configure the microphone
amixer sset 'Input Mux',0 'Both Mic'
amixer sset 'Mic Gain',0 '+20db'

# Configure the speaker
amixer sset 'Speaker Enable',0 on
amixer sset 'Speaker Gain',0 '+12.65db'
amixer sset 'Speaker Mixer SP',0 on
amixer sset 'ALC',0 on
```

### 2.2 Dropcam Pro (WM8974)
The Pro uses a pre-saved ALSA state file to restore mixer registers, or direct debug commands:
```bash
# Restore via state file
alsactl restore -f /tmp/connect/asound-crownroyal.state

# Manual register toggles (via SoC debug path)
echo 3 65 > /debug/asoc/DROPCAM/wm8974-codec.0-001a/codec_reg   # Enable DAC, SPKMIX, SPKP/N
echo 31 e > /debug/asoc/DROPCAM/wm8974-codec.0-001a/codec_reg   # Enable SPKBOOST and MONOBOOST
echo 2f 100 > /debug/asoc/DROPCAM/wm8974-codec.0-001a/codec_reg # Disable AUX and MIC boost
```

---

## 3. Local Audio Operations

### 3.1 Playing Local Audio (aplay)
Play a standard `.wav` or raw audio sample directly to the speaker:
```bash
aplay /tmp/connect/Q1_Connected_16K.wav
```

### 3.2 Recording Local Audio (arecord)
Record mono 16kHz 16-bit PCM audio from the microphone to a file:
```bash
arecord -d <duration_seconds> -f S16_LE -r 16000 -c 1 /tmp/test_mic.wav
```

---

## 4. Real-time Audio Streaming (Netcat Pipeline)

Because the stock RTSP server (`mediaserver`) is video-only, we stream audio over the network using `netcat` (nc).

### 4.1 Capturing Camera Mic -> Streaming to Host PC
Allows you to listen to the camera's microphone live.

1.  **On the Camera (Start listener):**
    ```bash
    arecord -f S16_LE -r 16000 -c 1 | nc -l -p 50005
    ```
2.  **On the Host PC (Receive and Play):**
    ```bash
    # Option A: Play using ALSA aplay
    nc <camera_ip> 50005 | aplay -f S16_LE -r 16000 -c 1
    
    # Option B: Play using low-latency ffplay
    ffplay -f s16le -ar 16000 -ac 1 -probesize 32 tcp://<camera_ip>:50005
    ```

### 4.2 Streaming Host PC Mic -> Playing on Camera Speaker
Allows you to speak into your PC microphone and have it play out of the camera speaker (intercom mode).

1.  **On the Camera (Start playback listener):**
    ```bash
    nc -l -p 50006 | aplay -f S16_LE -r 16000 -c 1
    ```
2.  **On the Host PC (Stream PC Mic to Camera):**
    ```bash
    arecord -f S16_LE -r 16000 -c 1 | nc <camera_ip> 50006
    ```
