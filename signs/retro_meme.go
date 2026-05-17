package signs

import (
	_ "embed"
	"image"
	"image/color"
	imgdraw "image/draw"
	"strings"

	"github.com/Beamer64/bb_images/internal/draw"
	"github.com/Beamer64/bb_images/internal/templates"
	"github.com/disintegration/imaging"
)

//go:embed res/retro-meme.png
var retroMemeBytes []byte

//go:embed res/markers/retro-meme.marker.png
var retroMemeMarkersBytes []byte

//go:embed res/fonts/impact-regular.ttf
var impactTTF []byte

// Marker colors for the retro-meme template — just two regions, the top
// and bottom text bands; the avatar fills the entire output canvas so
// there's no avatar marker for this template.
var retroMemeRoles = map[color.RGBA]string{
	{R: 255, G: 0, B: 255, A: 255}: "top",
	{R: 0, G: 255, B: 255, A: 255}: "bottom",
}

var (
	retroMemeFillColor   = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	retroMemeStrokeColor = color.RGBA{R: 0, G: 0, B: 0, A: 255}
)

const (
	retroMemeFontPt      = 50
	retroMemeStrokeWidth = 2
	retroMemeMaxChars    = 70
)

// RetroMeme renders a classic top/bottom-text meme: the user's avatar
// fills the whole canvas (sized from the template image), and the
// optional top/bottom strings are overlaid in Impact with a black stroke
// outline. Either text may be empty — only the regions with non-empty
// text are rendered, so callers can pass just one or even neither and
// get a clean avatar-only image back.
func RetroMeme(avatar image.Image, topText, bottomText string) ([]byte, error) {
	topText = strings.TrimSpace(topText)
	bottomText = strings.TrimSpace(bottomText)

	if r := []rune(topText); len(r) > retroMemeMaxChars {
		topText = string(r[:retroMemeMaxChars])
	}
	if r := []rune(bottomText); len(r) > retroMemeMaxChars {
		bottomText = string(r[:retroMemeMaxChars])
	}

	template, err := draw.Decode(retroMemeBytes)
	if err != nil {
		return nil, err
	}
	markersImg, err := draw.Decode(retroMemeMarkersBytes)
	if err != nil {
		return nil, err
	}

	tb := template.Bounds()
	w, h := tb.Dx(), tb.Dy()

	out := image.NewRGBA(image.Rect(0, 0, w, h))
	fitted := imaging.Fill(avatar, w, h, imaging.Center, imaging.Lanczos)
	imgdraw.Draw(out, out.Bounds(), fitted, image.Point{}, imgdraw.Src)

	if topText == "" && bottomText == "" {
		return draw.EncodePNG(out)
	}

	regions, err := templates.Detect(markersImg, retroMemeRoles, 10)
	if err != nil {
		return nil, err
	}

	if topText != "" {
		if err := renderOutlinedText(out, regions["top"].Bounds, topText); err != nil {
			return nil, err
		}
	}
	if bottomText != "" {
		if err := renderOutlinedText(out, regions["bottom"].Bounds, bottomText); err != nil {
			return nil, err
		}
	}

	return draw.EncodePNG(out)
}

// renderOutlinedText draws text centered in bounds with a black stroke
// outline (rendered as eight offset copies of the text) and a white fill
// on top — the classic meme-text look used by the retro-meme template.
func renderOutlinedText(out *image.RGBA, bounds image.Rectangle, text string) error {
	strokeOpts := textOpts{
		Color:    retroMemeStrokeColor,
		TTFBytes: impactTTF,
		FixedPt:  retroMemeFontPt,
		Align:    alignCenter,
	}
	for _, dx := range []int{-retroMemeStrokeWidth, 0, retroMemeStrokeWidth} {
		for _, dy := range []int{-retroMemeStrokeWidth, 0, retroMemeStrokeWidth} {
			if dx == 0 && dy == 0 {
				continue
			}
			canvas, err := renderFittedText(text, bounds.Dx(), bounds.Dy(), strokeOpts)
			if err != nil {
				return err
			}
			dst := image.Rect(
				bounds.Min.X+dx,
				bounds.Min.Y+dy,
				bounds.Min.X+dx+bounds.Dx(),
				bounds.Min.Y+dy+bounds.Dy(),
			)
			imgdraw.Draw(out, dst, canvas, image.Point{}, imgdraw.Over)
		}
	}
	fillOpts := textOpts{
		Color:    retroMemeFillColor,
		TTFBytes: impactTTF,
		FixedPt:  retroMemeFontPt,
		Align:    alignCenter,
	}
	return renderTextIntoRegion(out, bounds, text, fillOpts)
}
