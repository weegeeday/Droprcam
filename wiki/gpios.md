# Dropcam Hardware & GPIO Reference Wiki

This wiki documents the physical hardware mapping, GPIO pins, and register interfaces for the Dropcam HD and Dropcam Pro camera families. These mappings were reverse-engineered from the stock Nest client (`/usr/bin/connect`) and its runtime payload scripts (`a5s_boot.sh` and `set_ir_brightness.sh`).

---

## 1. Hardware Architecture Overview

Both Dropcam HD and Dropcam Pro cameras are built on the **Ambarella A5s System-on-Chip (SoC)** (typically `a5m` running ARMv6). 

*   **Dropcam HD (HWVER 1 & 2):** Uses the **OmniVision OV9710** image sensor (720p H.264 video). Audio is managed by the **Asahi Kasei AK4642** codec.
*   **Dropcam Pro (HWVER 3, 4, & 5):** Uses the **Aptina AR0330** (or MT9T002) image sensor (1080p H.264 video). Audio is managed by the **Wolfson Microelectronics WM8974** codec.

---

## 2. GPIO Mappings by Hardware Revision

The Ambarella GPIO controller contains three banks (Bank 0: 0-31, Bank 1: 32-63, Bank 2: 64-95). Many pins are write-masked by default and must be unmasked in the hardware registers (e.g. `0x70009028` for Bank 0) before writing.

### 2.1 Dropcam HD (Hardware Version 2)

| GPIO Pin | Function | Direction | Default State / Values |
| :--- | :--- | :--- | :--- |
| **GPIO 11** | LP5521 LED Driver I2C Latch | Output | Pulse `1` -> `0` -> `1` to commit I2C brightness settings |
| **GPIO 16** | IR LED Ring Master Enable | Output | `1` = Enable IR LED Ring power; `0` = Disable |
| **GPIO 31** | Ambient Light Sensor (Photodiode) | Input | Poll state: `0` = Night (Dark); `1` = Day (Light) |
| **GPIO 45** | IR Cut Filter Enable (Strobe) | Output | Pulse `1` for 100ms and then `0` to latch H-bridge state |
| **GPIO 46** | IR Cut Filter Mode (Direction) | Output | `0` = Day Mode (Filter IN); `1` = Night Mode (Filter OUT) |
| **GPIO 88** | Status LED Channel (Yellow/Blue/White) | Output | `1` = Blue; `0` = White (when combined with 89) |
| **GPIO 89** | Status LED Channel (Yellow/Blue/White) | Output | `1` = Yellow |
| **GPIO 90** | Status LED Channel (Red) | Output | `1` = Red |

### 2.2 Dropcam Pro (Hardware Version 3)

| GPIO Pin | Function | Direction | Default State / Values |
| :--- | :--- | :--- | :--- |
| **GPIO 11** | LP5521 LED Driver I2C Latch | Output | Pulse `1` -> `0` -> `1` to commit I2C brightness settings |
| **GPIO 16** | IR LED Ring Master Enable | Output | `1` = Enable IR LED Ring power; `0` = Disable |
| **GPIO 45** | Status LED Channel (Yellow) | Output | `1` = Yellow LED ON |
| **GPIO 46** | Status LED Channel (Blue / ALS Input) | Bi-directional | LED control & Ambient Light Sensor input |
| **GPIO 88** | IR Cut Filter Enable (Strobe) | Output | Pulse `1` for 100ms and then `0` to latch H-bridge state |
| **GPIO 89** | IR Cut Filter Mode (Direction) | Output | `0` = Day Mode (Filter IN); `1` = Night Mode (Filter OUT) |
| **GPIO 90** | Status LED Channel (Red) | Output | `1` = Red LED ON |

### 2.3 Dropcam Pro (Hardware Version 4 & 5)

| GPIO Pin | Function | Direction | Default State / Values |
| :--- | :--- | :--- | :--- |
| **GPIO 11** | LP5521 LED Driver I2C Latch | Output | Pulse `1` -> `0` -> `1` to commit I2C brightness settings |
| **GPIO 16** | IR LED Ring Master Enable | Output | `1` = Enable IR LED Ring power; `0` = Disable |
| **GPIO 45** | IR Cut Filter Enable (Strobe) | Output | Pulse `1` for 100ms and then `0` to latch H-bridge state |
| **GPIO 46** | IR Cut Filter Mode (Direction) | Output | `0` = Day Mode (Filter IN); `1` = Night Mode (Filter OUT) |
| **GPIO 90** | Status LED Channel (Red) | Output | `1` = Red LED ON |

---

## 3. Peripheral Controllers

### 3.1 Status LED Array (LP5521 / LP5523)
The status indicator on the front is controlled via I2C at **Address `0x05`** (I2C Bus 0). 
The brightness controls are exposed under the standard Linux sysfs class:
`/sys/class/leds/lp5521:channel0/brightness` (Green/Blue)
`/sys/class/leds/lp5521:channel1/brightness` (Red)
`/sys/class/leds/lp5521:channel2/brightness` (Blue)

To manually set the LED brightness using the driver register commands:
```bash
# Write directly to I2C slave 0x05
i2cset -y 0 0x05 0x01 0xff s     # Set register 1 (latch control)
i2cset -y 0 0x05 0x02 <value> s  # Set register 2 (brightness value 0-255)

# Pulse GPIO 11 to commit the change
echo 1 > /sys/class/gpio/gpio11/value
echo 0 > /sys/class/gpio/gpio11/value
echo 1 > /sys/class/gpio/gpio11/value
```

### 3.2 IR LED Ring Intensity (PWM Backlight)
The infrared LED ring brightness is driven by a Pulse Width Modulation (PWM0) signal mapped to the Linux PWM Backlight driver. It is controlled by writing a duty cycle between `0` (OFF) and `1000` (Max Brightness) to:
`/sys/class/backlight/pwm-backlight.0/brightness`

*Note: The hardware PWM block is gated by the video system clock. This control only works while active video encoding (like RTSP streaming) is running.*

### 3.3 Lens Motor Controller (Focus/Zoom)
On Dropcam HD (revision 1), the physical lens focusing stepper motor is controlled via the I2C bus at **Address `0x0c`** (I2C Bus 0). 
*   **Write Register:** Registers `0x03` (upper byte) and `0x04` (lower byte).
*   **Target Step Range:** `0` to `1023`.
