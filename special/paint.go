package special

import (
	"image"
	"image/color"
	imgdraw "image/draw"
	"math"

	"github.com/Beamer64/bb_images/internal/draw"
)

const paintRadius = 3

// Paint applies a Kuwahara filter: for each pixel, it picks the mean of whichever
// of four surrounding regions has the lowest luminance variance, producing a
// painterly, edge-preserving smoothing effect.
func Paint(src image.Image) ([]byte, error) {
	b := src.Bounds()
	w, h := b.Dx(), b.Dy()

	base := image.NewRGBA(image.Rect(0, 0, w, h))
	imgdraw.Draw(base, base.Bounds(), src, b.Min, imgdraw.Src)
	out := image.NewRGBA(image.Rect(0, 0, w, h))

	const r = paintRadius
	for y := r; y < h-r; y++ {
		for x := r; x < w-r; x++ {
			regions := [4][4]int{
				{x - r, y - r, x, y},
				{x, y - r, x + r, y},
				{x - r, y, x, y + r},
				{x, y, x + r, y + r},
			}

			bestVar := math.Inf(1)
			var bestR, bestG, bestB float64

			for _, reg := range regions {
				var sumR, sumG, sumB, sumSq float64
				var n int
				for py := reg[1]; py <= reg[3]; py++ {
					for px := reg[0]; px <= reg[2]; px++ {
						c := base.RGBAAt(px, py)
						sumR += float64(c.R)
						sumG += float64(c.G)
						sumB += float64(c.B)
						lum := 0.299*float64(c.R) + 0.587*float64(c.G) + 0.114*float64(c.B)
						sumSq += lum * lum
						n++
					}
				}
				if n == 0 {
					continue
				}
				fn := float64(n)
				meanR := sumR / fn
				meanG := sumG / fn
				meanB := sumB / fn
				meanLum := 0.299*meanR + 0.587*meanG + 0.114*meanB
				variance := sumSq/fn - meanLum*meanLum
				if variance < bestVar {
					bestVar = variance
					bestR, bestG, bestB = meanR, meanG, meanB
				}
			}

			out.SetRGBA(x, y, color.RGBA{
				R: uint8(bestR),
				G: uint8(bestG),
				B: uint8(bestB),
				A: 255,
			})
		}
	}

	return draw.EncodePNG(out)
}
