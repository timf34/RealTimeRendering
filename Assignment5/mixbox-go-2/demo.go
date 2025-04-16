package main

import (
	"fmt"
	"mixbox"
	"image"
	"image/color"
	"image/png"
	"os"
)

// Define some common pigment colors from the mixbox documentation
var (
	CadmiumYellow     = [3]uint8{254, 236, 0}
	HansaYellow       = [3]uint8{252, 211, 0}
	CadmiumOrange     = [3]uint8{255, 105, 0}
	CadmiumRed        = [3]uint8{255, 39, 2}
	QuinacridoneMagenta = [3]uint8{128, 2, 46}
	CobaltViolet      = [3]uint8{78, 0, 66}
	UltramarineBlue   = [3]uint8{25, 0, 89}
	CobaltBlue        = [3]uint8{0, 33, 133}
	PhthaloBlue       = [3]uint8{13, 27, 68}
	PhthaloGreen      = [3]uint8{0, 60, 50}
	PermanentGreen    = [3]uint8{7, 109, 22}
	SapGreen          = [3]uint8{107, 148, 4}
	BurntSienna       = [3]uint8{123, 72, 0}
)

// hexToRGB converts a hex color string to RGB
func hexToRGB(hex string) [3]uint8 {
	var r, g, b uint8
	fmt.Sscanf(hex, "#%02x%02x%02x", &r, &g, &b)
	return [3]uint8{r, g, b}
}

// rgbToHex converts RGB to a hex color string
func rgbToHex(rgb [3]uint8) string {
	return fmt.Sprintf("#%02x%02x%02x", rgb[0], rgb[1], rgb[2])
}

func main() {
	// Initialize Mixbox with the compressed LUT
	// In a real implementation, you would load the actual LUT data
	// For this demo, we're using a placeholder
	lutData := "xNrFm..." // Replace with actual compressed data
	err := mixbox.InitLUT(lutData)
	if err != nil {
		fmt.Println("Failed to initialize mixbox LUT:", err)
		return
	}

	// Simple linear interpolation demo
	color1 := CadmiumYellow
	color2 := UltramarineBlue

	fmt.Printf("Mixing colors using natural pigment simulation:\n")
	fmt.Printf("Color 1: %s (%v)\n", rgbToHex(color1), color1)
	fmt.Printf("Color 2: %s (%v)\n", rgbToHex(color2), color2)

	// Create gradients with traditional RGB linear interpolation and mixbox
	fmt.Println("\nTraditional RGB vs Mixbox Gradient:")
	
	// Create a gradient image
	width, height := 400, 100
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	
	// Draw the gradient
	for x := 0; x < width; x++ {
		t := float64(x) / float64(width-1)
		
		// Simple linear RGB interpolation (for comparison)
		linearRGB := [3]uint8{
			uint8(float64(color1[0])*(1-t) + float64(color2[0])*t),
			uint8(float64(color1[1])*(1-t) + float64(color2[1])*t),
			uint8(float64(color1[2])*(1-t) + float64(color2[2])*t),
		}
		
		// Mixbox interpolation
		mixboxRGB := mixbox.Lerp(color1, color2, t)
		
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
	f, err := os.Create("gradient.png")
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
	
	fmt.Println("Gradient image saved as 'gradient.png'")
	fmt.Println("Top half: Traditional RGB interpolation")
	fmt.Println("Bottom half: Mixbox pigment-based interpolation")

	// Multi-color mixing demo
	fmt.Println("\nMulti-color mixing demo:")
	red := CadmiumRed
	green := PermanentGreen
	blue := CobaltBlue
	
	// Convert to latent space
	redLatent := mixbox.RGBToLatent(red)
	greenLatent := mixbox.RGBToLatent(green)
	blueLatent := mixbox.RGBToLatent(blue)
	
	// Mix in latent space (30% red, 50% green, 20% blue)
	mixLatent := [mixbox.LATENT_SIZE]float64{}
	for i := 0; i < mixbox.LATENT_SIZE; i++ {
		mixLatent[i] = 0.3*redLatent[i] + 0.5*greenLatent[i] + 0.2*blueLatent[i]
	}
	
	// Convert back to RGB
	mixedColor := mixbox.LatentToRGB(mixLatent)
	
	fmt.Printf("Red:    %s (%v)\n", rgbToHex(red), red)
	fmt.Printf("Green:  %s (%v)\n", rgbToHex(green), green)
	fmt.Printf("Blue:   %s (%v)\n", rgbToHex(blue), blue)
	fmt.Printf("Mixed:  %s (%v) - 30%% Red, 50%% Green, 20%% Blue\n", 
		rgbToHex(mixedColor), mixedColor)
}