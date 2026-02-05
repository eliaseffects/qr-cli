package main

import (
	"image"
	_ "image/png"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestCLIEndToEnd(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	tmpDir := t.TempDir()
	outPath := filepath.Join(tmpDir, "test.png")

	cmd := exec.Command("go", "run", ".", "https://example.com", "-o", outPath)
	if err := cmd.Run(); err != nil {
		t.Fatalf("CLI failed: %v", err)
	}

	file, err := os.Open(outPath)
	if err != nil {
		t.Fatalf("failed to open output file: %v", err)
	}
	defer file.Close()

	_, format, err := image.Decode(file)
	if err != nil {
		t.Fatalf("failed to decode output image: %v", err)
	}
	if format != "png" {
		t.Fatalf("expected png output, got %s", format)
	}
}
