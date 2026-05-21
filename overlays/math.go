package overlays

import (
	_ "embed"
	"image"
	"image/color/palette"
	imgdraw "image/draw"

	"github.com/Beamer64/bb_images/internal/draw"
	"github.com/disintegration/imaging"
)

//go:embed res/math.gif
var mathBytes []byte

// Cached decoded template — frames are immutable, so a single decode per
// process amortizes across all invocations.
var mathTemplate = draw.NewLazyFrames(mathBytes)

// mathOpacity controls how visible the math overlay is over the source.
// 0.5 reads cleanly without obscuring the face.
const mathOpacity = 0.6

// Math returns src with the animated math overlay composited on top at
// mathOpacity. Output is an animated GIF whose frames share src's
// dimensions; the overlay frame is fit/cropped to match.
func Math(src image.Image) ([]byte, error) {
	b := src.Bounds()
	w, h := b.Dx(), b.Dy()

	return draw.AnimateOverGIF(
		mathTemplate, func(frame *image.RGBA) *image.Paletted {
			resized := imaging.Fill(frame, w, h, imaging.Center, imaging.Linear)
			composited := imaging.Overlay(src, resized, image.Point{}, mathOpacity)

			p := image.NewPaletted(composited.Bounds(), palette.Plan9)
			imgdraw.FloydSteinberg.Draw(p, composited.Bounds(), composited, image.Point{})
			return p
		},
	)
}
