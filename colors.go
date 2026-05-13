package bb_images

import (
	"errors"
	"image"
	"image/color"
	"sort"

	"github.com/Beamer64/bb_images/internal/draw"
)

const (
	colorsSwatchWidth  = 500
	colorsSwatchHeight = 100
)

// Colors returns a 500x100 PNG swatch of the k most dominant colors in src.
// Colors are quantized to 5 bits per channel, so visually-similar shades
// are grouped together before counting.
func Colors(src image.Image, k int) ([]byte, error) {
	if k < 1 {
		return nil, errors.New("colors: k must be >= 1")
	}

	counts := make(map[uint32]int)
	samples := make(map[uint32]color.RGBA)

	b := src.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			r, g, bl, a := src.At(x, y).RGBA()
			if a < 0x8000 {
				continue
			}
			qr := uint32(r>>11) & 0x1F
			qg := uint32(g>>11) & 0x1F
			qb := uint32(bl>>11) & 0x1F
			key := qr<<10 | qg<<5 | qb
			counts[key]++
			if _, ok := samples[key]; !ok {
				samples[key] = color.RGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(bl >> 8), A: 255}
			}
		}
	}

	if len(counts) == 0 {
		return nil, errors.New("colors: no opaque pixels in source")
	}

	type bucket struct {
		key   uint32
		count int
	}
	buckets := make([]bucket, 0, len(counts))
	for ky, c := range counts {
		buckets = append(buckets, bucket{ky, c})
	}
	sort.Slice(buckets, func(i, j int) bool { return buckets[i].count > buckets[j].count })

	if len(buckets) < k {
		k = len(buckets)
	}

	out := image.NewRGBA(image.Rect(0, 0, colorsSwatchWidth, colorsSwatchHeight))
	stripeW := colorsSwatchWidth / k
	for i := 0; i < k; i++ {
		col := samples[buckets[i].key]
		x0 := i * stripeW
		x1 := x0 + stripeW
		if i == k-1 {
			x1 = colorsSwatchWidth
		}
		for y := 0; y < colorsSwatchHeight; y++ {
			for x := x0; x < x1; x++ {
				out.Set(x, y, col)
			}
		}
	}

	return draw.EncodePNG(out)
}
