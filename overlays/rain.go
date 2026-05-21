package overlays

import (
	_ "embed"
	"image"
	"image/color/palette"
	imgdraw "image/draw"

	"github.com/Beamer64/bb_images/internal/draw"
	"github.com/disintegration/imaging"
)

//go:embed res/rain.gif
var rainBytes []byte

// Cached decoded template — frames are immutable, so a single decode per
// process amortizes across all invocations.
var rainTemplate = draw.NewLazyFrames(rainBytes)

// rainOpacity controls how visible the rain overlay is over the source.
// Tune up to drench, down for a light drizzle.
const rainOpacity = 0.5

// Rain returns src with the animated rain overlay composited on top at
// rainOpacity. Output is an animated GIF whose frames share src's
// dimensions; the overlay frame is fit/cropped to match.
//
// Note: a procedural rain implementation also lives in animated/rain.go
// (synthesized streaks rather than a pre-rendered GIF). They produce
// different aesthetics; this is the GIF-overlay version.
func Rain(src image.Image) ([]byte, error) {
	b := src.Bounds()
	w, h := b.Dx(), b.Dy()

	return draw.AnimateOverGIF(rainTemplate, func(frame *image.RGBA) *image.Paletted {
		resized := imaging.Fill(frame, w, h, imaging.Center, imaging.Linear)
		composited := imaging.Overlay(src, resized, image.Point{}, rainOpacity)

		p := image.NewPaletted(composited.Bounds(), palette.Plan9)
		imgdraw.FloydSteinberg.Draw(p, composited.Bounds(), composited, image.Point{})
		return p
	})
}
