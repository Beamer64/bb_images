package edges

import (
	"image"

	"github.com/Beamer64/bb_images/internal/draw"
	"github.com/disintegration/imaging"
)

// Sketch returns src as a clean pencil sketch (dark lines on white).
func Sketch(src image.Image) ([]byte, error) {
	return draw.EncodePNG(imaging.Invert(sobelEdges(src)))
}
