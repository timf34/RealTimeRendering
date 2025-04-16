package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"math"
)

// Simple implementation of Mixbox's polynomial evaluation for demonstration purposes
// This is a simplified version that doesn't require the LUT data
func simpleMixboxLerp(color1, color2 [3]uint8, t float64) [3]uint8 {
	// Convert to float for easier math
	r1, g1, b1 := float64(color1[0])/255.0, float64(color1[1])/255.0, float64(color1[2])/255.0
	r2, g2, b2 := float64(color2[0])/255.0, float64(color2[1])/255.0, float64(color2[2])/255.0
	
	// Simple implementation of a non-linear interpolation that approximates pigment mixing
	// This is NOT the actual Mixbox algorithm, but a simple approximation for demonstration
	
	// Apply gamma correction to simulate non-linear mixing
	gamma := 2.2
	
	// Convert to linear space
	r1Linear, g1Linear, b1Linear := math.Pow(r1, gamma), math.Pow(g1, gamma), math.Pow(b1, gamma)
	r2Linear, g2Linear, b2Linear := math.Pow(r2, gamma), math.Pow(g2, gamma), math.Pow(b2, gamma)
	
	// Mix in linear space with slight subtractive mixing effect
	// The 0.7 factor makes the mixed colors slightly darker to simulate subtractive color mixing
	mixR := (1-t)*r1Linear + t*r2Linear - 0.2*t*(1-t)
	mixG := (1-t)*g1Linear + t*g2Linear - 0.2*t*(1-t)
	mixB := (1-t)*b1Linear + t*b2Linear - 0.2*t*(1-t)
	
	// Clamp values
	mixR = math.Max(0, math.Min(1, mixR))
	mixG = math.Max(0, math.Min(1, mixG))
	mixB = math.Max(0, math.Min(1, mixB))
	
	// Convert back to sRGB
	mixR, mixG, mixB = math.Pow(mixR, 1/gamma), math.Pow(mixG, 1/gamma), math.Pow(mixB, 1/gamma)
	
	// Convert to 8-bit color
	return [3]uint8{
		uint8(math.Round(mixR * 255)),
		uint8(math.Round(mixG * 255)),
		uint8(math.Round(mixB * 255)),
	}
}

// RGB to Hex conversion
func rgbToHex(rgb [3]uint8) string {
	return fmt.Sprintf("#%02x%02x%02x", rgb[0], rgb[1], rgb[2])
}

func main() {
	// Define some pigment colors
	cadmiumYellow := [3]uint8{254, 236, 0}
	ultramarineBlue := [3]uint8{25, 0, 89}
	cadmiumRed := [3]uint8{255, 39, 2}
	phthaloGreen := [3]uint8{0, 60, 50}
	
	fmt.Println("=== Simple Mixbox Demo ===")
	fmt.Println("Note: This is a simplified approximation of the Mixbox algorithm.")
	fmt.Println("For accurate pigment mixing, the full implementation with LUT data is required.")
	fmt.Println()
	
	// Linear vs Mixbox comparison
	fmt.Printf("Yellow (%s) + Blue (%s):\n", 
		rgbToHex(cadmiumYellow), rgbToHex(ultramarineBlue))
	
	// 50% mix
	t := 0.5
	
	// Linear interpolation
	linearRGB := [3]uint8{
		uint8(float64(cadmiumYellow[0])*(1-t) + float64(ultramarineBlue[0])*t),
		uint8(float64(cadmiumYellow[1])*(1-t) + float64(ultramarineBlue[1])*t),
		uint8(float64(cadmiumYellow[2])*(1-t) + float64(ultramarineBlue[2])*t),
	}
	
	// Mixbox-like interpolation
	mixboxRGB := simpleMixboxLerp(cadmiumYellow, ultramarineBlue, t)
	
	fmt.Printf("  Linear RGB:  %s (%v)\n", rgbToHex(linearRGB), linearRGB)
	fmt.Printf("  Simple Mixbox-like: %s (%v)\n", rgbToHex(mixboxRGB), mixboxRGB)
	fmt.Println()
	
	// Red + Green example
	fmt.Printf("Red (%s) + Green (%s):\n", 
		rgbToHex(cadmiumRed), rgbToHex(phthaloGreen))
	
	// Linear interpolation
	linearRGB = [3]uint8{
		uint8(float64(cadmiumRed[0])*(1-t) + float64(phthaloGreen[0])*t),
		uint8(float64(cadmiumRed[1])*(1-t) + float64(phthaloGreen[1])*t),
		uint8(float64(cadmiumRed[2])*(1-t) + float64(phthaloGreen[2])*t),
	}
	
	// Mixbox-like interpolation
	mixboxRGB = simpleMixboxLerp(cadmiumRed, phthaloGreen, t)
	
	fmt.Printf("  Linear RGB:  %s (%v)\n", rgbToHex(linearRGB), linearRGB)
	fmt.Printf("  Simple Mixbox-like: %s (%v)\n", rgbToHex(mixboxRGB), mixboxRGB)
	
	// Create a gradient image to visualize the difference
	fmt.Println("\nCreating gradient image 'simple_gradient.png'...")
	width, height := 400, 100
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	
	// Draw the yellow-blue gradient
	for x := 0; x < width; x++ {
		t := float64(x) / float64(width-1)
		
		// Linear RGB interpolation (top half)
		linearRGB := [3]uint8{
			uint8(float64(cadmiumYellow[0])*(1-t) + float64(ultramarineBlue[0])*t),
			uint8(float64(cadmiumYellow[1])*(1-t) + float64(ultramarineBlue[1])*t),
			uint8(float64(cadmiumYellow[2])*(1-t) + float64(ultramarineBlue[2])*t),
		}
		
		// Mixbox interpolation (bottom half)
		mixboxRGB := simpleMixboxLerp(cadmiumYellow, ultramarineBlue, t)
		
		// Draw the linear RGB gradient on the top half
		for y := 0; y < height/2; y++ {
			img.Set(x, y, color.RGBA{linearRGB[0], linearRGB[1], linearRGB[2], 255})
		}
		
		// Draw the mixbox gradient on the bottom half
		for y := height/2; y < height; y++ {
			img.Set(x, y, color.RGBA{mixboxRGB[0], mixboxRGB[1], mixboxRGB[2], 255})
		}
	}
	
	// Save the image
	f, err := os.Create("simple_gradient.png")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer f.Close()
	
	err = png.Encode(f, img)
	if err != nil {
		fmt.Println("Error encoding PNG:", err)
		return
	}
	
	fmt.Println("Gradient image saved - showing Yellow to Blue mixing")
	fmt.Println("Top half: Traditional RGB interpolation")
	fmt.Println("Bottom half: Simple Mixbox-like pigment approximation")
	
	// Create a second gradient for red-green
	fmt.Println("\nCreating gradient image 'simple_gradient2.png'...")
	img2 := image.NewRGBA(image.Rect(0, 0, width, height))
	
	// Draw the red-green gradient
	for x := 0; x < width; x++ {
		t := float64(x) / float64(width-1)
		
		// Linear RGB interpolation (top half)
		linearRGB := [3]uint8{
			uint8(float64(cadmiumRed[0])*(1-t) + float64(phthaloGreen[0])*t),
			uint8(float64(cadmiumRed[1])*(1-t) + float64(phthaloGreen[1])*t),
			uint8(float64(cadmiumRed[2])*(1-t) + float64(phthaloGreen[2])*t),
		}
		
		// Mixbox interpolation (bottom half)
		mixboxRGB := simpleMixboxLerp(cadmiumRed, phthaloGreen, t)
		
		// Draw the linear RGB gradient on the top half
		for y := 0; y < height/2; y++ {
			img2.Set(x, y, color.RGBA{linearRGB[0], linearRGB[1], linearRGB[2], 255})
		}
		
		// Draw the mixbox gradient on the bottom half
		for y := height/2; y < height; y++ {
			img2.Set(x, y, color.RGBA{mixboxRGB[0], mixboxRGB[1], mixboxRGB[2], 255})
		}
	}
	
	// Save the second image
	f2, err := os.Create("simple_gradient2.png")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer f2.Close()
	
	err = png.Encode(f2, img2)
	if err != nil {
		fmt.Println("Error encoding PNG:", err)
		return
	}
	
	fmt.Println("Gradient image saved - showing Red to Green mixing")
	fmt.Println("Top half: Traditional RGB interpolation")
	fmt.Println("Bottom half: Simple Mixbox-like pigment approximation")
}