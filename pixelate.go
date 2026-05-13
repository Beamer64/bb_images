package bb_images

import (
	"errors"
	"image"

	"github.com/Beamer64/bb_images/internal/draw"
	"github.com/disintegration/imaging"
)

// Pixelate returns src pixelated by factor (e.g. 8 for an 8x8 mosaic) as PNG bytes.
func Pixelate(src image.Image, factor int) ([]byte, error) {
	if factor < 2 {
		return nil, errors.New("pixelate: factor must be >= 2")
	}
	b := src.Bounds()
	w, h := b.Dx(), b.Dy()
	if w < factor || h < factor {
		return nil, errors.New("pixelate: image smaller than factor")
	}
	small := imaging.Resize(src, w/factor, h/factor, imaging.NearestNeighbor)
	out := imaging.Resize(small, w, h, imaging.NearestNeighbor)
	return draw.EncodePNG(out)
}
