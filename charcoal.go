package bb_images

import (
	"image"

	"github.com/Beamer64/bb_images/internal/draw"
	"github.com/disintegration/imaging"
)

// Charcoal returns src as a soft charcoal drawing — softer edges from a pre-blur
// and a darker paper-toned background to distinguish from Sketch.
func Charcoal(src image.Image) ([]byte, error) {
	blurred := imaging.Blur(src, 1.5)
	edges := sobelEdges(blurred)
	boosted := imaging.AdjustContrast(edges, 30)
	inv := imaging.Invert(boosted)
	return draw.EncodePNG(imaging.AdjustBrightness(inv, -20))
}
