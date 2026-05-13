package bb_images

import (
	"image"
	"image/color"

	"github.com/Beamer64/bb_images/internal/draw"
)

const mosaicTileSize = 12

// Mosaic replaces each mosaicTileSize×mosaicTileSize block of src with the
// block's average color, separated by a dark grout line on the right/bottom edge.
func Mosaic(src image.Image) ([]byte, error) {
	b := src.Bounds()
	w, h := b.Dx(), b.Dy()

	out := image.NewRGBA(image.Rect(0, 0, w, h))
	grout := color.RGBA{R: 30, G: 30, B: 30, A: 255}

	for ty := 0; ty < h; ty += mosaicTileSize {
		for tx := 0; tx < w; tx += mosaicTileSize {
			var sumR, sumG, sumB, count uint64
			tyMax := ty + mosaicTileSize
			if tyMax > h {
				tyMax = h
			}
			txMax := tx + mosaicTileSize
			if txMax > w {
				txMax = w
			}
			for y := ty; y < tyMax; y++ {
				for x := tx; x < txMax; x++ {
					r, g, bl, _ := src.At(x, y).RGBA()
					sumR += uint64(r >> 8)
					sumG += uint64(g >> 8)
					sumB += uint64(bl >> 8)
					count++
				}
			}
			if count == 0 {
				continue
			}
			avg := color.RGBA{
				R: uint8(sumR / count),
				G: uint8(sumG / count),
				B: uint8(sumB / count),
				A: 255,
			}
			for y := ty; y < tyMax; y++ {
				for x := tx; x < txMax; x++ {
					if x == txMax-1 || y == tyMax-1 {
						out.SetRGBA(x, y, grout)
					} else {
						out.SetRGBA(x, y, avg)
					}
				}
			}
		}
	}

	return draw.EncodePNG(out)
}
