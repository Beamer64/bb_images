package bb_images

import (
	"image"
	"image/color"
	imgdraw "image/draw"

	"github.com/Beamer64/bb_images/internal/draw"
)

const triangleCellSize = 15

// Triangle reduces src to a low-poly mosaic: each cell is split along one of
// its diagonals (alternating per cell in a herringbone pattern) into two
// triangles, each filled with the average color of the source pixels inside.
func Triangle(src image.Image) ([]byte, error) {
	b := src.Bounds()
	w, h := b.Dx(), b.Dy()

	base := image.NewRGBA(image.Rect(0, 0, w, h))
	imgdraw.Draw(base, base.Bounds(), src, b.Min, imgdraw.Src)
	out := image.NewRGBA(image.Rect(0, 0, w, h))

	for ty := 0; ty < h; ty += triangleCellSize {
		for tx := 0; tx < w; tx += triangleCellSize {
			txEnd := tx + triangleCellSize
			tyEnd := ty + triangleCellSize
			if txEnd > w {
				txEnd = w
			}
			if tyEnd > h {
				tyEnd = h
			}
			cw := float64(txEnd - tx)
			ch := float64(tyEnd - ty)
			flip := ((tx/triangleCellSize)+(ty/triangleCellSize))%2 == 1

			inA := func(lx, ly float64) bool {
				if flip {
					return lx > ly
				}
				return lx+ly < 1
			}

			var sumAR, sumAG, sumAB uint64
			var sumBR, sumBG, sumBB uint64
			var nA, nB int
			for py := ty; py < tyEnd; py++ {
				for px := tx; px < txEnd; px++ {
					lx := float64(px-tx) / cw
					ly := float64(py-ty) / ch
					c := base.RGBAAt(px, py)
					if inA(lx, ly) {
						sumAR += uint64(c.R)
						sumAG += uint64(c.G)
						sumAB += uint64(c.B)
						nA++
					} else {
						sumBR += uint64(c.R)
						sumBG += uint64(c.G)
						sumBB += uint64(c.B)
						nB++
					}
				}
			}

			avgA := color.RGBA{A: 255}
			avgB := color.RGBA{A: 255}
			if nA > 0 {
				avgA.R = uint8(sumAR / uint64(nA))
				avgA.G = uint8(sumAG / uint64(nA))
				avgA.B = uint8(sumAB / uint64(nA))
			}
			if nB > 0 {
				avgB.R = uint8(sumBR / uint64(nB))
				avgB.G = uint8(sumBG / uint64(nB))
				avgB.B = uint8(sumBB / uint64(nB))
			}

			for py := ty; py < tyEnd; py++ {
				for px := tx; px < txEnd; px++ {
					lx := float64(px-tx) / cw
					ly := float64(py-ty) / ch
					if inA(lx, ly) {
						out.SetRGBA(px, py, avgA)
					} else {
						out.SetRGBA(px, py, avgB)
					}
				}
			}
		}
	}

	return draw.EncodePNG(out)
}
