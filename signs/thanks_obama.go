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

//go:embed res/ThanksObama.png
var thanksObamaBytes []byte

//go:embed res/markers/ThanksObama.marker.png
var thanksObamaMarkersBytes []byte

var thanksObamaAvatarColor = color.RGBA{G: 255, A: 255}

// ThanksObama composites the avatar onto the "Thanks, Obama" template,
// placing the avatar once into each green marker region so multiple slots
// in the template each get their own copy.
func ThanksObama(avatar image.Image) ([]byte, error) {
	if avatar == nil {
		return nil, errors.New("thanks_obama: nil avatar")
	}

	template, err := draw.Decode(thanksObamaBytes)
	if err != nil {
		return nil, err
	}
	markersImg, err := draw.Decode(thanksObamaMarkersBytes)
	if err != nil {
		return nil, err
	}

	tb := template.Bounds()
	out := image.NewRGBA(image.Rect(0, 0, tb.Dx(), tb.Dy()))
	imgdraw.Draw(out, out.Bounds(), template, tb.Min, imgdraw.Src)

	regions := templates.ConnectedRegions(markersImg, thanksObamaAvatarColor, 10)
	if len(regions) == 0 {
		return nil, errors.New("thanks_obama: no green marker regions found")
	}
	for _, r := range regions {
		placeMaskedImage(out, avatar, markersImg, thanksObamaAvatarColor, 10, r)
	}

	return draw.EncodePNG(out)
}
