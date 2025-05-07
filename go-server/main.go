package main

import (
	"crypto/rand"
	"fmt"
	"io"
	"net/http"
	"os"
	// "os/exec"
	"log"
	"path/filepath"
)

func createImageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("405 method not allowed"))
		return
	}
	
	// Get form values
	quarterNote := r.FormValue("quarterNote")
	halfNote := r.FormValue("halfNote")
	wholeNote := r.FormValue("wholeNote")
	fmt.Printf("Received values: %s and %s and %s\n", quarterNote, halfNote, wholeNote)

	// Execute the `date` command
	//cmd := exec.Command("date")
	//output, err := cmd.Output()
	//if err != nil {
	//	http.Error(w, fmt.Sprintf("Failed to execute date command: %v", err), http.StatusInternalServerError)
	//	return
	//}
	//fmt.Printf("Date: %s\n", output)

	// Create an image id that can be referenced later
	imageID := rand.Text()

	newUrl := fmt.Sprintf("/result?id=%s", imageID)
	http.Redirect(w, r, newUrl, http.StatusSeeOther)
}

func getImageHandler(w http.ResponseWriter, r *http.Request) {
	// filename := "/tmp/placeholder.png"
	filename := "placeholder.png"
	imageID := r.FormValue("id")
	if imageID != "" {
		//filename = fmt.Sprintf("/tmp/%s.png", imageID)
		filename = fmt.Sprintf("%s.png", imageID)
	}

	fmt.Printf("Checking for an image file at %s\n", filename)
	file, err := os.Open(filename)
	if err != nil {
		filePath := filepath.Join(".", "waiting.html")
		http.ServeFile(w, r, filePath)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Type", "image/png")

	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, "Error serving image", http.StatusInternalServerError)
		return
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	filePath := filepath.Join(".", "index.html")
	http.ServeFile(w, r, filePath)
}

func main() {
	fmt.Println("You are running version 0.0.1")
	
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/result", getImageHandler)
	http.HandleFunc("/api/create_image", createImageHandler)
	fmt.Println("Starting server at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}