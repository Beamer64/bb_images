package bb_images

import (
	"image"

	"github.com/Beamer64/bb_images/internal/draw"
)

// Ground tints src with a warm soil-brown palette.
func Ground(src image.Image) ([]byte, error) {
	return draw.EncodePNG(tinted(src, 160, 100, 60))
}
