package bb_images

import (
	"image"

	"github.com/Beamer64/bb_images/internal/draw"
	"github.com/disintegration/imaging"
)

const blurSigma = 5.0

// Blur applies a gaussian blur to src and returns PNG bytes.
func Blur(src image.Image) ([]byte, error) {
	return draw.EncodePNG(imaging.Blur(src, blurSigma))
}
