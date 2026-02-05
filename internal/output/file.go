package output

import (
	"errors"
	"os"
	"path/filepath"
)

// WriteFile writes bytes to the requested path, creating parent directories as needed.
func WriteFile(path string, data []byte) error {
	if path == "" {
		return errors.New("output path is empty")
	}

	dir := filepath.Dir(path)
	if dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}

	return os.WriteFile(path, data, 0o644)
}
