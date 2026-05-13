package bb_images

import (
	"image"

	"github.com/Beamer64/bb_images/internal/draw"
	"github.com/disintegration/imaging"
)

// Deepfry oversaturates, over-contrasts, and sharpens src for the classic meme look.
func Deepfry(src image.Image) ([]byte, error) {
	out := imaging.AdjustSaturation(src, 80)
	out = imaging.AdjustContrast(out, 50)
	return draw.EncodePNG(imaging.Sharpen(out, 5))
}
