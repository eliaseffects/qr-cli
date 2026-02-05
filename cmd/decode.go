package cmd

import (
	"errors"
	"fmt"

	"github.com/eliaseffects/qr-cli/internal/qr"
	"github.com/spf13/cobra"
)

var (
	decodeFile string

	decodeCmd = &cobra.Command{
		Use:   "decode <image>",
		Short: "Decode QR code(s) from an image",
		Args:  cobra.MaximumNArgs(1),
		RunE:  runDecode,
	}
)

func init() {
	decodeCmd.Flags().StringVarP(&decodeFile, "file", "f", "", "Image file to decode")
}

func runDecode(cmd *cobra.Command, args []string) error {
	path := decodeFile
	if len(args) > 0 {
		path = args[0]
	}
	if path == "" {
		return errors.New("image file is required")
	}

	results, err := qr.DecodeFile(path)
	if err != nil {
		return err
	}

	out := cmd.OutOrStdout()
	for _, value := range results {
		fmt.Fprintln(out, value)
	}

	return nil
}
