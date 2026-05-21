package overlays

import (
	_ "embed"
	"image"
	"image/color/palette"
	imgdraw "image/draw"

	"github.com/Beamer64/bb_images/internal/draw"
	"github.com/disintegration/imaging"
)

//go:embed res/elmo_burn.gif
var elmoBurnBytes []byte

// Cached decoded template — frames are immutable, so a single decode per
// process amortizes across all invocations.
var elmoBurnTemplate = draw.NewLazyFrames(elmoBurnBytes)

// elmoBurnOpacity controls how visible the burning-Elmo overlay is over
// the source. Tune up to make the flames pop, down to favor the face.
const elmoBurnOpacity = 0.5

// ElmoBurn returns src with the burning-Elmo animated overlay composited
// on top at elmoBurnOpacity. Output is an animated GIF whose frames
// share src's dimensions; the overlay frame is fit/cropped to match.
func ElmoBurn(src image.Image) ([]byte, error) {
	b := src.Bounds()
	w, h := b.Dx(), b.Dy()

	return draw.AnimateOverGIF(elmoBurnTemplate, func(frame *image.RGBA) *image.Paletted {
		resized := imaging.Fill(frame, w, h, imaging.Center, imaging.Linear)
		composited := imaging.Overlay(src, resized, image.Point{}, elmoBurnOpacity)

		p := image.NewPaletted(composited.Bounds(), palette.Plan9)
		imgdraw.FloydSteinberg.Draw(p, composited.Bounds(), composited, image.Point{})
		return p
	})
}
