package bb_images

import (
	"image"

	"github.com/Beamer64/bb_images/internal/draw"
	"github.com/disintegration/imaging"
)

// Invert returns src with all RGB channels negated as PNG bytes.
func Invert(src image.Image) ([]byte, error) {
	return draw.EncodePNG(imaging.Invert(src))
}
