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

//go:embed res/fedora.png
var fedoraBytes []byte

//go:embed res/markers/fedora.marker.png
var fedoraMarkersBytes []byte

var fedoraAvatarColor = color.RGBA{G: 255, A: 255}

// Fedora places the avatar at the green marker region(s) as the base
// layer, then draws the fedora template on top using its own alpha
// channel. This is the inverse layering of tweet/discord/retro_meme,
// where the template is the base and the avatar sits on top. Here the
// template's opaque pixels (the hat itself) cover whatever's behind them
// — so the brim can sit *in front of* the avatar where the two overlap,
// making the hat look like it's resting on the user's head.
func Fedora(avatar image.Image) ([]byte, error) {
	if avatar == nil {
		return nil, errors.New("fedora: nil avatar")
	}

	template, err := draw.Decode(fedoraBytes)
	if err != nil {
		return nil, err
	}
	markersImg, err := draw.Decode(fedoraMarkersBytes)
	if err != nil {
		return nil, err
	}

	tb := template.Bounds()
	out := image.NewRGBA(image.Rect(0, 0, tb.Dx(), tb.Dy()))

	// 1. Avatar in the marker region(s) — base layer.
	regions := templates.ConnectedRegions(markersImg, fedoraAvatarColor, 10)
	if len(regions) == 0 {
		return nil, errors.New("fedora: no green marker regions found")
	}
	for _, r := range regions {
		placeMaskedImage(out, avatar, markersImg, fedoraAvatarColor, 10, r)
	}

	// 2. Template on top. Anywhere the template is opaque, it covers the
	//    avatar; anywhere the template is transparent, the avatar shows
	//    through (or the canvas stays transparent if no avatar landed there).
	imgdraw.Draw(out, out.Bounds(), template, tb.Min, imgdraw.Over)

	return draw.EncodePNG(out)
}
