package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"

	"github.com/timf34/mixbox-go/mixbox"
)

func main() {
	// Load the LUT
	err := mixbox.LoadLUTFromFile("lut.dat")
	if err != nil {
		log.Fatalf("Error loading LUT: %v", err)
	}

	// Define the colors to mix
	c1 := [3]uint8{255, 255, 0} // Red
	c2 := [3]uint8{0, 255, 0} // Blue

	// Create a new image
	width := 256
	height := 50
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for x := 0; x < width; x++ {
		t := float64(x) / float64(width-1)
		mixed := mixbox.Lerp(c1, c2, t)
		col := color.RGBA{mixed[0], mixed[1], mixed[2], 255}

		for y := 0; y < height; y++ {
			img.Set(x, y, col)
		}
	}

	// Save the image
	file, err := os.Create("gradient.png")
	if err != nil {
		log.Fatalf("Error creating image file: %v", err)
	}
	defer file.Close()

	err = png.Encode(file, img)
	if err != nil {
		log.Fatalf("Error encoding PNG: %v", err)
	}

	fmt.Println("✅ Gradient saved as gradient.png")
}
