package bb_images

import (
	"image"
	"image/color"
	"math/rand/v2"

	"github.com/Beamer64/bb_images/internal/draw"
	"github.com/disintegration/imaging"
)

// Burn overlays a "Doom-fire" simulated flame field over src — seeded along
// the bottom row, propagated upward with random horizontal drift and cooling,
// then colored via a black→red→orange→yellow→white palette.
func Burn(src image.Image) ([]byte, error) {
	b := src.Bounds()
	w, h := b.Dx(), b.Dy()

	heat := make([]int, w*h)
	for x := 0; x < w; x++ {
		heat[(h-1)*w+x] = rand.IntN(256)
	}
	for y := h - 2; y >= 0; y-- {
		for x := 0; x < w; x++ {
			sx := x + rand.IntN(3) - 1
			if sx < 0 {
				sx = 0
			}
			if sx >= w {
				sx = w - 1
			}
			v := heat[(y+1)*w+sx] - rand.IntN(3)
			if v < 0 {
				v = 0
			}
			heat[y*w+x] = v
		}
	}

	fire := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			v := heat[y*w+x]
			if v == 0 {
				continue
			}
			r, g, bl := fireColor(v)
			alpha := v * 2
			if alpha > 255 {
				alpha = 255
			}
			fire.SetRGBA(x, y, color.RGBA{R: r, G: g, B: bl, A: uint8(alpha)})
		}
	}

	out := imaging.Overlay(src, fire, image.Point{}, 1.0)
	return draw.EncodePNG(out)
}

// fireColor maps a heat intensity (0..255) to the classic fire palette.
func fireColor(intensity int) (uint8, uint8, uint8) {
	switch {
	case intensity < 64:
		return uint8(intensity * 4), 0, 0
	case intensity < 128:
		return 255, uint8((intensity - 64) * 2), 0
	case intensity < 192:
		return 255, uint8(128 + (intensity-128)*2), 0
	default:
		return 255, 255, uint8((intensity - 192) * 4)
	}
}
