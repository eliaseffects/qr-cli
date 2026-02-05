package qr

import (
	"fmt"
	"image"
	"os"

	"github.com/liyue201/goqr"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

// DecodeImage extracts QR payloads from an image.
func DecodeImage(img image.Image) ([]string, error) {
	codes, err := goqr.Recognize(img)
	if err != nil {
		return nil, err
	}

	results := make([]string, 0, len(codes))
	for _, code := range codes {
		if code == nil {
			continue
		}
		results = append(results, string(code.Payload))
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("no QR code data found")
	}

	return results, nil
}

// DecodeFile loads an image file and extracts QR payloads.
func DecodeFile(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	return DecodeImage(img)
}
