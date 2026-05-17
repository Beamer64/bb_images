package special

import (
	"image"
	"image/color"

	"github.com/Beamer64/bb_images/internal/draw"
	"github.com/disintegration/imaging"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

const asciiChars = " .:-=+*#%@"

// Ascii renders src as ASCII art: each 7×13 cell is replaced with a character
// whose density matches the cell's mean luminance, drawn on a black background.
func Ascii(src image.Image) ([]byte, error) {
	sb := src.Bounds()
	w, h := sb.Dx(), sb.Dy()

	const cellW, cellH = 7, 13
	cols := w / cellW
	rows := h / cellH

	gray := imaging.Grayscale(src)

	out := image.NewRGBA(image.Rect(0, 0, w, h))
	bg := color.RGBA{A: 255}
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			out.SetRGBA(x, y, bg)
		}
	}

	d := &font.Drawer{
		Dst:  out,
		Src:  image.NewUniform(color.White),
		Face: basicfont.Face7x13,
	}

	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			cx := col*cellW + cellW/2
			cy := row*cellH + cellH/2
			intensity := float64(gray.NRGBAAt(cx, cy).R) / 255.0
			idx := int(intensity * float64(len(asciiChars)-1))
			d.Dot = fixed.Point26_6{
				X: fixed.I(col * cellW),
				Y: fixed.I((row + 1) * cellH),
			}
			d.DrawString(string(asciiChars[idx]))
		}
	}

	return draw.EncodePNG(out)
}
