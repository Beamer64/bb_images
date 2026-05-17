package color

import (
	"image"
	imgcolor "image/color"

	"github.com/disintegration/imaging"
)

// tinted maps each pixel of src to a monochromatic shade of the target RGB
// derived from the pixel's luminance.
func tinted(src image.Image, tr, tg, tb float64) image.Image {
	return imaging.AdjustFunc(src, func(c imgcolor.NRGBA) imgcolor.NRGBA {
		l := (0.299*float64(c.R) + 0.587*float64(c.G) + 0.114*float64(c.B)) / 255.0
		return imgcolor.NRGBA{
			R: clamp8(tr * l),
			G: clamp8(tg * l),
			B: clamp8(tb * l),
			A: c.A,
		}
	})
}

func clamp8(v float64) uint8 {
	if v < 0 {
		return 0
	}
	if v > 255 {
		return 255
	}
	return uint8(v)
}
