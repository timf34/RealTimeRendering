package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"mixbox"
	"net/http"
	"strconv"
)

// Hex to RGB conversion
func hexToRGB(hex string) ([3]uint8, error) {
	var r, g, b uint8
	if hex[0] == '#' {
		hex = hex[1:]
	}
	if len(hex) != 6 {
		return [3]uint8{}, fmt.Errorf("invalid hex color format: %s", hex)
	}
	r64, err := strconv.ParseUint(hex[0:2], 16, 8)
	if err != nil {
		return [3]uint8{}, err
	}
	g64, err := strconv.ParseUint(hex[2:4], 16, 8)
	if err != nil {
		return [3]uint8{}, err
	}
	b64, err := strconv.ParseUint(hex[4:6], 16, 8)
	if err != nil {
		return [3]uint8{}, err
	}
	return [3]uint8{uint8(r64), uint8(g64), uint8(b64)}, nil
}

// RGB to Hex conversion
func rgbToHex(rgb [3]uint8) string {
	return fmt.Sprintf("#%02x%02x%02x", rgb[0], rgb[1], rgb[2])
}

// ColorResult represents the result of a color mixing operation
type ColorResult struct {
	MixedColor   string `json:"mixedColor"`
	LinearColor  string `json:"linearColor"`
	MixedRGB     [3]uint8 `json:"mixedRGB"`
	LinearRGB    [3]uint8 `json:"linearRGB"`
}

func main() {
	// Initialize Mixbox with the compressed LUT
	// In a real implementation, you would load the actual LUT data
	lutData := "xNrFm..." // Replace with actual compressed data
	err := mixbox.InitLUT(lutData)
	if err != nil {
		log.Fatalf("Failed to initialize mixbox LUT: %v", err)
	}

	// Define HTML template for the web interface
	const htmlTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>Mixbox Go Demo</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            max-width: 800px;
            margin: 0 auto;
            padding: 20px;
        }
        h1 {
            color: #333;
        }
        .color-inputs {
            display: flex;
            gap: 20px;
            margin-bottom: 20px;
        }
        .color-input {
            display: flex;
            flex-direction: column;
            align-items: center;
        }
        .color-display {
            width: 100px;
            height: 100px;
            border: 1px solid #ccc;
            margin-top: 10px;
        }
        .results {
            display: flex;
            gap: 20px;
            margin-top: 20px;
        }
        .result {
            display: flex;
            flex-direction: column;
            align-items: center;
        }
        .slider-container {
            margin: 20px 0;
            width: 100%;
        }
        input[type="range"] {
            width: 100%;
        }
        .gradient {
            height: 50px;
            width: 100%;
            margin: 20px 0;
            background: linear-gradient(to right, #ffffff, #000000);
        }
    </style>
</head>
<body>
    <h1>Mixbox Go Demo</h1>
    <p>This demo shows the difference between traditional linear RGB color mixing and physically-based pigment mixing using Mixbox.</p>
    
    <div class="color-inputs">
        <div class="color-input">
            <label for="color1">Color 1:</label>
            <input type="color" id="color1" value="#ffeC00">
            <div class="color-display" id="color1-display" style="background-color: #ffeC00;"></div>
        </div>
        
        <div class="color-input">
            <label for="color2">Color 2:</label>
            <input type="color" id="color2" value="#0021aa">
            <div class="color-display" id="color2-display" style="background-color: #0021aa;"></div>
        </div>
    </div>
    
    <div class="slider-container">
        <label for="mixing-ratio">Mixing Ratio:</label>
        <input type="range" id="mixing-ratio" min="0" max="100" value="50">
        <span id="ratio-value">50%</span>
    </div>
    
    <div class="gradient" id="mixbox-gradient"></div>
    <div class="gradient" id="linear-gradient"></div>
    
    <div class="results">
        <div class="result">
            <h3>Linear RGB Mixing</h3>
            <div class="color-display" id="linear-result"></div>
            <span id="linear-hex">#000000</span>
        </div>
        
        <div class="result">
            <h3>Mixbox Pigment Mixing</h3>
            <div class="color-display" id="mixbox-result"></div>
            <span id="mixbox-hex">#000000</span>
        </div>
    </div>
    
    <script>
        // DOM elements
        const color1Input = document.getElementById('color1');
        const color2Input = document.getElementById('color2');
        const color1Display = document.getElementById('color1-display');
        const color2Display = document.getElementById('color2-display');
        const mixingRatio = document.getElementById('mixing-ratio');
        const ratioValue = document.getElementById('ratio-value');
        const linearResult = document.getElementById('linear-result');
        const mixboxResult = document.getElementById('mixbox-result');
        const linearHex = document.getElementById('linear-hex');
        const mixboxHex = document.getElementById('mixbox-hex');
        const mixboxGradient = document.getElementById('mixbox-gradient');
        const linearGradient = document.getElementById('linear-gradient');

        // Update UI based on inputs
        function updateColor1() {
            color1Display.style.backgroundColor = color1Input.value;
            updateMixing();
        }

        function updateColor2() {
            color2Display.style.backgroundColor = color2Input.value;
            updateMixing();
        }

        function updateRatio() {
            ratioValue.textContent = mixingRatio.value + '%';
            updateMixing();
        }

        // Fetch mixed colors from the server
        async function updateMixing() {
            const color1 = color1Input.value;
            const color2 = color2Input.value;
            const ratio = mixingRatio.value / 100;

            try {
                const response = await fetch(`/mix?color1=${encodeURIComponent(color1)}&color2=${encodeURIComponent(color2)}&ratio=${ratio}`);
                const result = await response.json();
                
                // Update result displays
                mixboxResult.style.backgroundColor = result.mixedColor;
                linearResult.style.backgroundColor = result.linearColor;
                mixboxHex.textContent = result.mixedColor;
                linearHex.textContent = result.linearColor;
                
                // Update gradients
                updateGradients(color1, color2);
            } catch (error) {
                console.error('Error fetching color mix:', error);
            }
        }

        // Create gradient displays
        function updateGradients(color1, color2) {
            linearGradient.style.background = `linear-gradient(to right, ${color1}, ${color2})`;
            
            // For the mixbox gradient, we'll use a series of 10 color stops
            let mixboxGradientStops = [];
            for (let i = 0; i <= 10; i++) {
                let ratio = i / 10;
                fetch(`/mix?color1=${encodeURIComponent(color1)}&color2=${encodeURIComponent(color2)}&ratio=${ratio}`)
                    .then(response => response.json())
                    .then(result => {
                        mixboxGradientStops.push({ ratio: ratio, color: result.mixedColor });
                        if (mixboxGradientStops.length === 11) {
                            // Sort by ratio
                            mixboxGradientStops.sort((a, b) => a.ratio - b.ratio);
                            
                            // Build gradient string
                            let gradientStr = 'linear-gradient(to right';
                            for (const stop of mixboxGradientStops) {
                                gradientStr += `, ${stop.color} ${stop.ratio * 100}%`;
                            }
                            gradientStr += ')';
                            
                            mixboxGradient.style.background = gradientStr;
                        }
                    });
            }
        }

        // Event listeners
        color1Input.addEventListener('input', updateColor1);
        color2Input.addEventListener('change', updateColor1);
        color2Input.addEventListener('input', updateColor2);
        color2Input.addEventListener('change', updateColor2);
        mixingRatio.addEventListener('input', updateRatio);
        mixingRatio.addEventListener('change', updateRatio);

        // Initialize
        updateMixing();
    </script>
</body>
</html>
`

	// Create template
	tmpl, err := template.New("index").Parse(htmlTemplate)
	if err != nil {
		log.Fatalf("Failed to parse template: %v", err)
	}

	// Define HTTP handlers
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, nil)
	})

	http.HandleFunc("/mix", func(w http.ResponseWriter, r *http.Request) {
		// Get parameters from the request
		color1Hex := r.URL.Query().Get("color1")
		color2Hex := r.URL.Query().Get("color2")
		ratioStr := r.URL.Query().Get("ratio")

		ratio, err := strconv.ParseFloat(ratioStr, 64)
		if err != nil {
			http.Error(w, "Invalid ratio", http.StatusBadRequest)
			return
		}

		// Convert hex colors to RGB
		color1, err := hexToRGB(color1Hex)
		if err != nil {
			http.Error(w, "Invalid color1", http.StatusBadRequest)
			return
		}

		color2, err := hexToRGB(color2Hex)
		if err != nil {
			http.Error(w, "Invalid color2", http.StatusBadRequest)
			return
		}

		// Mix colors using mixbox
		mixboxRGB := mixbox.Lerp(color1, color2, ratio)

		// Mix colors using linear interpolation for comparison
		linearRGB := [3]uint8{
			uint8(float64(color1[0])*(1-ratio) + float64(color2[0])*ratio),
			uint8(float64(color1[1])*(1-ratio) + float64(color2[1])*ratio),
			uint8(float64(color1[2])*(1-ratio) + float64(color2[2])*ratio),
		}

		// Convert back to hex
		mixboxHex := rgbToHex(mixboxRGB)
		linearHex := rgbToHex(linearRGB)

		// Create result
		result := ColorResult{
			MixedColor:  mixboxHex,
			LinearColor: linearHex,
			MixedRGB:    mixboxRGB,
			LinearRGB:   linearRGB,
		}

		// Send JSON response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	})

	// Start HTTP server
	fmt.Println("Starting Mixbox demo server on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}