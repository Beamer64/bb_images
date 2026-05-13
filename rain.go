package bb_images

import (
	"image"
	"image/color/palette"
	imgdraw "image/draw"
	"image/gif"
	"math/rand/v2"

	"github.com/Beamer64/bb_images/internal/draw"
)

const (
	rainFrames      = 8
	rainDelay       = 5
	rainStreakCount = 40
)

// Rain overlays randomly-placed pale-blue streaks per frame to simulate rain,
// returning an animated GIF.
func Rain(src image.Image) ([]byte, error) {
	sb := src.Bounds()
	w, h := sb.Dx(), sb.Dy()

	base := image.NewRGBA(image.Rect(0, 0, w, h))
	imgdraw.Draw(base, base.Bounds(), src, sb.Min, imgdraw.Src)

	g := &gif.GIF{LoopCount: 0}
	for f := 0; f < rainFrames; f++ {
		frame := image.NewRGBA(base.Bounds())
		imgdraw.Draw(frame, frame.Bounds(), base, image.Point{}, imgdraw.Src)

		for s := 0; s < rainStreakCount; s++ {
			startX := rand.IntN(w)
			startY := rand.IntN(h)
			length := 8 + rand.IntN(15)
			for l := 0; l < length; l++ {
				x := startX - l/8
				y := startY + l
				if x < 0 || x >= w || y < 0 || y >= h {
					continue
				}
				c := frame.RGBAAt(x, y)
				c.R = uint8((uint32(c.R) + 200) / 2)
				c.G = uint8((uint32(c.G) + 220) / 2)
				c.B = uint8((uint32(c.B) + 255) / 2)
				frame.SetRGBA(x, y, c)
			}
		}

		paletted := image.NewPaletted(frame.Bounds(), palette.Plan9)
		imgdraw.FloydSteinberg.Draw(paletted, frame.Bounds(), frame, image.Point{})
		g.Image = append(g.Image, paletted)
		g.Delay = append(g.Delay, rainDelay)
	}

	return draw.EncodeGIF(g)
}
