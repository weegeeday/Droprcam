package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
)

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func sendJSON(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(Response{
		Status:  "ok",
		Message: msg,
	})
}

func StartAPI() error {
	http.HandleFunc("/night_vision/on", func(w http.ResponseWriter, r *http.Request) {
		EnableNightVision()
		sendJSON(w, http.StatusOK, "Night vision enabled")
	})

	http.HandleFunc("/night_vision/off", func(w http.ResponseWriter, r *http.Request) {
		DisableNightVision()
		sendJSON(w, http.StatusOK, "Night vision disabled")
	})

	http.HandleFunc("/led/", func(w http.ResponseWriter, r *http.Request) {
		color := r.URL.Path[len("/led/"):]
		SetStatusLED(color)
		sendJSON(w, http.StatusOK, fmt.Sprintf("LED set to %s", color))
	})

	http.HandleFunc("/intercom", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Pipe the HTTP request body directly to aplay
		// We expect RAW S16_LE 16kHz audio or similar. 
		// You might need to adjust the format parameters depending on what the client sends.
		cmd := exec.Command("aplay", "-f", "S16_LE", "-r", "16000", "-c", "1")
		cmd.Stdin = r.Body
		if err := cmd.Run(); err != nil {
			http.Error(w, fmt.Sprintf("aplay failed: %v", err), http.StatusInternalServerError)
			return
		}
		
		sendJSON(w, http.StatusOK, "Audio played successfully")
	})

	http.HandleFunc("/mic", func(w http.ResponseWriter, r *http.Request) {
		// Set headers for raw PCM audio stream
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Connection", "keep-alive")

		// Bind the command strictly to the HTTP request context.
		// When the client disconnects, the context cancels and kills arecord automatically!
		cmd := exec.CommandContext(r.Context(), "arecord", "-f", "S16_LE", "-r", "16000", "-c", "1")
		cmd.Stdout = w
		
		fmt.Println("[+] Client connected to /mic. Streaming audio...")
		if err := cmd.Run(); err != nil {
			fmt.Printf("[-] arecord stopped: %v\n", err)
		}
		fmt.Println("[+] Client disconnected from /mic.")
	})

	fmt.Println("[+] Starting HTTP API on port 8080...")
	return http.ListenAndServe(":8080", nil)
}
