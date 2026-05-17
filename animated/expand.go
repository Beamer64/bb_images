package animated

import (
	"image"
	"image/color/palette"
	imgdraw "image/draw"
	"image/gif"
	"math"

	"github.com/Beamer64/bb_images/internal/draw"
	"github.com/disintegration/imaging"
)

const (
	expandFrames   = 12
	expandDelay    = 5
	expandMaxScale = 1.5
)

// Expand returns an animation that horizontally stretches src out and back via a
// sinusoidal scale curve, encoded as a GIF.
func Expand(src image.Image) ([]byte, error) {
	sb := src.Bounds()
	w, h := sb.Dx(), sb.Dy()

	g := &gif.GIF{LoopCount: 0}
	for f := 0; f < expandFrames; f++ {
		t := float64(f) / float64(expandFrames-1)
		scale := 1.0 + (expandMaxScale-1.0)*math.Sin(t*math.Pi)
		newW := int(float64(w) * scale)

		stretched := imaging.Resize(src, newW, h, imaging.Lanczos)
		frame := image.NewRGBA(image.Rect(0, 0, w, h))
		offX := (w - newW) / 2
		imgdraw.Draw(frame, image.Rect(offX, 0, offX+newW, h), stretched, image.Point{}, imgdraw.Src)

		paletted := image.NewPaletted(frame.Bounds(), palette.Plan9)
		imgdraw.FloydSteinberg.Draw(paletted, frame.Bounds(), frame, image.Point{})

		g.Image = append(g.Image, paletted)
		g.Delay = append(g.Delay, expandDelay)
	}

	return draw.EncodeGIF(g)
}
