# Mixbox for Go

This is a Go implementation of the Mixbox color mixing library, which simulates natural pigment mixing for more realistic color blending.

## About Mixbox

Mixbox is a library for natural color mixing based on real pigments, created by Secret Weapons. The original implementation is available for various platforms, and this is a Go port of the algorithm.

The key feature of Mixbox is that it simulates how actual pigments mix, rather than using simple RGB linear interpolation. This leads to more realistic and aesthetically pleasing color blends, especially for applications related to digital painting, design, and color theory.

## Color Mixing Comparison

Traditional RGB color mixing is additive and works by averaging color components. This often results in desaturated or muddy colors when mixing complementary colors. Mixbox, on the other hand, simulates subtractive color mixing similar to how physical pigments behave, producing more vibrant and natural-looking mixtures.

Examples of the difference:
- Yellow + Blue: In RGB mixing, this produces a grayish color. With Mixbox, you get a more natural green.
- Red + Green: In RGB mixing, this yields a muddy brown. With Mixbox, the result is closer to what you'd expect when mixing actual paint.

## Usage

### Basic Color Interpolation

```go
import "mixbox"

// Initialize the lookup table (required once before using the library)
lutData := "..." // The compressed LUT data (not included in this demo)
mixbox.InitLUT(lutData)

// Define two colors in RGB format
yellow := [3]uint8{254, 236, 0}  // Cadmium Yellow
blue := [3]uint8{25, 0, 89}      // Ultramarine Blue

// Mix the colors with a ratio (0.0 to 1.0)
// 0.0 means 100% yellow, 1.0 means 100% blue
mixedColor := mixbox.Lerp(yellow, blue, 0.5)  // 50% mixture

fmt.Printf("Mixed color: RGB(%d, %d, %d)\n", 
    mixedColor[0], mixedColor[1], mixedColor[2])
```

### Multi-Color Mixing

For mixing more than two colors with different proportions:

```go
// Convert colors to latent space
redLatent := mixbox.RGBToLatent(red)
greenLatent := mixbox.RGBToLatent(green)
blueLatent := mixbox.RGBToLatent(blue)

// Mix in latent space with different proportions
mixLatent := [mixbox.LATENT_SIZE]float64{}
for i := 0; i < mixbox.LATENT_SIZE; i++ {
    mixLatent[i] = 0.3*redLatent[i] + 0.5*greenLatent[i] + 0.2*blueLatent[i]
}

// Convert back to RGB
mixedColor := mixbox.LatentToRGB(mixLatent)
```

### Working with Floating Point Colors

If you need to work with floating-point color values (0.0 to 1.0 range):

```go
// Float colors (0.0 to 1.0 range)
color1 := [3]float64{1.0, 0.5, 0.2}
color2 := [3]float64{0.2, 0.8, 0.9}

// Mix using floating point functions
mixedColor := mixbox.LerpFloat(color1, color2, 0.5)
```

### Linear Color Space

For working with linear color spaces:

```go
// Linear space colors (0.0 to 1.0 range, not in sRGB space)
linearColor1 := [3]float64{1.0, 0.5, 0.2}
linearColor2 := [3]float64{0.2, 0.8, 0.9}

// Mix in linear space
mixedLinearColor := mixbox.LerpLinearFloat(linearColor1, linearColor2, 0.5)
```

## Predefined Pigment Colors

The library includes several predefined pigment colors that can be used as reference:

```go
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
```

## Demo Applications

This repository includes two demo applications:

1. **simple_demo.go**: A command-line application that demonstrates the difference between RGB linear interpolation and Mixbox pigment mixing. It also generates gradient images to visualize the difference.

2. **web_demo.go**: A web-based interactive demo that allows you to pick colors and see the difference between traditional RGB mixing and Mixbox mixing in real-time.

### Running the Simple Demo

```bash
go run simple_demo.go
```

This will output color comparisons and generate gradient images.

### Running the Web Demo

```bash
go run web_demo.go
```

Then open a browser and go to http://localhost:8080 to interact with the web interface.

## Notes

- This implementation requires the Mixbox LUT (Look-Up Table) data to work correctly. Due to licensing restrictions, the LUT data is not included in this repository.
- The original Mixbox library is licensed under Creative Commons Attribution-NonCommercial 4.0. For commercial use, you need to contact Secret Weapons at mixbox@scrtwpns.com.

## License

This Go port of Mixbox follows the same license as the original:
Creative Commons Attribution-NonCommercial 4.0

Original authors: Sarka Sochorova and Ondrej Jamriska