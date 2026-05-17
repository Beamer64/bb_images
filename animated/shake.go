package animated

import (
	"image"
	"image/color/palette"
	imgdraw "image/draw"
	"image/gif"
	"math/rand/v2"

	"github.com/Beamer64/bb_images/internal/draw"
	"github.com/disintegration/imaging"
)

const (
	shakeFrames = 10
	shakeDelay  = 3
	shakeJitter = 8
)

// Shake returns src jittered randomly across frames as an animated GIF.
func Shake(src image.Image) ([]byte, error) {
	sb := src.Bounds()
	w, h := sb.Dx(), sb.Dy()

	enlarged := imaging.Resize(src, w+shakeJitter*2, h+shakeJitter*2, imaging.Lanczos)

	g := &gif.GIF{LoopCount: 0}
	for f := 0; f < shakeFrames; f++ {
		offX := rand.IntN(shakeJitter*2 + 1)
		offY := rand.IntN(shakeJitter*2 + 1)

		frame := image.NewRGBA(image.Rect(0, 0, w, h))
		cropped := imaging.Crop(enlarged, image.Rect(offX, offY, offX+w, offY+h))
		imgdraw.Draw(frame, frame.Bounds(), cropped, image.Point{}, imgdraw.Src)

		paletted := image.NewPaletted(frame.Bounds(), palette.Plan9)
		imgdraw.FloydSteinberg.Draw(paletted, frame.Bounds(), frame, image.Point{})

		g.Image = append(g.Image, paletted)
		g.Delay = append(g.Delay, shakeDelay)
	}

	return draw.EncodeGIF(g)
}
