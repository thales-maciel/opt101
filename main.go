package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"os"
	"path/filepath"
	"runtime/pprof"
	"sync"

	"github.com/fogleman/gg"
)

const SRC = "image.jpg"
const WIDTH = 1280
const HEIGHT = 720
const NUM_WORKERS = 8
var opts jpeg.Options = jpeg.Options{Quality: 40}

func createFrame(img image.Image, t float64, frameIndex int, errChan chan error) {
	// Create a new canvas
	dc := gg.NewContext(WIDTH, HEIGHT)

	// Draw the image onto the canvas
	dc.DrawImage(img, int(1280 * t), 0)

	// Define the output path
	outPath := filepath.Join("frames", fmt.Sprintf("frame_%06d.jpg", frameIndex))

	// Create the output file
	outFile, err := os.Create(outPath)
	if err != nil {
		errChan <- fmt.Errorf("failed to create output file: %v", err)
	}
	defer outFile.Close()

	// Encode the image as JPEG with quality of 40
	err = jpeg.Encode(outFile, dc.Image(), &opts)
	if err != nil {
		errChan <- fmt.Errorf("failed to encode image: %v", err)
	}
}

func main() {
	cpuFile, err := os.Create("cpu.pprof")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(cpuFile)
	defer pprof.StopCPUProfile()

	frames := 1000

	// Load the image
	imgFile, err := os.Open(SRC)
	if err != nil {
		log.Fatalf("failed to open image: %v", err)
	}
	defer imgFile.Close()

	img, _, err := image.Decode(imgFile)
	if err != nil {
		log.Fatalf("failed to decode image: %v", err)
	}

	var wg sync.WaitGroup
	errChan := make(chan error, frames)

	workerChan := make(chan struct{}, NUM_WORKERS)

	for i := 0; i < frames; i++ {
		wg.Add(1)
		t := float64(i) / float64(frames)

		workerChan <- struct{}{}

		go func(t float64, i int) {
			defer func() { <-workerChan }() // Release worker slot
			defer wg.Done()
			createFrame(img, t, i, errChan)
		}(t, i)

	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		log.Println(err)
	}
}
