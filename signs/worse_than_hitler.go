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

//go:embed res/worse-than-hitler.png
var worseThanHitlerBytes []byte

//go:embed res/markers/worse-than-hitler.marker.png
var worseThanHitlerMarkersBytes []byte

var worseThanHitlerAvatarColor = color.RGBA{G: 255, A: 255}

// WorseThanHitler composites the avatar onto the "worse than Hitler"
// template, placing the avatar once into each green marker region.
func WorseThanHitler(avatar image.Image) ([]byte, error) {
	if avatar == nil {
		return nil, errors.New("worse_than_hitler: nil avatar")
	}

	template, err := draw.Decode(worseThanHitlerBytes)
	if err != nil {
		return nil, err
	}
	markersImg, err := draw.Decode(worseThanHitlerMarkersBytes)
	if err != nil {
		return nil, err
	}

	tb := template.Bounds()
	out := image.NewRGBA(image.Rect(0, 0, tb.Dx(), tb.Dy()))
	imgdraw.Draw(out, out.Bounds(), template, tb.Min, imgdraw.Src)

	regions := templates.ConnectedRegions(markersImg, worseThanHitlerAvatarColor, 10)
	if len(regions) == 0 {
		return nil, errors.New("worse_than_hitler: no green marker regions found")
	}
	for _, r := range regions {
		placeMaskedImage(out, avatar, markersImg, worseThanHitlerAvatarColor, 10, r)
	}

	return draw.EncodePNG(out)
}
