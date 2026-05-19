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

//go:embed res/whyAreYouGay.png
var whyAreYouGayBytes []byte

//go:embed res/markers/whyAreYouGay.marker.png
var whyAreYouGayMarkersBytes []byte

var (
	whyAreYouGayIntervieweeColor = color.RGBA{G: 255, A: 255}
	whyAreYouGayInterviewerColor = color.RGBA{B: 250, A: 255}
)

// WhyAreYouGay composites the interviewee and interviewer avatars onto the
// "why are you gay" interview template. Green marker pixels carry the
// interviewee avatar, blue marker pixels carry the interviewer avatar.
// Each connected component of marker pixels gets its own placement, so the
// template can have multiple regions per role without stretching.
func WhyAreYouGay(interviewee, interviewer image.Image) ([]byte, error) {
	if interviewee == nil || interviewer == nil {
		return nil, errors.New("why_are_you_gay: nil avatar")
	}

	template, err := draw.Decode(whyAreYouGayBytes)
	if err != nil {
		return nil, err
	}
	markersImg, err := draw.Decode(whyAreYouGayMarkersBytes)
	if err != nil {
		return nil, err
	}

	tb := template.Bounds()
	out := image.NewRGBA(image.Rect(0, 0, tb.Dx(), tb.Dy()))
	imgdraw.Draw(out, out.Bounds(), template, tb.Min, imgdraw.Src)

	intervieweeRegions := templates.ConnectedRegions(markersImg, whyAreYouGayIntervieweeColor, 10)
	if len(intervieweeRegions) == 0 {
		return nil, errors.New("why_are_you_gay: no green marker regions found")
	}
	for _, r := range intervieweeRegions {
		placeMaskedImage(out, interviewee, markersImg, whyAreYouGayIntervieweeColor, 10, r)
	}

	interviewerRegions := templates.ConnectedRegions(markersImg, whyAreYouGayInterviewerColor, 10)
	if len(interviewerRegions) == 0 {
		return nil, errors.New("why_are_you_gay: no blue marker regions found")
	}
	for _, r := range interviewerRegions {
		placeMaskedImage(out, interviewer, markersImg, whyAreYouGayInterviewerColor, 10, r)
	}

	return draw.EncodePNG(out)
}
