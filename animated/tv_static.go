package animated

import (
	"image"
	"image/color"
	"image/color/palette"
	imgdraw "image/draw"
	"image/gif"
	"math/rand/v2"

	"github.com/Beamer64/bb_images/internal/draw"
	"github.com/disintegration/imaging"
)

const (
	tvStaticFrames  = 6
	tvStaticDelay   = 5
	tvStaticOpacity = 0.7
)

// TvStatic returns src obscured by per-frame random grayscale noise as an animated GIF.
func TvStatic(src image.Image) ([]byte, error) {
	sb := src.Bounds()
	w, h := sb.Dx(), sb.Dy()

	g := &gif.GIF{LoopCount: 0}
	for f := 0; f < tvStaticFrames; f++ {
		noise := image.NewRGBA(image.Rect(0, 0, w, h))
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				v := uint8(rand.IntN(256))
				noise.SetRGBA(x, y, color.RGBA{R: v, G: v, B: v, A: 255})
			}
		}
		out := imaging.Overlay(src, noise, image.Point{}, tvStaticOpacity)

		paletted := image.NewPaletted(out.Bounds(), palette.Plan9)
		imgdraw.FloydSteinberg.Draw(paletted, out.Bounds(), out, image.Point{})

		g.Image = append(g.Image, paletted)
		g.Delay = append(g.Delay, tvStaticDelay)
	}

	return draw.EncodeGIF(g)
}
