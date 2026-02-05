package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/eliaseffects/qr-cli/internal/output"
	"github.com/eliaseffects/qr-cli/internal/qr"
	"github.com/spf13/cobra"
)

type batchFlags struct {
	File   string
	Dir    string
	Size   int
	Format string
	Prefix string
	Quiet  bool
}

var (
	batchCfg = batchFlags{}

	batchCmd = &cobra.Command{
		Use:   "batch",
		Short: "Generate QR codes from a file of data (one per line)",
		RunE:  runBatch,
	}
)

func init() {
	batchCmd.Flags().StringVarP(&batchCfg.File, "file", "f", "", "Input file with data (one per line, required)")
	batchCmd.Flags().StringVarP(&batchCfg.Dir, "dir", "d", "./qr-output", "Output directory")
	batchCmd.Flags().IntVarP(&batchCfg.Size, "size", "s", 256, "Image size in pixels")
	batchCmd.Flags().StringVar(&batchCfg.Format, "format", "png", "Output format: png, svg")
	batchCmd.Flags().StringVar(&batchCfg.Prefix, "prefix", "qr-", "Filename prefix")
	batchCmd.Flags().BoolVarP(&batchCfg.Quiet, "quiet", "q", false, "Suppress non-error output")
	_ = batchCmd.MarkFlagRequired("file")

	bindBatchFlags(batchCmd)
}

func runBatch(cmd *cobra.Command, args []string) error {
	applyBatchConfig(cmd)

	if strings.TrimSpace(batchCfg.File) == "" {
		return errors.New("input file is required")
	}

	format := strings.ToLower(strings.TrimSpace(batchCfg.Format))
	if format == "" {
		format = "png"
	}
	if format != "png" && format != "svg" {
		return fmt.Errorf("unsupported format: %s", batchCfg.Format)
	}

	if batchCfg.Size <= 0 {
		return errors.New("size must be greater than zero")
	}

	file, err := os.Open(batchCfg.File)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	if len(lines) == 0 {
		return errors.New("no data found in input file")
	}

	if err := os.MkdirAll(batchCfg.Dir, 0o755); err != nil {
		return err
	}

	opts := qr.DefaultOptions()
	opts.Size = batchCfg.Size

	for i, line := range lines {
		var payload []byte
		if format == "svg" {
			payload, err = qr.SVG(line, opts)
		} else {
			payload, err = qr.PNG(line, opts)
		}
		if err != nil {
			return fmt.Errorf("line %d: %w", i+1, err)
		}

		filename := fmt.Sprintf("%s%03d.%s", batchCfg.Prefix, i+1, format)
		path := filepath.Join(batchCfg.Dir, filename)
		if err := output.WriteFile(path, payload); err != nil {
			return err
		}
	}

	if !batchCfg.Quiet {
		fmt.Printf("✓ Generated %d QR codes in %s\n", len(lines), batchCfg.Dir)
	}

	return nil
}
