package mixbox

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"io"
)

// InitLUT initializes the MixboxLUT global variable with the decompressed lookup table data
func InitLUT(compressedBase64 string) error {
	// Decode from base64
	compressed, err := base64.StdEncoding.DecodeString(compressedBase64)
	if err != nil {
		return err
	}

	// Create a zlib reader
	zlibReader, err := zlib.NewReader(bytes.NewReader(compressed))
	if err != nil {
		return err
	}
	defer zlibReader.Close()

	// Read the decompressed data
	decompressed, err := io.ReadAll(zlibReader)
	if err != nil {
		return err
	}

	// Apply the delta encoding used in the original implementation
	for i := range decompressed {
		if (i & 63) != 0 {
			decompressed[i] = (decompressed[i-1]) + (decompressed[i] - 127)
		} else {
			decompressed[i] = 127 + (decompressed[i] - 127)
		}
	}

	// Pad with 4161 zero bytes as in the original implementation
	padding := make([]byte, 4161)
	decompressed = append(decompressed, padding...)

	// Save to the global LUT
	MixboxLUT = decompressed

	return nil
}