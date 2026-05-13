package bb_images

import (
	"image"
	"image/color"

	"github.com/Beamer64/bb_images/internal/draw"
	"github.com/disintegration/imaging"
)

const posterizeLevels = 4

// Posterize reduces src to posterizeLevels values per RGB channel, evenly spread across 0–255.
func Posterize(src image.Image) ([]byte, error) {
	const n = posterizeLevels
	bandSize := 256 / n
	quant := func(v uint8) uint8 {
		idx := int(v) / bandSize
		if idx >= n {
			idx = n - 1
		}
		return uint8(idx * 255 / (n - 1))
	}
	out := imaging.AdjustFunc(src, func(c color.NRGBA) color.NRGBA {
		return color.NRGBA{R: quant(c.R), G: quant(c.G), B: quant(c.B), A: c.A}
	})
	return draw.EncodePNG(out)
}
