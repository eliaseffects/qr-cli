package qr

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"strings"
)

// PNG renders a QR code into PNG bytes.
func PNG(data string, opts Options) ([]byte, error) {
	bitmap, err := Bitmap(data, opts)
	if err != nil {
		return nil, err
	}

	if opts.Size <= 0 {
		opts.Size = DefaultOptions().Size
	}

	totalModules := len(bitmap)
	if totalModules == 0 {
		return nil, fmt.Errorf("empty QR bitmap")
	}

	size := opts.Size
	if size < totalModules {
		size = totalModules
	}

	scale := size / totalModules
	if scale < 1 {
		scale = 1
	}
	actual := totalModules * scale
	pad := (size - actual) / 2

	img := image.NewRGBA(image.Rect(0, 0, size, size))
	bg := colorToRGBA(opts.BackgroundColor)
	fg := colorToRGBA(opts.ForegroundColor)

	draw.Draw(img, img.Bounds(), &image.Uniform{C: bg}, image.Point{}, draw.Src)

	for y := 0; y < totalModules; y++ {
		for x := 0; x < totalModules; x++ {
			if !bitmap[y][x] {
				continue
			}
			startX := pad + x*scale
			startY := pad + y*scale
			rect := image.Rect(startX, startY, startX+scale, startY+scale)
			draw.Draw(img, rect, &image.Uniform{C: fg}, image.Point{}, draw.Src)
		}
	}

	if opts.LogoPath != "" {
		if err := overlayLogoPNG(img, opts, size); err != nil {
			return nil, err
		}
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// SVG renders a QR code into SVG bytes.
func SVG(data string, opts Options) ([]byte, error) {
	bitmap, err := Bitmap(data, opts)
	if err != nil {
		return nil, err
	}

	if opts.Size <= 0 {
		opts.Size = DefaultOptions().Size
	}

	totalModules := len(bitmap)
	if totalModules == 0 {
		return nil, fmt.Errorf("empty QR bitmap")
	}

	fg := colorToHex(opts.ForegroundColor)
	bg := colorToHex(opts.BackgroundColor)

	var b strings.Builder
	b.WriteString(fmt.Sprintf(
		`<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d" shape-rendering="crispEdges">`,
		opts.Size, opts.Size, totalModules, totalModules,
	))
	b.WriteString(fmt.Sprintf(`<rect width="100%%" height="100%%" fill="%s"/>`, bg))
	b.WriteString(fmt.Sprintf(`<g fill="%s">`, fg))

	for y := 0; y < totalModules; y++ {
		for x := 0; x < totalModules; x++ {
			if bitmap[y][x] {
				b.WriteString(fmt.Sprintf(`<rect x="%d" y="%d" width="1" height="1"/>`, x, y))
			}
		}
	}

	b.WriteString(`</g>`)

	if opts.LogoPath != "" {
		element, err := svgLogoElement(opts, totalModules)
		if err != nil {
			return nil, err
		}
		b.WriteString(element)
	}

	b.WriteString(`</svg>`)
	return []byte(b.String()), nil
}

func colorToRGBA(c color.Color) color.RGBA {
	return color.RGBAModel.Convert(c).(color.RGBA)
}

func colorToHex(c color.Color) string {
	rgba := colorToRGBA(c)
	return fmt.Sprintf("#%02x%02x%02x", rgba.R, rgba.G, rgba.B)
}
