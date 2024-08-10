package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"os"
	"path/filepath"
	"runtime/pprof"

	"github.com/fogleman/gg"
)

func createFrame(t float64, imagePath string, frameIndex int) error {
	// Load the image
	imgFile, err := os.Open(imagePath)
	if err != nil {
		return fmt.Errorf("failed to open image: %v", err)
	}
	defer imgFile.Close()

	img, _, err := image.Decode(imgFile)
	if err != nil {
		return fmt.Errorf("failed to decode image: %v", err)
	}

	// Create a new canvas
	const width = 1280
	const height = 720
	dc := gg.NewContext(width, height)

	// Calculate the x offset
	x := int(1280 * t)

	// Draw the image onto the canvas
	dc.DrawImage(img, x, 0)

	// Define the output path
	outPath := filepath.Join("frames", fmt.Sprintf("frame_%06d.jpg", frameIndex))

	// Create the output file
	outFile, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer outFile.Close()

	// Encode the image as JPEG with quality of 40
	opts := jpeg.Options{Quality: 40}
	err = jpeg.Encode(outFile, dc.Image(), &opts)
	if err != nil {
		return fmt.Errorf("failed to encode image: %v", err)
	}

	return nil
}

func main() {
	cpuFile, err := os.Create("cpu.pprof")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(cpuFile)
	defer pprof.StopCPUProfile()

	frames := 1000
	imagePath := "image.jpg"

	// Ensure the output directory exists
	if err := os.MkdirAll("frames", os.ModePerm); err != nil {
		log.Fatalf("failed to create frames directory: %v", err)
	}

	for i := 0; i < frames; i++ {
		t := float64(i) / float64(frames)

		err := createFrame(t, imagePath, i)
		if err != nil {
			log.Fatalf("failed to create frame %d: %v", i, err)
		}
	}
}
