# Droprcam HTTP API Reference

Droprcam exposes a lightweight HTTP REST API on port `8080`. This API is designed to allow smart home hubs (like Home Assistant or Scrypted) to control the camera's hardware and ingest its audio streams directly, bypassing the need for complex RTSP/ONVIF multiplexing on the low-power camera hardware.

## Base URL
```
http://<camera_ip>:8080
```

---

## 1. Hardware Control

### Toggle Night Vision ON
Manually engages the infrared cut filter (flips the filter OUT to allow IR light into the sensor) and turns on the front IR LED ring. 
*Note: Droprcam does not auto-switch based on light levels. You should use your smart home hub's sunset/sunrise triggers or a dedicated light sensor to hit this endpoint.*

* **URL:** `/night_vision/on`
* **Method:** `GET`
* **Success Response:**
  * **Code:** 200 OK
  * **Content:** `{"status":"ok", "message":"Night vision enabled"}`

### Toggle Night Vision OFF
Disengages the infrared cut filter (flips the filter IN to block IR light for normal daytime colors) and turns off the IR LED ring.

* **URL:** `/night_vision/off`
* **Method:** `GET`
* **Success Response:**
  * **Code:** 200 OK
  * **Content:** `{"status":"ok", "message":"Night vision disabled"}`

### Set Status LED
Changes the color of the front-facing status LED array. Supported colors on the Dropcam HD are `blue`, `yellow`, `red`, and `white`. Passing an unknown color (or `off`) will disable the LEDs. LED Colors are currently not accurate. It does change the led color, but not to the right one.

* **URL:** `/led/{color}`
* **Method:** `GET`
* **Success Response:**
  * **Code:** 200 OK
  * **Content:** `{"status":"ok", "message":"LED set to {color}"}`

---

## 2. Audio Streams

### Live Microphone Feed
Provides a continuous, real-time audio feed directly from the camera's AK4642 microphone. When this endpoint is accessed, Droprcam spawns an `arecord` subprocess and pipes the raw PCM audio bytes directly into the HTTP response.
This stream can be ingested natively by `ffmpeg` to mux with the video track.

* **URL:** `/mic`
* **Method:** `GET`
* **Stream Format:** Raw PCM, 16-bit Little Endian (S16_LE), 16000Hz, Mono.
* **Testing Command:**
  ```bash
  ffplay -f s16le -ar 16000 -ac 1 http://<camera_ip>:8080/mic
  ```

### Two-Way Intercom (Speaker Playback)
Allows you to play audio out of the camera's built-in speaker. You must POST a raw PCM audio payload in the body of the request. Droprcam pipes the request body directly into `aplay`.

* **URL:** `/intercom`
* **Method:** `POST`
* **Payload Format:** Raw PCM, 16-bit Little Endian (S16_LE), 16000Hz, Mono.
* **Success Response:**
  * **Code:** 200 OK
  * **Content:** `{"status":"ok", "message":"Audio played successfully"}`
* **Testing Command:**
  ```bash
  curl -X POST --data-binary @audio_file.raw http://<camera_ip>:8080/intercom
  ```
