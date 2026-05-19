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

//go:embed res/BatmanSlap.png
var batmanSlapBytes []byte

//go:embed res/markers/BatmanSlap.marker.png
var batmanSlapMarkersBytes []byte

var (
	batmanSlapBatmanColor = color.RGBA{G: 255, A: 255}
	batmanSlapRobinColor  = color.RGBA{B: 250, A: 255}
)

// BatmanSlap composites the batman and robin avatars onto the classic
// "Batman slapping Robin" meme template. Green marker pixels carry batman,
// blue marker pixels carry robin. Each connected component of marker
// pixels gets its own placement so the template can have multiple regions
// per role without stretching.
func BatmanSlap(batman, robin image.Image) ([]byte, error) {
	if batman == nil || robin == nil {
		return nil, errors.New("batman_slap: nil avatar")
	}

	template, err := draw.Decode(batmanSlapBytes)
	if err != nil {
		return nil, err
	}
	markersImg, err := draw.Decode(batmanSlapMarkersBytes)
	if err != nil {
		return nil, err
	}

	tb := template.Bounds()
	out := image.NewRGBA(image.Rect(0, 0, tb.Dx(), tb.Dy()))
	imgdraw.Draw(out, out.Bounds(), template, tb.Min, imgdraw.Src)

	batmanRegions := templates.ConnectedRegions(markersImg, batmanSlapBatmanColor, 10)
	if len(batmanRegions) == 0 {
		return nil, errors.New("batman_slap: no green marker regions found")
	}
	for _, r := range batmanRegions {
		placeMaskedImage(out, batman, markersImg, batmanSlapBatmanColor, 10, r)
	}

	robinRegions := templates.ConnectedRegions(markersImg, batmanSlapRobinColor, 10)
	if len(robinRegions) == 0 {
		return nil, errors.New("batman_slap: no blue marker regions found")
	}
	for _, r := range robinRegions {
		placeMaskedImage(out, robin, markersImg, batmanSlapRobinColor, 10, r)
	}

	return draw.EncodePNG(out)
}
