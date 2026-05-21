package overlays

import (
	_ "embed"
	"image"
	"image/color/palette"
	imgdraw "image/draw"

	"github.com/Beamer64/bb_images/internal/draw"
	"github.com/disintegration/imaging"
)

//go:embed res/communism.gif
var communismBytes []byte

// Cached decoded template — frames are immutable, so a single decode per
// process amortizes across all invocations.
var communismTemplate = draw.NewLazyFrames(communismBytes)

// communismOpacity controls how visible the hammer-and-sickle / red wash
// is over the source. 0.5 reads cleanly without obscuring the face.
const communismOpacity = 0.5

// Communism returns src with the animated communism overlay composited
// on top at communismOpacity. Output is an animated GIF whose frames
// share src's dimensions; the overlay is fit/cropped to match.
func Communism(src image.Image) ([]byte, error) {
	b := src.Bounds()
	w, h := b.Dx(), b.Dy()

	return draw.AnimateOverGIF(communismTemplate, func(frame *image.RGBA) *image.Paletted {
		// Bilinear resize is sufficient here: any high-frequency detail
		// in the overlay would be lost in the 256-color dither anyway,
		// and Linear is ~2× faster than Lanczos.
		resized := imaging.Fill(frame, w, h, imaging.Center, imaging.Linear)
		composited := imaging.Overlay(src, resized, image.Point{}, communismOpacity)

		// Quantize back to a paletted image for GIF encoding. Plan9's
		// 256-color palette is the same one animated/ uses — good
		// general coverage with stdlib-only Floyd-Steinberg dithering.
		p := image.NewPaletted(composited.Bounds(), palette.Plan9)
		imgdraw.FloydSteinberg.Draw(p, composited.Bounds(), composited, image.Point{})
		return p
	})
}
