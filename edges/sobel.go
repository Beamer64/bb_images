package edges

import (
	"image"
	"image/color"
	"math"

	"github.com/Beamer64/bb_images/internal/draw"
	"github.com/disintegration/imaging"
)

// Sobel returns src as a grayscale edge-detection map (bright edges on black).
func Sobel(src image.Image) ([]byte, error) {
	return draw.EncodePNG(sobelEdges(src))
}

// sobelEdges runs a 3x3 Sobel kernel over a grayscale of src and returns the
// per-pixel edge magnitude as an *image.Gray. Shared by Sobel/Sketch/Charcoal.
func sobelEdges(src image.Image) *image.Gray {
	gray := imaging.Grayscale(src)
	b := gray.Bounds()
	w, h := b.Dx(), b.Dy()

	out := image.NewGray(image.Rect(0, 0, w, h))

	at := func(x, y int) float64 { return float64(gray.NRGBAAt(x, y).R) }

	for y := 1; y < h-1; y++ {
		for x := 1; x < w-1; x++ {
			tl, tc, tr := at(x-1, y-1), at(x, y-1), at(x+1, y-1)
			ml, mr := at(x-1, y), at(x+1, y)
			bl, bc, br := at(x-1, y+1), at(x, y+1), at(x+1, y+1)

			gx := -tl + tr - 2*ml + 2*mr - bl + br
			gy := -tl - 2*tc - tr + bl + 2*bc + br

			mag := math.Sqrt(gx*gx + gy*gy)
			if mag > 255 {
				mag = 255
			}
			out.SetGray(x, y, color.Gray{Y: uint8(mag)})
		}
	}

	return out
}
