package bb_images

import (
	"image"
	"image/color"
	"math"

	"github.com/Beamer64/bb_images/internal/draw"
	"github.com/disintegration/imaging"
)

const (
	hogCellSize = 10
	hogBins     = 10
)

// Hog renders the histogram-of-oriented-gradients visualization of src — each
// cell of the output shows oriented edge-direction lines whose length and
// brightness reflect how much gradient mass falls in each orientation bin.
func Hog(src image.Image) ([]byte, error) {
	gray := imaging.Grayscale(src)
	b := gray.Bounds()
	w, h := b.Dx(), b.Dy()

	gx := make([]float64, w*h)
	gy := make([]float64, w*h)
	for y := 1; y < h-1; y++ {
		for x := 1; x < w-1; x++ {
			tl := float64(gray.NRGBAAt(x-1, y-1).R)
			tc := float64(gray.NRGBAAt(x, y-1).R)
			tr := float64(gray.NRGBAAt(x+1, y-1).R)
			ml := float64(gray.NRGBAAt(x-1, y).R)
			mr := float64(gray.NRGBAAt(x+1, y).R)
			bl := float64(gray.NRGBAAt(x-1, y+1).R)
			bc := float64(gray.NRGBAAt(x, y+1).R)
			br := float64(gray.NRGBAAt(x+1, y+1).R)
			gx[y*w+x] = -tl + tr - 2*ml + 2*mr - bl + br
			gy[y*w+x] = -tl - 2*tc - tr + bl + 2*bc + br
		}
	}

	out := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			out.SetRGBA(x, y, color.RGBA{A: 255})
		}
	}

	halfSize := float64(hogCellSize) / 2
	for cy := 0; cy+hogCellSize <= h; cy += hogCellSize {
		for cx := 0; cx+hogCellSize <= w; cx += hogCellSize {
			bins := make([]float64, hogBins)
			for py := cy; py < cy+hogCellSize; py++ {
				for px := cx; px < cx+hogCellSize; px++ {
					gxv, gyv := gx[py*w+px], gy[py*w+px]
					mag := math.Sqrt(gxv*gxv + gyv*gyv)
					theta := math.Atan2(gyv, gxv)
					if theta < 0 {
						theta += math.Pi
					}
					idx := int(theta / math.Pi * float64(hogBins))
					if idx >= hogBins {
						idx = hogBins - 1
					}
					bins[idx] += mag
				}
			}

			cellMax := 0.0
			for _, v := range bins {
				if v > cellMax {
					cellMax = v
				}
			}
			if cellMax == 0 {
				continue
			}

			ccx := float64(cx) + halfSize
			ccy := float64(cy) + halfSize
			for i, weight := range bins {
				if weight == 0 {
					continue
				}
				strength := weight / cellMax
				length := halfSize * strength
				// Edges are perpendicular to the gradient direction.
				edgeAngle := (float64(i)+0.5)*math.Pi/float64(hogBins) + math.Pi/2
				dx := length * math.Cos(edgeAngle)
				dy := length * math.Sin(edgeAngle)
				shade := uint8(255 * strength)
				drawLine(
					out,
					int(ccx-dx), int(ccy-dy),
					int(ccx+dx), int(ccy+dy),
					color.RGBA{R: shade, G: shade, B: shade, A: 255},
				)
			}
		}
	}

	return draw.EncodePNG(out)
}

// drawLine rasterizes a line via Bresenham. Pixels outside img's bounds are skipped.
func drawLine(img *image.RGBA, x0, y0, x1, y1 int, c color.RGBA) {
	dx := absInt(x1 - x0)
	dy := absInt(y1 - y0)
	sx := 1
	if x0 >= x1 {
		sx = -1
	}
	sy := 1
	if y0 >= y1 {
		sy = -1
	}
	err := dx - dy
	bounds := img.Bounds()
	for {
		if x0 >= bounds.Min.X && x0 < bounds.Max.X && y0 >= bounds.Min.Y && y0 < bounds.Max.Y {
			img.SetRGBA(x0, y0, c)
		}
		if x0 == x1 && y0 == y1 {
			return
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x0 += sx
		}
		if e2 < dx {
			err += dx
			y0 += sy
		}
	}
}

func absInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
