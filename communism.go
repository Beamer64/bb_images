package bb_images

import (
	_ "embed"
	"image"

	"github.com/Beamer64/bb_images/internal/draw"
	"github.com/disintegration/imaging"
)

//go:embed res/overlays/communism.png
var communismBytes []byte

// Communism returns src with a translucent hammer-and-sickle/red overlay (0.5 opacity) as PNG bytes.
func Communism(src image.Image) ([]byte, error) {
	overlay, err := draw.Decode(communismBytes)
	if err != nil {
		return nil, err
	}
	b := src.Bounds()
	resized := imaging.Fill(overlay, b.Dx(), b.Dy(), imaging.Center, imaging.Lanczos)
	out := imaging.Overlay(src, resized, image.Point{}, 0.5)
	return draw.EncodePNG(out)
}
