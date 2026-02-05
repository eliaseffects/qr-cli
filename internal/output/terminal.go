package output

import (
	"fmt"
	"image/color"
	"strings"

	"github.com/eliaseffects/qr-cli/internal/qr"
)

// TerminalOptions controls ANSI color output and inversion.
type TerminalOptions struct {
	UseColor bool
	Invert   bool
}

// ToTerminal renders QR code as Unicode block characters.
func ToTerminal(data string, opts qr.Options, term TerminalOptions) (string, error) {
	bitmap, err := qr.Bitmap(data, opts)
	if err != nil {
		return "", err
	}

	var sb strings.Builder
	height := len(bitmap)
	if height == 0 {
		return "", nil
	}
	width := len(bitmap[0])

	if term.UseColor {
		fg := toRGB(opts.ForegroundColor)
		bg := toRGB(opts.BackgroundColor)
		lastFg := ""
		lastBg := ""

		for y := 0; y < height; y += 2 {
			for x := 0; x < width; x++ {
				upper := bitmap[y][x]
				lower := false
				if y+1 < height {
					lower = bitmap[y+1][x]
				}
				if term.Invert {
					upper = !upper
					lower = !lower
				}

				char, fgColor, bgColor := blockFor(upper, lower, fg, bg)
				if fgColor != lastFg {
					sb.WriteString(fgColor)
					lastFg = fgColor
				}
				if bgColor != lastBg {
					sb.WriteString(bgColor)
					lastBg = bgColor
				}
				sb.WriteRune(char)
			}
			sb.WriteString("\x1b[0m")
			sb.WriteRune('\n')
			lastFg = ""
			lastBg = ""
		}
	} else {
		for y := 0; y < height; y += 2 {
			for x := 0; x < width; x++ {
				upper := bitmap[y][x]
				lower := false
				if y+1 < height {
					lower = bitmap[y+1][x]
				}
				if term.Invert {
					upper = !upper
					lower = !lower
				}

				switch {
				case upper && lower:
					sb.WriteRune('█')
				case upper && !lower:
					sb.WriteRune('▀')
				case !upper && lower:
					sb.WriteRune('▄')
				default:
					sb.WriteRune(' ')
				}
			}
			sb.WriteRune('\n')
		}
	}

	return sb.String(), nil
}

func toRGB(c color.Color) [3]int {
	rgba := color.RGBAModel.Convert(c).(color.RGBA)
	return [3]int{int(rgba.R), int(rgba.G), int(rgba.B)}
}

func ansiFg(rgb [3]int) string {
	return fmt.Sprintf("\x1b[38;2;%d;%d;%dm", rgb[0], rgb[1], rgb[2])
}

func ansiBg(rgb [3]int) string {
	return fmt.Sprintf("\x1b[48;2;%d;%d;%dm", rgb[0], rgb[1], rgb[2])
}

func blockFor(upper, lower bool, fg, bg [3]int) (rune, string, string) {
	switch {
	case upper && lower:
		return '█', ansiFg(fg), ansiBg(bg)
	case upper && !lower:
		return '▀', ansiFg(fg), ansiBg(bg)
	case !upper && lower:
		return '▄', ansiFg(fg), ansiBg(bg)
	default:
		return ' ', ansiFg(bg), ansiBg(bg)
	}
}
