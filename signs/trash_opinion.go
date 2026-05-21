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

//go:embed res/trash.png
var trashOpinionBytes []byte

//go:embed res/markers/trash.marker.png
var trashOpinionMarkersBytes []byte

var trashOpinionAvatarColor = color.RGBA{G: 255, A: 255}

// TrashOpinion composites the avatar onto the "trash opinion" template,
// placing the avatar once into each green marker region.
func TrashOpinion(avatar image.Image) ([]byte, error) {
	if avatar == nil {
		return nil, errors.New("trash_opinion: nil avatar")
	}

	template, err := draw.Decode(trashOpinionBytes)
	if err != nil {
		return nil, err
	}
	markersImg, err := draw.Decode(trashOpinionMarkersBytes)
	if err != nil {
		return nil, err
	}

	tb := template.Bounds()
	out := image.NewRGBA(image.Rect(0, 0, tb.Dx(), tb.Dy()))
	imgdraw.Draw(out, out.Bounds(), template, tb.Min, imgdraw.Src)

	regions := templates.ConnectedRegions(markersImg, trashOpinionAvatarColor, 10)
	if len(regions) == 0 {
		return nil, errors.New("trash_opinion: no green marker regions found")
	}
	for _, r := range regions {
		placeMaskedImage(out, avatar, markersImg, trashOpinionAvatarColor, 10, r)
	}

	return draw.EncodePNG(out)
}
