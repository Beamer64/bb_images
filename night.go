package bb_images

import (
	"image"

	"github.com/Beamer64/bb_images/internal/draw"
)

// Night tints src with a deep navy palette.
func Night(src image.Image) ([]byte, error) {
	return draw.EncodePNG(tinted(src, 60, 80, 140))
}
