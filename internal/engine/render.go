package engine

import (
	"bytes"
	"image"

	_ "image/jpeg"
	_ "image/png"

	"github.com/dolmen-go/kittyimg"
	"golang.org/x/image/draw"
)

func RenderImage(imageData []byte, maxCellWidth, maxCellHeight int) (string, int, error) {
	if maxCellWidth <= 0 || maxCellHeight <= 0 {
		return "", 0, nil
	}

	// Decode the image
	src, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		return "", 0, err
	}
	srcBounds := src.Bounds()

	// cell aspect ratio is mostly 0.5 (width/height)
	const cellAspectRatio = 0.5

	imageRatio := float64(srcBounds.Dx()) / float64(srcBounds.Dy())

	wFromH := float64(maxCellHeight) / cellAspectRatio * imageRatio
	hFromW := float64(maxCellWidth) * cellAspectRatio / imageRatio

	var targetCols, targetRows int

	if wFromH <= float64(maxCellWidth) {
		targetCols = int(wFromH)
		targetRows = maxCellHeight
	} else {
		targetCols = maxCellWidth
		targetRows = int(hFromW)
	}

	if targetCols == 0 {
		targetCols = 1
	}
	if targetRows == 0 {
		targetRows = 1
	}

	// We map 1 cell = ~10x20 pixels roughly for quality.
	pixelWidth := targetCols * 10
	pixelHeight := targetRows * 20

	dst := image.NewRGBA(image.Rect(0, 0, pixelWidth, pixelHeight))
	draw.BiLinear.Scale(dst, dst.Bounds(), src, srcBounds, draw.Over, nil)

	// 4. Encode to Kitty Protocol
	var buf bytes.Buffer
	if err := kittyimg.Fprint(&buf, dst); err != nil {
		return "", 0, err
	}

	return buf.String(), targetRows, nil
}

func ClearAllImages() string {
	return "\x1b_Ga=d\x1b\\"
}
