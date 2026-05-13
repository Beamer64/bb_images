package bb_images

import (
	_ "embed"
	"image"
	"image/color"

	"github.com/Beamer64/bb_images/internal/draw"
	"github.com/disintegration/imaging"
)

//go:embed res/overlays/wasted.png
var wastedBytes []byte

// Wasted desaturates src to grayscale, then composites a dark-grey backdrop
// and the still-colored WASTED overlay on top — source goes monochrome while
// the title text keeps its colors.
func Wasted(src image.Image) ([]byte, error) {
	overlay, err := draw.Decode(wastedBytes)
	if err != nil {
		return nil, err
	}
	gray := imaging.Grayscale(src)
	b := src.Bounds()
	w, h := b.Dx(), b.Dy()

	fitted := imaging.Fit(overlay, w, h, imaging.Lanczos)
	fw, fh := fitted.Bounds().Dx(), fitted.Bounds().Dy()
	offX, offY := (w-fw)/2, (h-fh)/2

	barH := fh * 4 / 3
	barOffY := offY - (barH-fh)/2
	if barOffY < 0 {
		barH += barOffY
		barOffY = 0
	}
	if barOffY+barH > h {
		barH = h - barOffY
	}

	bar := imaging.New(fw, barH, color.RGBA{R: 20, G: 20, B: 20, A: 255})
	out := imaging.Overlay(gray, bar, image.Point{X: offX, Y: barOffY}, 0.65)
	out = imaging.Overlay(out, fitted, image.Point{X: offX, Y: offY}, 1.0)
	return draw.EncodePNG(out)
}
