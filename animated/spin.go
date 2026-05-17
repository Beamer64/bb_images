package animated

import (
	"image"
	"image/color"
	"image/color/palette"
	imgdraw "image/draw"
	"image/gif"

	"github.com/Beamer64/bb_images/internal/draw"
	"github.com/disintegration/imaging"
)

const (
	spinFrames = 12
	spinDelay  = 5
)

// Spin rotates src 360° across frames as an animated GIF; corners are clipped.
func Spin(src image.Image) ([]byte, error) {
	sb := src.Bounds()
	w, h := sb.Dx(), sb.Dy()

	g := &gif.GIF{LoopCount: 0}
	for f := 0; f < spinFrames; f++ {
		angle := float64(f) * 360.0 / float64(spinFrames)
		rotated := imaging.Rotate(src, angle, color.Transparent)
		cropped := imaging.CropCenter(rotated, w, h)

		paletted := image.NewPaletted(cropped.Bounds(), palette.Plan9)
		imgdraw.FloydSteinberg.Draw(paletted, cropped.Bounds(), cropped, cropped.Bounds().Min)

		g.Image = append(g.Image, paletted)
		g.Delay = append(g.Delay, spinDelay)
	}

	return draw.EncodeGIF(g)
}
