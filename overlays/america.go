package overlays

import (
	_ "embed"
	"image"
	"image/color/palette"
	imgdraw "image/draw"

	"github.com/Beamer64/bb_images/internal/draw"
	"github.com/disintegration/imaging"
)

//go:embed res/USflag.gif
var americaFlagBytes []byte

// Cached decoded template — frames are immutable, so a single decode per
// process amortizes across all invocations.
var americaTemplate = draw.NewLazyFrames(americaFlagBytes)

// americaOpacity controls how visible the flag is over the source.
// 0.4 keeps the face clearly readable while still showing the flag.
const americaOpacity = 0.4

// America returns src with the animated American flag overlay composited
// on top at americaOpacity. Output is an animated GIF whose frames share
// src's dimensions; the flag frame is fit/cropped to match.
func America(src image.Image) ([]byte, error) {
	b := src.Bounds()
	w, h := b.Dx(), b.Dy()

	return draw.AnimateOverGIF(americaTemplate, func(frame *image.RGBA) *image.Paletted {
		resized := imaging.Fill(frame, w, h, imaging.Center, imaging.Linear)
		composited := imaging.Overlay(src, resized, image.Point{}, americaOpacity)

		p := image.NewPaletted(composited.Bounds(), palette.Plan9)
		imgdraw.FloydSteinberg.Draw(p, composited.Bounds(), composited, image.Point{})
		return p
	})
}
