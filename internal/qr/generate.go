package qr

import (
	"errors"
	"strings"

	qrcode "github.com/skip2/go-qrcode"
)

// Generate creates a QRCode instance from input data.
func Generate(data string, opts Options) (*qrcode.QRCode, error) {
	if strings.TrimSpace(data) == "" {
		return nil, errors.New("data is empty")
	}

	code, err := qrcode.New(data, opts.Level)
	if err != nil {
		return nil, err
	}

	// We'll handle borders ourselves for consistent sizing.
	code.DisableBorder = true
	code.ForegroundColor = opts.ForegroundColor
	code.BackgroundColor = opts.BackgroundColor

	return code, nil
}

// Bitmap returns a QR bitmap with the configured border size applied.
func Bitmap(data string, opts Options) ([][]bool, error) {
	if opts.BorderSize < 0 {
		opts.BorderSize = 0
	}

	code, err := Generate(data, opts)
	if err != nil {
		return nil, err
	}

	bitmap := code.Bitmap()
	if opts.BorderSize == 0 {
		return bitmap, nil
	}

	return addBorder(bitmap, opts.BorderSize), nil
}

func addBorder(bitmap [][]bool, border int) [][]bool {
	if len(bitmap) == 0 || border == 0 {
		return bitmap
	}

	width := len(bitmap[0])
	total := width + border*2

	bordered := make([][]bool, total)
	for y := 0; y < total; y++ {
		row := make([]bool, total)
		if y >= border && y < border+len(bitmap) {
			copy(row[border:border+width], bitmap[y-border])
		}
		bordered[y] = row
	}

	return bordered
}
