// ==========================================================
//  MIXBOX 2.0 (c) 2022 Secret Weapons. All rights reserved.
//  License: Creative Commons Attribution-NonCommercial 4.0
//  Authors: Sarka Sochorova and Ondrej Jamriska
//  Go implementation: [Your name]
// ==========================================================
package mixbox

import (
	"math"
)

// LATENT_SIZE is the size of the latent color representation
const LATENT_SIZE = 7

// The LUT data would normally be imported from a file
// In a real implementation, you would need to include the LUT data
// For brevity in this demo, we'll use a placeholder and assume
// the actual data would be imported elsewhere
var MixboxLUT []byte

// Lerp blends two colors together using natural pigment mixing
func Lerp(rgb1, rgb2 [3]uint8, t float64) [3]uint8 {
	latent1 := RGBToLatent(rgb1)
	latent2 := RGBToLatent(rgb2)

	latentMix := [LATENT_SIZE]float64{}
	for i := 0; i < LATENT_SIZE; i++ {
		latentMix[i] = (1.0-t)*latent1[i] + t*latent2[i]
	}

	return LatentToRGB(latentMix)
}

// LerpFloat blends two floating-point color values
func LerpFloat(rgb1, rgb2 [3]float64, t float64) [3]float64 {
	latent1 := FloatRGBToLatent(rgb1)
	latent2 := FloatRGBToLatent(rgb2)

	latentMix := [LATENT_SIZE]float64{}
	for i := 0; i < LATENT_SIZE; i++ {
		latentMix[i] = (1.0-t)*latent1[i] + t*latent2[i]
	}

	return LatentToFloatRGB(latentMix)
}

// LerpLinearFloat blends two linear-space floating-point color values
func LerpLinearFloat(rgb1, rgb2 [3]float64, t float64) [3]float64 {
	latent1 := LinearFloatRGBToLatent(rgb1)
	latent2 := LinearFloatRGBToLatent(rgb2)

	latentMix := [LATENT_SIZE]float64{}
	for i := 0; i < LATENT_SIZE; i++ {
		latentMix[i] = (1.0-t)*latent1[i] + t*latent2[i]
	}

	return LatentToLinearFloatRGB(latentMix)
}

// RGBToLatent converts an RGB color to its latent representation
func RGBToLatent(rgb [3]uint8) [LATENT_SIZE]float64 {
	return FloatRGBToLatent([3]float64{
		float64(rgb[0]) / 255.0,
		float64(rgb[1]) / 255.0,
		float64(rgb[2]) / 255.0,
	})
}

// LatentToRGB converts a latent representation back to RGB
func LatentToRGB(latent [LATENT_SIZE]float64) [3]uint8 {
	rgb := evalPolynomial(latent[0], latent[1], latent[2], latent[3])
	return [3]uint8{
		uint8(math.Round(clamp01(rgb[0]+latent[4]) * 255.0)),
		uint8(math.Round(clamp01(rgb[1]+latent[5]) * 255.0)),
		uint8(math.Round(clamp01(rgb[2]+latent[6]) * 255.0)),
	}
}

// FloatRGBToLatent converts a floating-point RGB color to its latent representation
func FloatRGBToLatent(rgb [3]float64) [LATENT_SIZE]float64 {
	r := clamp01(rgb[0])
	g := clamp01(rgb[1])
	b := clamp01(rgb[2])

	x := r * 63.0
	y := g * 63.0
	z := b * 63.0

	ix := int(x)
	iy := int(y)
	iz := int(z)

	tx := x - float64(ix)
	ty := y - float64(iy)
	tz := z - float64(iz)

	xyz := (ix + iy*64 + iz*64*64) & 0x3FFFF

	c0 := 0.0
	c1 := 0.0
	c2 := 0.0

	// In a real implementation, you would access the LUT here
	// This is just a placeholder to show how it would work
	if len(MixboxLUT) > 0 {
		w := (1.0 - tx) * (1.0 - ty) * (1.0 - tz)
		c0 += w * float64(MixboxLUT[xyz+192]) 
		c1 += w * float64(MixboxLUT[xyz+262336])
		c2 += w * float64(MixboxLUT[xyz+524480])

		w = tx * (1.0 - ty) * (1.0 - tz)
		c0 += w * float64(MixboxLUT[xyz+193])
		c1 += w * float64(MixboxLUT[xyz+262337])
		c2 += w * float64(MixboxLUT[xyz+524481])

		w = (1.0 - tx) * ty * (1.0 - tz)
		c0 += w * float64(MixboxLUT[xyz+256])
		c1 += w * float64(MixboxLUT[xyz+262400])
		c2 += w * float64(MixboxLUT[xyz+524544])

		w = tx * ty * (1.0 - tz)
		c0 += w * float64(MixboxLUT[xyz+257])
		c1 += w * float64(MixboxLUT[xyz+262401])
		c2 += w * float64(MixboxLUT[xyz+524545])

		w = (1.0 - tx) * (1.0 - ty) * tz
		c0 += w * float64(MixboxLUT[xyz+4288])
		c1 += w * float64(MixboxLUT[xyz+266432])
		c2 += w * float64(MixboxLUT[xyz+528576])

		w = tx * (1.0 - ty) * tz
		c0 += w * float64(MixboxLUT[xyz+4289])
		c1 += w * float64(MixboxLUT[xyz+266433])
		c2 += w * float64(MixboxLUT[xyz+528577])

		w = (1.0 - tx) * ty * tz
		c0 += w * float64(MixboxLUT[xyz+4352])
		c1 += w * float64(MixboxLUT[xyz+266496])
		c2 += w * float64(MixboxLUT[xyz+528640])

		w = tx * ty * tz
		c0 += w * float64(MixboxLUT[xyz+4353])
		c1 += w * float64(MixboxLUT[xyz+266497])
		c2 += w * float64(MixboxLUT[xyz+528641])

		c0 /= 255.0
		c1 /= 255.0
		c2 /= 255.0
	}

	c3 := 1.0 - (c0 + c1 + c2)

	mixrgb := evalPolynomial(c0, c1, c2, c3)

	return [LATENT_SIZE]float64{
		c0,
		c1,
		c2,
		c3,
		r - mixrgb[0],
		g - mixrgb[1],
		b - mixrgb[2],
	}
}

// LatentToFloatRGB converts a latent representation to floating-point RGB
func LatentToFloatRGB(latent [LATENT_SIZE]float64) [3]float64 {
	rgb := evalPolynomial(latent[0], latent[1], latent[2], latent[3])
	return [3]float64{
		clamp01(rgb[0] + latent[4]),
		clamp01(rgb[1] + latent[5]),
		clamp01(rgb[2] + latent[6]),
	}
}

// LinearFloatRGBToLatent converts a linear-space floating-point RGB color to latent
func LinearFloatRGBToLatent(rgb [3]float64) [LATENT_SIZE]float64 {
	return FloatRGBToLatent([3]float64{
		linearToSRGB(rgb[0]),
		linearToSRGB(rgb[1]),
		linearToSRGB(rgb[2]),
	})
}

// LatentToLinearFloatRGB converts a latent representation to linear RGB
func LatentToLinearFloatRGB(latent [LATENT_SIZE]float64) [3]float64 {
	rgb := LatentToFloatRGB(latent)
	return [3]float64{
		srgbToLinear(rgb[0]),
		srgbToLinear(rgb[1]),
		srgbToLinear(rgb[2]),
	}
}

// clamp01 clamps a value between 0 and 1
func clamp01(x float64) float64 {
	if x < 0.0 {
		return 0.0
	}
	if x > 1.0 {
		return 1.0
	}
	return x
}

// srgbToLinear converts an sRGB value to linear space
func srgbToLinear(x float64) float64 {
	if x >= 0.04045 {
		return math.Pow((x+0.055)/1.055, 2.4)
	}
	return x / 12.92
}

// linearToSRGB converts a linear space value to sRGB
func linearToSRGB(x float64) float64 {
	if x >= 0.0031308 {
		return 1.055*math.Pow(x, 1.0/2.4) - 0.055
	}
	return 12.92 * x
}

// evalPolynomial evaluates the polynomial mixing function
func evalPolynomial(c0, c1, c2, c3 float64) [3]float64 {
	c00 := c0 * c0
	c11 := c1 * c1
	c22 := c2 * c2
	c33 := c3 * c3
	c01 := c0 * c1
	c02 := c0 * c2
	c12 := c1 * c2

	r := 0.0
	g := 0.0
	b := 0.0

	w := c0 * c00
	r += +0.07717053 * w
	g += +0.02826978 * w
	b += +0.24832992 * w

	w = c1 * c11
	r += +0.95912302 * w
	g += +0.80256528 * w
	b += +0.03561839 * w

	w = c2 * c22
	r += +0.74683774 * w
	g += +0.04868586 * w
	b += +0.00000000 * w

	w = c3 * c33
	r += +0.99518138 * w
	g += +0.99978149 * w
	b += +0.99704802 * w

	w = c00 * c1
	r += +0.04819146 * w
	g += +0.83363781 * w
	b += +0.32515377 * w

	w = c01 * c1
	r += -0.68146950 * w
	g += +1.46107803 * w
	b += +1.06980936 * w

	w = c00 * c2
	r += +0.27058419 * w
	g += -0.15324870 * w
	b += +1.98735057 * w

	w = c02 * c2
	r += +0.80478189 * w
	g += +0.67093710 * w
	b += +0.18424500 * w

	w = c00 * c3
	r += -0.35031003 * w
	g += +1.37855826 * w
	b += +3.68865000 * w

	w = c0 * c33
	r += +1.05128046 * w
	g += +1.97815239 * w
	b += +2.82989073 * w

	w = c11 * c2
	r += +3.21607125 * w
	g += +0.81270228 * w
	b += +1.03384539 * w

	w = c1 * c22
	r += +2.78893374 * w
	g += +0.41565549 * w
	b += -0.04487295 * w

	w = c11 * c3
	r += +3.02162577 * w
	g += +2.55374103 * w
	b += +0.32766114 * w

	w = c1 * c33
	r += +2.95124691 * w
	g += +2.81201112 * w
	b += +1.17578442 * w

	w = c22 * c3
	r += +2.82677043 * w
	g += +0.79933038 * w
	b += +1.81715262 * w

	w = c2 * c33
	r += +2.99691099 * w
	g += +1.22593053 * w
	b += +1.80653661 * w

	w = c01 * c2
	r += +1.87394106 * w
	g += +2.05027182 * w
	b += -0.29835996 * w

	w = c01 * c3
	r += +2.56609566 * w
	g += +7.03428198 * w
	b += +0.62575374 * w

	w = c02 * c3
	r += +4.08329484 * w
	g += -1.40408358 * w
	b += +2.14995522 * w

	w = c12 * c3
	r += +6.00078678 * w
	g += +2.55552042 * w
	b += +1.90739502 * w

	return [3]float64{r, g, b}
}