package spatial

import (
	"image"

	"github.com/Beamer64/bb_images/internal/draw"
	"github.com/disintegration/imaging"
)

// Mirror returns src flipped horizontally (left-right) as PNG bytes.
func Mirror(src image.Image) ([]byte, error) {
	return draw.EncodePNG(imaging.FlipH(src))
}
