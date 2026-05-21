package overlays

import (
	_ "embed"
	"image"

	"github.com/Beamer64/bb_images/internal/draw"
	"github.com/disintegration/imaging"
)

//go:embed res/broken-glass.png
var shatterBytes []byte

// Shatter composites the broken-glass overlay over src at full opacity.
// The PNG ships with a transparent background and only the glass shards
// are opaque, so a full-opacity overlay reveals src everywhere except
// where the cracks/shards sit. src keeps its native dimensions; the
// overlay is fit/cropped to match using imaging.Fill so the cracks stay
// in proportion.
func Shatter(src image.Image) ([]byte, error) {
	overlay, err := draw.Decode(shatterBytes)
	if err != nil {
		return nil, err
	}
	b := src.Bounds()
	resized := imaging.Fill(overlay, b.Dx(), b.Dy(), imaging.Center, imaging.Lanczos)
	out := imaging.Overlay(src, resized, image.Point{}, 1.0)
	return draw.EncodePNG(out)
}
