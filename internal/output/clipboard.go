package output

import (
	"fmt"

	"golang.design/x/clipboard"
)

// CopyPNG writes PNG bytes to the system clipboard.
func CopyPNG(pngData []byte) error {
	if err := clipboard.Init(); err != nil {
		return fmt.Errorf("clipboard init failed: %w", err)
	}
	clipboard.Write(clipboard.FmtImage, pngData)
	return nil
}

// CopyText writes text to the system clipboard.
func CopyText(text string) error {
	if err := clipboard.Init(); err != nil {
		return fmt.Errorf("clipboard init failed: %w", err)
	}
	clipboard.Write(clipboard.FmtText, []byte(text))
	return nil
}
