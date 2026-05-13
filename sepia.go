package bb_images

import (
	"image"
	"image/color"

	"github.com/Beamer64/bb_images/internal/draw"
	"github.com/disintegration/imaging"
)

// Sepia applies the classic warm-tan sepia transform to src.
func Sepia(src image.Image) ([]byte, error) {
	out := imaging.AdjustFunc(src, func(c color.NRGBA) color.NRGBA {
		r, g, b := float64(c.R), float64(c.G), float64(c.B)
		return color.NRGBA{
			R: clamp8(0.393*r + 0.769*g + 0.189*b),
			G: clamp8(0.349*r + 0.686*g + 0.168*b),
			B: clamp8(0.272*r + 0.534*g + 0.131*b),
			A: c.A,
		}
	})
	return draw.EncodePNG(out)
}
