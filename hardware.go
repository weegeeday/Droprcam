package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

func InitHardware() error {
	fmt.Println("[+] Initializing hardware...")

	// 1. Run the camera init script
	cmd := exec.Command("/usr/local/bin/init.sh", "--ov9710")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("[-] Warning: init.sh failed: %v\n", err)
	}

	// 2. Load the microcode
	cmd = exec.Command("/usr/local/bin/load_ucode", "/lib/firmware")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("[-] Warning: load_ucode failed: %v\n", err)
	}

	// 3. Setup encode.cfg
	encodeCfg := `vin_mode        = 65524
vin_framerate   = 17083750
vout_type       = 1
vout_mode       = 65520
vin_mirror      = 0
vin_bayer       = 1
s0_type         = 1
s0_width        = 1280
s0_height       = 720
s0_brc          = 1
s0_vbr_min_bps  = 500000
s0_vbr_max_bps  = 1500000
s0_N            = 30
s0_start        = 1
s1_type         = 0
s2_type         = 0
s3_type         = 0
`
	if err := os.WriteFile("/tmp/encode.cfg", []byte(encodeCfg), 0644); err != nil {
		return fmt.Errorf("failed to write encode.cfg: %v", err)
	}

	// 4. Start mediaserver and image_server
	exec.Command("killall", "-9", "mediaserver", "image_server", "rtsp_server").Run()
	time.Sleep(1 * time.Second)

	fmt.Println("[+] Starting mediaserver...")
	mediaCmd := exec.Command("/usr/local/bin/mediaserver", "-a", "-f", "/tmp/encode.cfg")
	mediaCmd.Start()

	fmt.Println("[+] Starting image_server...")
	imageCmd := exec.Command("/usr/local/bin/image_server")
	imageCmd.Start()

	// 5. Configure ALSA for AK4642
	fmt.Println("[+] Configuring ALSA...")
	
	alsaScript := `
		amixer -q sset 'Input Mux',0 'Both Mic'
		amixer -q sset 'Mic Gain',0 '+20db'
		amixer -q sset 'Speaker Enable',0 on
		amixer -q sset 'Speaker Gain',0 '+12.65db'
		amixer -q sset 'Speaker Mixer SP',0 on
		amixer -q sset 'ALC',0 on
	`
	alsaCmd := exec.Command("sh", "-c", alsaScript)
	alsaCmd.Stdout = os.Stdout
	alsaCmd.Stderr = os.Stderr
	if err := alsaCmd.Run(); err != nil {
		fmt.Printf("[-] Warning: ALSA amixer init failed: %v\n", err)
	}

	// 6. Export GPIOs
	exportGPIO("45")
	exportGPIO("46")
	exportGPIO("16")
	exportGPIO("88")
	exportGPIO("89")
	exportGPIO("90")
	exportGPIO("11")

	setGPIODirection("45", "out")
	setGPIODirection("46", "out")
	setGPIODirection("16", "out")
	setGPIODirection("88", "out")
	setGPIODirection("89", "out")
	setGPIODirection("90", "out")
	setGPIODirection("11", "out")
	
	exportGPIO("31")
	setGPIODirection("31", "in")

	// Set Default State (Day Mode, LED Blue)
	DisableNightVision()
	SetStatusLED("blue")

	return nil
}

func exportGPIO(pin string) {
	os.WriteFile("/sys/class/gpio/export", []byte(pin), 0644)
}

func setGPIODirection(pin, dir string) {
	os.WriteFile(fmt.Sprintf("/sys/class/gpio/gpio%s/direction", pin), []byte(dir), 0644)
}

func setGPIOValue(pin, val string) {
	os.WriteFile(fmt.Sprintf("/sys/class/gpio/gpio%s/value", pin), []byte(val), 0644)
}

func EnableNightVision() {
	// Filter OUT
	setGPIOValue("46", "1")
	setGPIOValue("45", "1")
	time.Sleep(100 * time.Millisecond)
	setGPIOValue("45", "0")
	// IR LEDs ON
	setGPIOValue("16", "1")
	os.WriteFile("/sys/class/backlight/pwm-backlight.0/brightness", []byte("1000"), 0644)
}

func DisableNightVision() {
	// Filter IN
	setGPIOValue("46", "0")
	setGPIOValue("45", "1")
	time.Sleep(100 * time.Millisecond)
	setGPIOValue("45", "0")
	// IR LEDs OFF
	os.WriteFile("/sys/class/backlight/pwm-backlight.0/brightness", []byte("0"), 0644)
	setGPIOValue("16", "0")
}

func SetStatusLED(color string) {
	// Default off
	setGPIOValue("88", "0")
	setGPIOValue("89", "0")
	setGPIOValue("90", "0")

	switch color {
	case "blue":
		setGPIOValue("88", "1")
	case "yellow":
		setGPIOValue("89", "1")
	case "red":
		setGPIOValue("90", "1")
	case "white":
		// For white, 88=0 and 89=0? We saw blue=1(88), yellow=1(89), white=0(88) + 0(89).
		setGPIOValue("88", "0")
		setGPIOValue("89", "0")
	}
}


