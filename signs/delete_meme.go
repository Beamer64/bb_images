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

//go:embed res/delete.png
var deleteMemeBytes []byte

//go:embed res/markers/delete.marker.png
var deleteMemeMarkersBytes []byte

var deleteMemeImageColor = color.RGBA{G: 255, A: 255}

// DeleteMeme composites an arbitrary meme image onto the green marker
// region(s) of the Windows-error template. Unlike most signs commands,
// the placed image isn't a Discord avatar — callers pass any decoded
// image (typically fetched from a user-supplied URL), so the same
// template can be reused with random meme imagery instead of always
// using the invoker's avatar.
func DeleteMeme(meme image.Image) ([]byte, error) {
	if meme == nil {
		return nil, errors.New("delete_meme: nil image")
	}

	template, err := draw.Decode(deleteMemeBytes)
	if err != nil {
		return nil, err
	}
	markersImg, err := draw.Decode(deleteMemeMarkersBytes)
	if err != nil {
		return nil, err
	}

	tb := template.Bounds()
	out := image.NewRGBA(image.Rect(0, 0, tb.Dx(), tb.Dy()))
	imgdraw.Draw(out, out.Bounds(), template, tb.Min, imgdraw.Src)

	regions := templates.ConnectedRegions(markersImg, deleteMemeImageColor, 10)
	if len(regions) == 0 {
		return nil, errors.New("delete_meme: no green marker regions found")
	}
	for _, r := range regions {
		placeMaskedImage(out, meme, markersImg, deleteMemeImageColor, 10, r)
	}

	return draw.EncodePNG(out)
}
