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
	"strconv"
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
	var err error
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Generate a unique ID
	id := fmt.Sprintf("%d", rand.Intn(1e9))
	imageFile := filepath.Join("..", "output", id+".png")
	midiFile := filepath.Join("..", "go-server", "out", "midi", id+".mid")

	// Map numeric dynamic levels to string names
	dynMark := map[string]string{"1": "pp", "2": "p", "3": "mp", "4": "mf", "5": "f", "6": "ff"}
	dynStartStr := dynMark[req.DynamicStart]
	dynEndStr := dynMark[req.DynamicEnd]

	// Parse note count
	var noteCount int
	noteCount, err = strconv.Atoi(req.NoteCount)
	if err != nil {
		http.Error(w, "Invalid note count", http.StatusBadRequest)
		return
	}

	// Build note list (same pitch repeated for now)
	noteList := []map[string]interface{}{}
	for i := 0; i < noteCount; i++ {
		note := map[string]interface{}{
			"pitch":        req.Pitch + "4",
			"duration":     1.0,
			"articulation": req.Articulation,
			"accent":       req.Accent,
		}
		noteList = append(noteList, note)
	}

	// Build phrase JSON
	phraseJson := map[string]interface{}{
		"phrase": map[string]interface{}{
			"slur":      false, // can add toggle later
			"dyn_start": dynStartStr,
			"dyn_end":   dynEndStr,
		},
		"notes": noteList,
	}

	// Write phrase JSON to file
	jsonPath := filepath.Join("..", "output", "phrase_"+id+".json")
	jsonFile, err := os.Create(jsonPath)
	if err != nil {
		http.Error(w, "Failed to write phrase JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(jsonFile).Encode(phraseJson)
	jsonFile.Close()

	// Call real dquant_phrase_cli
	cliPath := filepath.Join("..", "cli", "dquant_cli.exe") // now actually dquant_phrase_cli.exe
	args := []string{
		"--input_json", jsonPath,
		"--song_json", filepath.Join("..", "cli", "song_context.json"),
		"--output_id", id,
	}

	fmt.Printf("Calling CLI: %s %v\n", cliPath, args)
	cmd := exec.Command(cliPath, args...)
	err = cmd.Run()
	if err != nil {
		http.Error(w, "CLI execution failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	fmt.Println("[DEBUG] CLI finished, checking if output MIDI exists:", midiFile)
	if _, err := os.Stat(midiFile); err != nil {
		fmt.Println("[ERROR] MIDI file not found:", midiFile)
	} else {
		fmt.Println("[DEBUG] MIDI file exists:", midiFile)
	}

	// Copy generated MIDI to public output folder
	finalMidiPath := filepath.Join("..", "output", id+".mid")
	fmt.Println("[DEBUG] Copying to public output folder:", finalMidiPath)
	err = copyFile(midiFile, finalMidiPath)
	if err != nil {
		fmt.Println("[ERROR] Copy failed:", err)
		http.Error(w, "Failed to copy MIDI to output folder: "+err.Error(), http.StatusInternalServerError)
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
				if _, err := os.Stat(finalMidiPath); err == nil {
					foundMidi = true
					fmt.Println("[DEBUG] Found final MIDI output:", finalMidiPath)
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

func copyFile(src, dst string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, input, 0644)
}
