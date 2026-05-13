package bb_images

import (
	"image"
	"image/color"
	"math"

	"github.com/Beamer64/bb_images/internal/draw"
	"github.com/disintegration/imaging"
)

const (
	stringifyNails      = 200
	stringifyIterations = 1000
	stringifyLineDarken = 30
)

// Stringify approximates src as classic string art: nail points on a circle
// inscribed in the image, then a greedy loop adds the line whose path most
// reduces the residual "darkness needed" field.
//
// Per-line sums and a pixel→line reverse index are cached so each iteration
// only updates the lines that share pixels with the selected line instead of
// rescanning every candidate. This keeps the inner loop O(line_length × avg
// lines-per-pixel) rather than O(numLines × line_length).
func Stringify(src image.Image) ([]byte, error) {
	gray := imaging.Grayscale(src)
	b := gray.Bounds()
	w, h := b.Dx(), b.Dy()

	cx := float64(w) / 2
	cy := float64(h) / 2
	radius := math.Min(cx, cy) - 1

	nails := make([][2]int, stringifyNails)
	for i := 0; i < stringifyNails; i++ {
		angle := 2 * math.Pi * float64(i) / float64(stringifyNails)
		nails[i] = [2]int{
			int(cx + radius*math.Cos(angle)),
			int(cy + radius*math.Sin(angle)),
		}
	}

	residual := make([]int, w*h)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			residual[y*w+x] = 255 - int(gray.NRGBAAt(x, y).R)
		}
	}

	numLines := stringifyNails * (stringifyNails - 1) / 2
	lines := make([][]int, numLines)
	lineSum := make([]int, numLines)
	pixelToLines := make([][]int, w*h)

	idx := 0
	for i := 0; i < stringifyNails; i++ {
		for j := i + 1; j < stringifyNails; j++ {
			pixels := bresenhamPath(nails[i][0], nails[i][1], nails[j][0], nails[j][1], w, h)
			lines[idx] = pixels
			sum := 0
			for _, p := range pixels {
				sum += residual[p]
				pixelToLines[p] = append(pixelToLines[p], idx)
			}
			lineSum[idx] = sum
			idx++
		}
	}

	out := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			out.SetRGBA(x, y, color.RGBA{R: 255, G: 255, B: 255, A: 255})
		}
	}

	for it := 0; it < stringifyIterations; it++ {
		bestSum := 0
		bestIdx := -1
		for i, s := range lineSum {
			if s > bestSum {
				bestSum = s
				bestIdx = i
			}
		}
		if bestIdx < 0 {
			break
		}

		for _, p := range lines[bestIdx] {
			oldR := residual[p]
			newR := oldR - stringifyLineDarken
			if newR < 0 {
				newR = 0
			}
			delta := newR - oldR
			residual[p] = newR

			for _, li := range pixelToLines[p] {
				lineSum[li] += delta
			}

			y := p / w
			x := p % w
			c := out.RGBAAt(x, y)
			v := int(c.R) - 60
			if v < 0 {
				v = 0
			}
			out.SetRGBA(x, y, color.RGBA{R: uint8(v), G: uint8(v), B: uint8(v), A: 255})
		}
	}

	return draw.EncodePNG(out)
}

// bresenhamPath returns the in-bounds pixel indices (y*w+x) along the line
// from (x0,y0) to (x1,y1) in a w×h image.
func bresenhamPath(x0, y0, x1, y1, w, h int) []int {
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
	out := make([]int, 0, dx+dy)
	for {
		if x0 >= 0 && x0 < w && y0 >= 0 && y0 < h {
			out = append(out, y0*w+x0)
		}
		if x0 == x1 && y0 == y1 {
			return out
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
