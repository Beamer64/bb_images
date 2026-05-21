package overlays

import (
	_ "embed"
	"image"
	"image/color/palette"
	imgdraw "image/draw"

	"github.com/Beamer64/bb_images/internal/draw"
	"github.com/disintegration/imaging"
)

//go:embed res/bomb.gif
var bombBytes []byte

// Cached decoded template — frames are immutable, so a single decode per
// process amortizes across all invocations.
var bombTemplate = draw.NewLazyFrames(bombBytes)

// bombOpacity controls how visible the explosion overlay is over the
// source. Tune up to make the blast pop, down to favor the face.
const bombOpacity = 0.6

// Bomb returns src with the animated explosion overlay composited on
// top at bombOpacity. Output is an animated GIF whose frames share
// src's dimensions; the overlay frame is fit/cropped to match.
func Bomb(src image.Image) ([]byte, error) {
	b := src.Bounds()
	w, h := b.Dx(), b.Dy()

	return draw.AnimateOverGIF(
		bombTemplate, func(frame *image.RGBA) *image.Paletted {
			resized := imaging.Fill(frame, w, h, imaging.Center, imaging.Linear)
			composited := imaging.Overlay(src, resized, image.Point{}, bombOpacity)

			p := image.NewPaletted(composited.Bounds(), palette.Plan9)
			imgdraw.FloydSteinberg.Draw(p, composited.Bounds(), composited, image.Point{})
			return p
		},
	)
}
