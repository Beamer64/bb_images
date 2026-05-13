package bb_images

import (
	"image"
	"image/color"
	"image/color/palette"
	imgdraw "image/draw"
	"image/gif"
	"math/rand/v2"

	"github.com/Beamer64/bb_images/internal/draw"
)

const (
	glitchStaticFrames        = 8
	glitchStaticDelay         = 6
	glitchStaticNoiseFraction = 0.3 // 30% of pixels overwritten with random noise
)

// GlitchStatic combines glitch row-shifting with TV-static speckle for an
// extra-noisy animated GIF.
func GlitchStatic(src image.Image) ([]byte, error) {
	sb := src.Bounds()
	w, h := sb.Dx(), sb.Dy()

	base := image.NewRGBA(image.Rect(0, 0, w, h))
	imgdraw.Draw(base, base.Bounds(), src, sb.Min, imgdraw.Src)

	g := &gif.GIF{LoopCount: 0}
	for f := 0; f < glitchStaticFrames; f++ {
		frame := image.NewRGBA(base.Bounds())
		imgdraw.Draw(frame, frame.Bounds(), base, image.Point{}, imgdraw.Src)

		// Row-band shifts.
		for b := 0; b < glitchBandCount; b++ {
			bandY := rand.IntN(h)
			bandH := 5 + rand.IntN(15)
			if bandY+bandH > h {
				bandH = h - bandY
			}
			offset := rand.IntN(glitchMaxOffset*2+1) - glitchMaxOffset
			for y := bandY; y < bandY+bandH; y++ {
				for x := 0; x < w; x++ {
					sx := ((x-offset)%w + w) % w
					frame.SetRGBA(x, y, base.RGBAAt(sx, y))
				}
			}
		}

		// Speckle noise on top.
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				if rand.Float64() < glitchStaticNoiseFraction {
					v := uint8(rand.IntN(256))
					frame.SetRGBA(x, y, color.RGBA{R: v, G: v, B: v, A: 255})
				}
			}
		}

		paletted := image.NewPaletted(frame.Bounds(), palette.Plan9)
		imgdraw.FloydSteinberg.Draw(paletted, frame.Bounds(), frame, image.Point{})
		g.Image = append(g.Image, paletted)
		g.Delay = append(g.Delay, glitchStaticDelay)
	}

	return draw.EncodeGIF(g)
}
