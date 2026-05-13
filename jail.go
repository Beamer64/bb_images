package bb_images

import (
	_ "embed"
	"image"

	"github.com/Beamer64/bb_images/internal/draw"
	"github.com/disintegration/imaging"
)

//go:embed res/overlays/jail.png
var jailBytes []byte

// Jail composites the prison-bars overlay over src at full opacity, fitting
// the bars to cover the whole source while preserving the overlay's aspect.
func Jail(src image.Image) ([]byte, error) {
	overlay, err := draw.Decode(jailBytes)
	if err != nil {
		return nil, err
	}
	b := src.Bounds()
	resized := imaging.Fill(overlay, b.Dx(), b.Dy(), imaging.Center, imaging.Lanczos)
	out := imaging.Overlay(src, resized, image.Point{}, 1.0)
	return draw.EncodePNG(out)
}
