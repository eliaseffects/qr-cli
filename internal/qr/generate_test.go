package qr_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/eliaseffects/qr-cli/internal/qr"
)

func TestGenerate(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		wantErr bool
	}{
		{"simple URL", "https://example.com", false},
		{"unicode text", "こんにちは", false},
		{"empty string", "", true},
		{"long text", strings.Repeat("a", 2000), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := qr.DefaultOptions()
			_, err := qr.Generate(tt.data, opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("Generate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPNG(t *testing.T) {
	opts := qr.DefaultOptions()
	pngData, err := qr.PNG("https://example.com", opts)
	if err != nil {
		t.Fatalf("PNG() error = %v", err)
	}

	if len(pngData) < 8 {
		t.Fatal("PNG too short")
	}

	magic := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
	for i, b := range magic {
		if pngData[i] != b {
			t.Fatalf("invalid PNG magic byte at %d", i)
		}
	}
}

func TestPNGWithLogo(t *testing.T) {
	logoPath := filepath.Join("..", "..", "testdata", "logo.png")
	if _, err := os.Stat(logoPath); err != nil {
		t.Fatalf("logo test asset missing: %v", err)
	}

	opts := qr.DefaultOptions()
	opts.LogoPath = logoPath
	pngData, err := qr.PNG("https://example.com", opts)
	if err != nil {
		t.Fatalf("PNG() with logo error = %v", err)
	}
	if len(pngData) == 0 {
		t.Fatal("PNG output is empty")
	}
}

func TestSVGWithLogo(t *testing.T) {
	logoPath := filepath.Join("..", "..", "testdata", "logo.png")
	if _, err := os.Stat(logoPath); err != nil {
		t.Fatalf("logo test asset missing: %v", err)
	}

	opts := qr.DefaultOptions()
	opts.LogoPath = logoPath
	svgData, err := qr.SVG("https://example.com", opts)
	if err != nil {
		t.Fatalf("SVG() with logo error = %v", err)
	}
	if !strings.Contains(string(svgData), "data:image/png;base64") {
		t.Fatal("SVG output missing embedded logo data")
	}
}
