package signs

import (
	_ "embed"
	"errors"
	"image"
	"image/color"
	imgdraw "image/draw"

	"github.com/Beamer64/bb_images/internal/draw"
	"github.com/Beamer64/bb_images/internal/templates"
)

//go:embed res/bad.png
var badBytes []byte

//go:embed res/markers/bad.marker.png
var badMarkersBytes []byte

var badAvatarColor = color.RGBA{G: 255, A: 255}

// Bad composites the avatar onto the "bad" template, placing the avatar
// once into each green marker region.
func Bad(avatar image.Image) ([]byte, error) {
	if avatar == nil {
		return nil, errors.New("bad: nil avatar")
	}

	template, err := draw.Decode(badBytes)
	if err != nil {
		return nil, err
	}
	markersImg, err := draw.Decode(badMarkersBytes)
	if err != nil {
		return nil, err
	}

	tb := template.Bounds()
	out := image.NewRGBA(image.Rect(0, 0, tb.Dx(), tb.Dy()))
	imgdraw.Draw(out, out.Bounds(), template, tb.Min, imgdraw.Src)

	regions := templates.ConnectedRegions(markersImg, badAvatarColor, 10)
	if len(regions) == 0 {
		return nil, errors.New("bad: no green marker regions found")
	}
	for _, r := range regions {
		placeMaskedImage(out, avatar, markersImg, badAvatarColor, 10, r)
	}

	return draw.EncodePNG(out)
}
