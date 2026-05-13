package bb_images

import (
	_ "embed"
	"image"

	"github.com/Beamer64/bb_images/internal/draw"
	"github.com/disintegration/imaging"
)

//go:embed res/overlays/USflag.png
var americaFlagBytes []byte

// America returns src with a translucent American flag overlay (0.5 opacity) as PNG bytes.
func America(src image.Image) ([]byte, error) {
	flag, err := draw.Decode(americaFlagBytes)
	if err != nil {
		return nil, err
	}
	b := src.Bounds()
	resized := imaging.Fill(flag, b.Dx(), b.Dy(), imaging.Center, imaging.Lanczos)
	out := imaging.Overlay(src, resized, image.Point{}, 0.4)
	return draw.EncodePNG(out)
}
