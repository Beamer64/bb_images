package color

import (
	"image"

	"github.com/Beamer64/bb_images/internal/draw"
)

// Freeze tints src with a pale ice-blue palette.
func Freeze(src image.Image) ([]byte, error) {
	return draw.EncodePNG(tinted(src, 150, 200, 230))
}
