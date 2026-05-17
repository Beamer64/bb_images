package animated

import (
	"image"
	"image/color/palette"
	imgdraw "image/draw"
	"image/gif"
	"math/rand/v2"

	"github.com/Beamer64/bb_images/internal/draw"
)

const (
	glitchFrames    = 8
	glitchDelay     = 8
	glitchBandCount = 5
	glitchMaxOffset = 30
)

// Glitch returns src with random horizontal row shifts per frame as an animated GIF.
func Glitch(src image.Image) ([]byte, error) {
	sb := src.Bounds()
	w, h := sb.Dx(), sb.Dy()

	base := image.NewRGBA(image.Rect(0, 0, w, h))
	imgdraw.Draw(base, base.Bounds(), src, sb.Min, imgdraw.Src)

	g := &gif.GIF{LoopCount: 0}
	for f := 0; f < glitchFrames; f++ {
		frame := image.NewRGBA(base.Bounds())
		imgdraw.Draw(frame, frame.Bounds(), base, image.Point{}, imgdraw.Src)

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

		paletted := image.NewPaletted(frame.Bounds(), palette.Plan9)
		imgdraw.FloydSteinberg.Draw(paletted, frame.Bounds(), frame, image.Point{})

		g.Image = append(g.Image, paletted)
		g.Delay = append(g.Delay, glitchDelay)
	}

	return draw.EncodeGIF(g)
}
