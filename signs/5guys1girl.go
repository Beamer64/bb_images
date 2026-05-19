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

//go:embed res/5guys1girl.png
var fiveGuysOneGirlBytes []byte

//go:embed res/markers/5guys1girl.marker.png
var fiveGuysOneGirlMarkersBytes []byte

var (
	fiveGuysOneGirlGuysColor = color.RGBA{G: 255, A: 255}
	fiveGuysOneGirlGirlColor = color.RGBA{B: 250, A: 255}
)

// FiveGuysOneGirl composites the guys and girl avatars onto the 5guys1girl
// template. The marker image contains multiple disjoint green regions (one
// per "guy") plus blue region(s) for the girl. Each connected component of
// marker pixels gets its own placement of the corresponding avatar — so the
// guys avatar is drawn into each green square individually instead of being
// stretched across their union bounding box.
func FiveGuysOneGirl(guys, girl image.Image) ([]byte, error) {
	if guys == nil || girl == nil {
		return nil, errors.New("5guys1girl: nil avatar")
	}

	template, err := draw.Decode(fiveGuysOneGirlBytes)
	if err != nil {
		return nil, err
	}
	markersImg, err := draw.Decode(fiveGuysOneGirlMarkersBytes)
	if err != nil {
		return nil, err
	}

	tb := template.Bounds()
	out := image.NewRGBA(image.Rect(0, 0, tb.Dx(), tb.Dy()))
	imgdraw.Draw(out, out.Bounds(), template, tb.Min, imgdraw.Src)

	guysRegions := templates.ConnectedRegions(markersImg, fiveGuysOneGirlGuysColor, 10)
	if len(guysRegions) == 0 {
		return nil, errors.New("5guys1girl: no green marker regions found")
	}
	for _, r := range guysRegions {
		placeMaskedImage(out, guys, markersImg, fiveGuysOneGirlGuysColor, 10, r)
	}

	girlRegions := templates.ConnectedRegions(markersImg, fiveGuysOneGirlGirlColor, 10)
	if len(girlRegions) == 0 {
		return nil, errors.New("5guys1girl: no blue marker regions found")
	}
	for _, r := range girlRegions {
		placeMaskedImage(out, girl, markersImg, fiveGuysOneGirlGirlColor, 10, r)
	}

	return draw.EncodePNG(out)
}
