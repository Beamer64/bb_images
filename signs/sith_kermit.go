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

//go:embed res/sith-kermit.png
var sithKermitBytes []byte

//go:embed res/markers/sith-kermit.marker.png
var sithKermitMarkersBytes []byte

var (
	sithKermitSithColor   = color.RGBA{G: 255, A: 255}
	sithKermitKermitColor = color.RGBA{B: 250, A: 255}
)

// SithKermit composites the sith and kermit avatars onto the template.
// Green marker pixels carry sith, blue marker pixels carry kermit. Each
// connected component of marker pixels gets its own placement, so the
// template can have multiple regions per role without stretching.
func SithKermit(sith, kermit image.Image) ([]byte, error) {
	if sith == nil || kermit == nil {
		return nil, errors.New("sith_kermit: nil avatar")
	}

	template, err := draw.Decode(sithKermitBytes)
	if err != nil {
		return nil, err
	}
	markersImg, err := draw.Decode(sithKermitMarkersBytes)
	if err != nil {
		return nil, err
	}

	tb := template.Bounds()
	out := image.NewRGBA(image.Rect(0, 0, tb.Dx(), tb.Dy()))
	imgdraw.Draw(out, out.Bounds(), template, tb.Min, imgdraw.Src)

	sithRegions := templates.ConnectedRegions(markersImg, sithKermitSithColor, 10)
	if len(sithRegions) == 0 {
		return nil, errors.New("sith_kermit: no green marker regions found")
	}
	for _, r := range sithRegions {
		placeMaskedImage(out, sith, markersImg, sithKermitSithColor, 10, r)
	}

	kermitRegions := templates.ConnectedRegions(markersImg, sithKermitKermitColor, 10)
	if len(kermitRegions) == 0 {
		return nil, errors.New("sith_kermit: no blue marker regions found")
	}
	for _, r := range kermitRegions {
		placeMaskedImage(out, kermit, markersImg, sithKermitKermitColor, 10, r)
	}

	return draw.EncodePNG(out)
}
