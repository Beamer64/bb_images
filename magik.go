package bb_images

import (
	"image"
	imgdraw "image/draw"
	"math"

	"github.com/Beamer64/bb_images/internal/draw"
)

const (
	magikAmplitude = 10.0
	magikFreq      = 4.0
)

// Magik applies a wavy "liquid" distortion to src by displacing each pixel
// according to two superimposed sinusoidal fields. Cheap approximation of
// content-aware liquid rescaling.
func Magik(src image.Image) ([]byte, error) {
	b := src.Bounds()
	w, h := b.Dx(), b.Dy()

	base := image.NewRGBA(image.Rect(0, 0, w, h))
	imgdraw.Draw(base, base.Bounds(), src, b.Min, imgdraw.Src)
	out := image.NewRGBA(image.Rect(0, 0, w, h))

	fw := float64(w)
	fh := float64(h)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			fy := float64(y)
			fx := float64(x)
			dx := magikAmplitude*math.Sin(fy*magikFreq*math.Pi/fh) +
				(magikAmplitude/2)*math.Sin(fy*magikFreq*3*math.Pi/fh)
			dy := magikAmplitude*math.Cos(fx*magikFreq*math.Pi/fw) +
				(magikAmplitude/2)*math.Cos(fx*magikFreq*3*math.Pi/fw)
			sx := x + int(dx)
			sy := y + int(dy)
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
