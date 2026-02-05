package qr_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/eliaseffects/qr-cli/internal/qr"
)

func TestDecodeFile(t *testing.T) {
	payload := "https://example.com"
	opts := qr.DefaultOptions()
	opts.Size = 512

	pngData, err := qr.PNG(payload, opts)
	if err != nil {
		t.Fatalf("PNG() error = %v", err)
	}

	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "code.png")
	if err := os.WriteFile(path, pngData, 0o644); err != nil {
		t.Fatalf("failed to write test image: %v", err)
	}

	results, err := qr.DecodeFile(path)
	if err != nil {
		t.Fatalf("DecodeFile() error = %v", err)
	}

	if len(results) == 0 {
		t.Fatal("expected decoded payloads")
	}
	if results[0] != payload {
		t.Fatalf("expected %q, got %q", payload, results[0])
	}
}
