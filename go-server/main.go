package main

import (
	"encoding/json"
	"os/exec"
	"time"
	"math/rand"
	"fmt"
	"net/http"
	"os"
	"log"
	"path/filepath"
)

type GenerateNoteRequest struct {
	NoteCount    string `json:"noteCount"`
	Accent       string `json:"accent"`
	Articulation string `json:"articulation"`
	Pitch        string `json:"pitch"`
	DynamicStart string `json:"dynamicStart"`
	DynamicEnd   string `json:"dynamicEnd"`
}

type GenerateNoteResponse struct {
	ImageURL string `json:"imageUrl"`
	MidiURL  string `json:"midiUrl"`
}

func generateNoteHandler(w http.ResponseWriter, r *http.Request) {
    // Add this to allow CORS from anywhere (for dev only)
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
    w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")

    // Handle preflight (OPTIONS) request
    if r.Method == http.MethodOptions {
        w.WriteHeader(http.StatusNoContent)
        return
    }

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req GenerateNoteRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Generate a unique ID
	id := fmt.Sprintf("%d", rand.Intn(1e9))
	imageFile := filepath.Join("..", "output", id+".png")
	midiFile := filepath.Join("..", "output", id+".mid")

	// Construct CLI args
	cliPath := filepath.Join("..", "cli", "dquant_cli.exe")
	args := []string{
		"--output_id", id,
		"--notes", req.NoteCount,
		"--accent", req.Accent,
		"--articulation", req.Articulation,
		"--pitch", req.Pitch,
		"--dyn_start", req.DynamicStart,
		"--dyn_end", req.DynamicEnd,
	}

	fmt.Printf("Calling CLI: %s %v\n", cliPath, args)

	cmd := exec.Command(cliPath, args...)
	err = cmd.Run()
	if err != nil {
		http.Error(w, "CLI execution failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Wait up to 3 seconds for files to appear
	timeout := time.After(3 * time.Second)
	tick := time.Tick(100 * time.Millisecond)
	foundImage, foundMidi := false, false

	for !(foundImage && foundMidi) {
		select {
		case <-timeout:
			http.Error(w, "Timed out waiting for files", http.StatusInternalServerError)
			return
		case <-tick:
			if !foundImage {
				if _, err := os.Stat(imageFile); err == nil {
					foundImage = true
				}
			}
			if !foundMidi {
				if _, err := os.Stat(midiFile); err == nil {
					foundMidi = true
				}
			}
		}
	}

	resp := GenerateNoteResponse{
		ImageURL: "/output/" + id + ".png",
		MidiURL:  "/output/" + id + ".mid",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func main() {
	fmt.Println("You are running version 0.0.1")
	
	rand.Seed(time.Now().UnixNano())
	http.Handle("/html/", http.StripPrefix("/html/", http.FileServer(http.Dir("../html"))))
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("../images"))))
	http.HandleFunc("/api/generate_note", generateNoteHandler)
	fs := http.FileServer(http.Dir("../output"))
	http.Handle("/output/", http.StripPrefix("/output/", fs))

	fmt.Println("Starting server at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}