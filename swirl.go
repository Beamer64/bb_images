package bb_images

import (
	"image"
	imgdraw "image/draw"
	"math"

	"github.com/Beamer64/bb_images/internal/draw"
)

const swirlStrength = 1.0

// Swirl applies a polar swirl distortion centered on the image. Pixels near
// the center rotate up to ±swirlStrength*π radians; the edge stays in place.
func Swirl(src image.Image) ([]byte, error) {
	b := src.Bounds()
	w, h := b.Dx(), b.Dy()

	base := image.NewRGBA(image.Rect(0, 0, w, h))
	imgdraw.Draw(base, base.Bounds(), src, b.Min, imgdraw.Src)

	cx, cy := float64(w)/2, float64(h)/2
	maxR := math.Sqrt(cx*cx + cy*cy)

	out := image.NewRGBA(image.Rect(0, 0, w, h))

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			dx := float64(x) - cx
			dy := float64(y) - cy
			r := math.Sqrt(dx*dx + dy*dy)

			theta := math.Atan2(dy, dx)
			newTheta := theta - math.Pi*swirlStrength*(1-r/maxR)

			sx := int(math.Round(cx + r*math.Cos(newTheta)))
			sy := int(math.Round(cy + r*math.Sin(newTheta)))
			if sx < 0 {
				sx = 0
			}
			if sx >= w {
				sx = w - 1
			}
			if sy < 0 {
				sy = 0
			}
			if sy >= h {
				sy = h - 1
			}
			out.SetRGBA(x, y, base.RGBAAt(sx, sy))
		}
	}

	return draw.EncodePNG(out)
}
