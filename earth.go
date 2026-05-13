package bb_images

import (
	"image"

	"github.com/Beamer64/bb_images/internal/draw"
)

// Earth tints src with a green-blue earth palette (vegetation/water tones).
func Earth(src image.Image) ([]byte, error) {
	return draw.EncodePNG(tinted(src, 80, 180, 100))
}
