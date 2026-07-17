package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	fmt.Println("======================================")
	fmt.Println("      Dropcam Connect Replacement     ")
	fmt.Println("======================================")

	// 1. Initialize Hardware (ALSA, GPIOs, mediaserver)
	if err := InitHardware(); err != nil {
		log.Fatalf("Failed to initialize hardware: %v", err)
	}

	// 2. Start RTSP Multiplexer in a goroutine
	go func() {
		if err := StartStreamMuxer(); err != nil {
			log.Printf("Stream muxer stopped: %v", err)
		}
	}()

	// 3. Start HTTP API
	go func() {
		if err := StartAPI(); err != nil {
			log.Printf("HTTP API stopped: %v", err)
		}
	}()

	// Wait for termination signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	fmt.Println("\nShutting down Dropcam Daemon...")
}
