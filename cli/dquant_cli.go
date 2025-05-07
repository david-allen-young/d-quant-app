package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"time"
)

func main() {
	var outputID, notes, accent, articulation, pitch, dynStart, dynEnd string

	flag.StringVar(&outputID, "output_id", "", "Unique output file ID")
	flag.StringVar(&notes, "notes", "", "Number of notes")
	flag.StringVar(&accent, "accent", "", "Accent type")
	flag.StringVar(&articulation, "articulation", "", "Articulation style")
	flag.StringVar(&pitch, "pitch", "", "Pitch")
	flag.StringVar(&dynStart, "dyn_start", "", "Starting dynamic")
	flag.StringVar(&dynEnd, "dyn_end", "", "Ending dynamic")
	flag.Parse()

	fmt.Printf("Received args:\n")
	fmt.Printf("output_id: %s\n", outputID)
	fmt.Printf("notes: %s\n", notes)
	fmt.Printf("accent: %s\n", accent)
	fmt.Printf("articulation: %s\n", articulation)
	fmt.Printf("pitch: %s\n", pitch)
	fmt.Printf("dyn_start: %s\n", dynStart)
	fmt.Printf("dyn_end: %s\n", dynEnd)

	// Simulate generation delay
	time.Sleep(1 * time.Second)

	outputDir := filepath.Join("..", "output")
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create output directory: %v\n", err)
		os.Exit(1)
	}

	// Write dummy MIDI file
	midiPath := filepath.Join(outputDir, outputID+".mid")
	if err := os.WriteFile(midiPath, []byte("MIDI-DATA"), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write MIDI file: %v\n", err)
		os.Exit(1)
	}

	// Write dummy PNG file
	pngPath := filepath.Join(outputDir, outputID+".png")
	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	img.Set(0, 0, color.RGBA{0, 0, 0, 0}) // fully transparent

	f, err := os.Create(pngPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create PNG file: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	err = png.Encode(f, img)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to encode PNG: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Dummy generation complete.")
}
