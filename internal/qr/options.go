package qr

import (
	"image/color"

	qrcode "github.com/skip2/go-qrcode"
)

// Options configures QR code generation.
type Options struct {
	Size            int
	Level           qrcode.RecoveryLevel
	ForegroundColor color.Color
	BackgroundColor color.Color
	BorderSize      int
	LogoPath        string
	LogoScale       float64
}

// DefaultOptions returns sensible defaults for QR code generation.
func DefaultOptions() Options {
	return Options{
		Size:            256,
		Level:           qrcode.Medium,
		ForegroundColor: color.Black,
		BackgroundColor: color.White,
		BorderSize:      4,
		LogoScale:       0.2,
	}
}
