package cmd

import (
	"errors"
	"fmt"
	"image/color"
	"path/filepath"
	"strings"

	"github.com/eliaseffects/qr-cli/internal/output"
	"github.com/eliaseffects/qr-cli/internal/qr"
	qrcode "github.com/skip2/go-qrcode"
	"github.com/spf13/cobra"
)

type OutputFlags struct {
	OutputPath string
	Size       int
	Format     string
	Level      string
	FgColor    string
	BgColor    string
	Border     int
	LogoPath   string
	LogoScale  float64
	Invert     bool
	TermColor  bool
	Terminal   bool
	OpenViewer bool
	CopyClip   bool
	Quiet      bool
}

func addOutputFlags(cmd *cobra.Command, flags *OutputFlags, includeTerminal bool) {
	cmd.Flags().StringVarP(&flags.OutputPath, "output", "o", "", "Output file path")
	cmd.Flags().IntVarP(&flags.Size, "size", "s", 256, "Image size in pixels")
	cmd.Flags().StringVarP(&flags.Format, "format", "f", "png", "Output format: png, svg, terminal")
	cmd.Flags().StringVarP(&flags.Level, "level", "l", "M", "Error correction: L, M, Q, H")
	cmd.Flags().StringVar(&flags.FgColor, "fg", "#000000", "Foreground color (hex)")
	cmd.Flags().StringVar(&flags.BgColor, "bg", "#ffffff", "Background color (hex)")
	cmd.Flags().IntVar(&flags.Border, "border", 4, "Border size in modules")
	cmd.Flags().StringVar(&flags.LogoPath, "logo", "", "Path to logo image to overlay")
	cmd.Flags().Float64Var(&flags.LogoScale, "logo-scale", 0.2, "Logo size as fraction of QR (0.05-0.4)")
	cmd.Flags().BoolVar(&flags.Invert, "invert", false, "Invert terminal rendering colors")
	cmd.Flags().BoolVar(&flags.TermColor, "terminal-color", false, "Use ANSI colors when rendering in terminal")
	if includeTerminal {
		cmd.Flags().BoolVarP(&flags.Terminal, "terminal", "t", false, "Render in terminal")
	}
	cmd.Flags().BoolVar(&flags.OpenViewer, "open", false, "Open in system viewer")
	cmd.Flags().BoolVar(&flags.CopyClip, "copy", false, "Copy to clipboard")
	cmd.Flags().BoolVarP(&flags.Quiet, "quiet", "q", false, "Suppress non-error output")
}

func runGenerate(data string, flags OutputFlags, formatSet bool) error {
	if strings.TrimSpace(data) == "" {
		return errors.New("no data provided")
	}

	opts, err := flags.toOptions()
	if err != nil {
		return err
	}

	format := strings.ToLower(flags.Format)
	if !formatSet && flags.OutputPath != "" {
		switch strings.ToLower(filepath.Ext(flags.OutputPath)) {
		case ".svg":
			format = "svg"
		case ".png":
			format = "png"
		}
	}

	if !flags.Terminal && format != "terminal" && (flags.Invert || flags.TermColor) {
		return errors.New("terminal color/invert requires --terminal or --format terminal")
	}

	if flags.Terminal || format == "terminal" {
		if flags.LogoPath != "" {
			return errors.New("logo overlay is not supported for terminal rendering")
		}
		if flags.CopyClip {
			return errors.New("clipboard output is not supported for terminal rendering")
		}
		result, err := output.ToTerminal(data, opts, output.TerminalOptions{
			UseColor: flags.TermColor,
			Invert:   flags.Invert,
		})
		if err != nil {
			return err
		}
		fmt.Print(result)
		return nil
	}

	switch format {
	case "png", "svg":
	default:
		return fmt.Errorf("unsupported format: %s", flags.Format)
	}

	outPath := flags.OutputPath
	if outPath == "" {
		if format == "svg" {
			outPath = "qr.svg"
		} else {
			outPath = "qr.png"
		}
	}

	var payload []byte
	if format == "svg" {
		payload, err = qr.SVG(data, opts)
	} else {
		payload, err = qr.PNG(data, opts)
	}
	if err != nil {
		return err
	}

	if err := output.WriteFile(outPath, payload); err != nil {
		return err
	}

	if !flags.Quiet {
		fmt.Printf("✓ QR code saved to %s\n", outPath)
	}

	if flags.CopyClip {
		if format == "svg" {
			if err := output.CopyText(string(payload)); err != nil {
				return err
			}
		} else {
			if err := output.CopyPNG(payload); err != nil {
				return err
			}
		}
	}

	if flags.OpenViewer {
		if err := output.OpenInViewer(outPath); err != nil {
			return fmt.Errorf("failed to open viewer: %w", err)
		}
	}

	return nil
}

func (flags OutputFlags) toOptions() (qr.Options, error) {
	opts := qr.DefaultOptions()
	if flags.Size <= 0 {
		return opts, errors.New("size must be greater than zero")
	}
	if flags.Border < 0 {
		return opts, errors.New("border size must be zero or positive")
	}

	opts.Size = flags.Size
	opts.Level = parseLevel(flags.Level)

	fg, err := parseColor(flags.FgColor)
	if err != nil {
		return opts, err
	}
	bg, err := parseColor(flags.BgColor)
	if err != nil {
		return opts, err
	}

	opts.ForegroundColor = fg
	opts.BackgroundColor = bg
	opts.BorderSize = flags.Border
	opts.LogoPath = strings.TrimSpace(flags.LogoPath)
	opts.LogoScale = flags.LogoScale

	return opts, nil
}

func parseLevel(s string) qrcode.RecoveryLevel {
	switch strings.ToUpper(s) {
	case "L":
		return qrcode.Low
	case "M":
		return qrcode.Medium
	case "Q":
		return qrcode.High
	case "H":
		return qrcode.Highest
	default:
		return qrcode.Medium
	}
}

func parseColor(hex string) (color.Color, error) {
	hex = strings.TrimSpace(strings.TrimPrefix(hex, "#"))
	if len(hex) != 6 {
		return nil, fmt.Errorf("invalid color: %s", hex)
	}

	var r, g, b uint8
	if _, err := fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b); err != nil {
		return nil, fmt.Errorf("invalid color: %s", hex)
	}

	return color.RGBA{R: r, G: g, B: b, A: 255}, nil
}
