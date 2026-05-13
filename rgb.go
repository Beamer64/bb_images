package bb_images

import (
	"image"
	"image/color"

	"github.com/Beamer64/bb_images/internal/draw"
)

const (
	rgbChartW = 512
	rgbChartH = 256
)

// RGB returns a 512×256 PNG histogram of src's R/G/B channel distributions,
// with additive blending where bars overlap (yellow/cyan/magenta/white).
func RGB(src image.Image) ([]byte, error) {
	b := src.Bounds()

	var histR, histG, histB [256]int
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			r, g, bl, _ := src.At(x, y).RGBA()
			histR[r>>8]++
			histG[g>>8]++
			histB[bl>>8]++
		}
	}

	maxCount := 1
	for i := 0; i < 256; i++ {
		if histR[i] > maxCount {
			maxCount = histR[i]
		}
		if histG[i] > maxCount {
			maxCount = histG[i]
		}
		if histB[i] > maxCount {
			maxCount = histB[i]
		}
	}

	out := image.NewRGBA(image.Rect(0, 0, rgbChartW, rgbChartH))
	bg := color.RGBA{A: 255}
	for y := 0; y < rgbChartH; y++ {
		for x := 0; x < rgbChartW; x++ {
			out.SetRGBA(x, y, bg)
		}
	}

	barWidth := rgbChartW / 256
	for i := 0; i < 256; i++ {
		rh := histR[i] * rgbChartH / maxCount
		gh := histG[i] * rgbChartH / maxCount
		bh := histB[i] * rgbChartH / maxCount
		for dx := 0; dx < barWidth; dx++ {
			x := i*barWidth + dx
			for y := rgbChartH - rh; y < rgbChartH; y++ {
				c := out.RGBAAt(x, y)
				c.R = 255
				out.SetRGBA(x, y, c)
			}
			for y := rgbChartH - gh; y < rgbChartH; y++ {
				c := out.RGBAAt(x, y)
				c.G = 255
				out.SetRGBA(x, y, c)
			}
			for y := rgbChartH - bh; y < rgbChartH; y++ {
				c := out.RGBAAt(x, y)
				c.B = 255
				out.SetRGBA(x, y, c)
			}
		}
	}

	return draw.EncodePNG(out)
}
