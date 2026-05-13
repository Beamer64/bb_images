package bb_images

import (
	"image"

	"github.com/Beamer64/bb_images/internal/draw"
	"github.com/disintegration/imaging"
)

// Flip returns src flipped vertically (upside-down) as PNG bytes.
func Flip(src image.Image) ([]byte, error) {
	return draw.EncodePNG(imaging.FlipV(src))
}
