package qr

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"math"
	"os"

	xdraw "golang.org/x/image/draw"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

const (
	defaultLogoScale = 0.2
	minLogoScale     = 0.05
	maxLogoScale     = 0.4
	logoPaddingScale = 0.12
)

func clampLogoScale(scale float64) float64 {
	if scale <= 0 {
		return defaultLogoScale
	}
	if scale < minLogoScale {
		return minLogoScale
	}
	if scale > maxLogoScale {
		return maxLogoScale
	}
	return scale
}

func loadLogoImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func scaledLogoPNG(path string, size int) ([]byte, error) {
	if size <= 0 {
		return nil, fmt.Errorf("invalid logo size")
	}

	img, err := loadLogoImage(path)
	if err != nil {
		return nil, err
	}

	scaled := scaleToSquare(img, size)

	var buf bytes.Buffer
	if err := png.Encode(&buf, scaled); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func scaleToSquare(src image.Image, size int) *image.RGBA {
	dst := image.NewRGBA(image.Rect(0, 0, size, size))

	bounds := src.Bounds()
	srcW := bounds.Dx()
	srcH := bounds.Dy()
	if srcW == 0 || srcH == 0 {
		return dst
	}

	scale := math.Min(float64(size)/float64(srcW), float64(size)/float64(srcH))
	newW := int(math.Round(float64(srcW) * scale))
	newH := int(math.Round(float64(srcH) * scale))
	if newW < 1 {
		newW = 1
	}
	if newH < 1 {
		newH = 1
	}

	tmp := image.NewRGBA(image.Rect(0, 0, newW, newH))
	xdraw.CatmullRom.Scale(tmp, tmp.Bounds(), src, bounds, draw.Over, nil)

	offsetX := (size - newW) / 2
	offsetY := (size - newH) / 2
	draw.Draw(dst, image.Rect(offsetX, offsetY, offsetX+newW, offsetY+newH), tmp, image.Point{}, draw.Over)

	return dst
}

func overlayLogoPNG(img *image.RGBA, opts Options, canvasSize int) error {
	if opts.LogoPath == "" {
		return nil
	}

	scale := clampLogoScale(opts.LogoScale)
	logoSize := int(math.Round(float64(canvasSize) * scale))
	if logoSize < 1 {
		return nil
	}

	padding := int(math.Round(float64(logoSize) * logoPaddingScale))
	bgSize := logoSize + padding*2
	if bgSize > canvasSize {
		bgSize = canvasSize
	}

	bgX := (canvasSize - bgSize) / 2
	bgY := (canvasSize - bgSize) / 2
	draw.Draw(img, image.Rect(bgX, bgY, bgX+bgSize, bgY+bgSize), &image.Uniform{C: opts.BackgroundColor}, image.Point{}, draw.Src)

	logo, err := loadLogoImage(opts.LogoPath)
	if err != nil {
		return err
	}

	scaled := scaleToSquare(logo, logoSize)
	logoX := (canvasSize - logoSize) / 2
	logoY := (canvasSize - logoSize) / 2
	draw.Draw(img, image.Rect(logoX, logoY, logoX+logoSize, logoY+logoSize), scaled, image.Point{}, draw.Over)

	return nil
}

func svgLogoElement(opts Options, totalModules int) (string, error) {
	if opts.LogoPath == "" {
		return "", nil
	}

	scale := clampLogoScale(opts.LogoScale)
	logoUnits := float64(totalModules) * scale
	paddingUnits := logoUnits * logoPaddingScale
	bgUnits := logoUnits + paddingUnits*2

	bgX := (float64(totalModules) - bgUnits) / 2
	bgY := (float64(totalModules) - bgUnits) / 2

	logoX := (float64(totalModules) - logoUnits) / 2
	logoY := (float64(totalModules) - logoUnits) / 2

	pixelSize := int(math.Round(float64(opts.Size) * scale))
	if pixelSize < 1 {
		pixelSize = 1
	}

	pngData, err := scaledLogoPNG(opts.LogoPath, pixelSize)
	if err != nil {
		return "", err
	}

	encoded := base64.StdEncoding.EncodeToString(pngData)

	return fmt.Sprintf(
		`<rect x="%.2f" y="%.2f" width="%.2f" height="%.2f" fill="%s"/>`+
			`<image href="data:image/png;base64,%s" x="%.2f" y="%.2f" width="%.2f" height="%.2f"/>`,
		bgX, bgY, bgUnits, bgUnits, colorToHex(opts.BackgroundColor),
		encoded, logoX, logoY, logoUnits, logoUnits,
	), nil
}
