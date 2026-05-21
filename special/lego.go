package special

import (
	"image"
	"image/color"

	"github.com/Beamer64/bb_images/internal/draw"
)

// Tuning constants for the LEGO mosaic look. Adjust to taste:
//   - smaller brick size → more detail, less "blocky" feel
//   - larger stud radius → bumpier
//   - bevel/highlight multipliers control how strongly the 3D effect reads
//
// Stud radius is roughly 30% of the brick size — encoded as an integer
// directly because Go forbids constant float→int conversions with
// fractional results. If you change legoBrickSize, retune legoStudR by
// hand (≈ brickSize * 0.30).
const (
	legoBrickSize     = 25
	legoStudR         = 7
	legoBevelLighten  = 1.15
	legoBevelDarken   = 0.85
	legoStudLighten   = 1.10
	legoStudHighlight = 1.25
)

// Lego renders src as a mosaic of LEGO-style bricks. Each
// legoBrickSize-sized cell is filled with the average color of the
// corresponding region of src, with subtle bevels (lighter top/left,
// darker bottom/right) and a shaded circular stud at the cell's center.
// The two-tone stud rim — light on the upper-left, dark on the lower-
// right — makes the bumps read as spherical 3D studs instead of flat
// circles.
//
// Output dimensions match src. Cells at the right/bottom edge that
// don't fit a full legoBrickSize get a clipped brick with no stud
// (better than a stud that overflows into the next non-cell area).
func Lego(src image.Image) ([]byte, error) {
	b := src.Bounds()
	w, h := b.Dx(), b.Dy()
	out := image.NewRGBA(image.Rect(0, 0, w, h))

	for by := 0; by < h; by += legoBrickSize {
		for bx := 0; bx < w; bx += legoBrickSize {
			cw := legoBrickSize
			if bx+cw > w {
				cw = w - bx
			}
			ch := legoBrickSize
			if by+ch > h {
				ch = h - by
			}

			base := legoAverageColor(src, b.Min.X+bx, b.Min.Y+by, cw, ch)
			light := legoScaleColor(base, legoBevelLighten)
			dark := legoScaleColor(base, legoBevelDarken)

			legoFillRect(out, bx, by, cw, ch, base)

			// 1px bevels on each edge. Corners get overwritten but always
			// to the same color (light wins on the two top corners after
			// the second pass), which is fine visually.
			for x := bx; x < bx+cw; x++ {
				out.SetRGBA(x, by, light)
				out.SetRGBA(x, by+ch-1, dark)
			}
			for y := by; y < by+ch; y++ {
				out.SetRGBA(bx, y, light)
				out.SetRGBA(bx+cw-1, y, dark)
			}

			// Only render a stud on cells that fit a full brick. Edge
			// cells with partial dimensions skip the stud — clipped studs
			// look worse than missing ones.
			if cw == legoBrickSize && ch == legoBrickSize {
				cx, cy := bx+cw/2, by+ch/2
				studBase := legoScaleColor(base, legoStudLighten)
				studHigh := legoScaleColor(base, legoStudHighlight)
				legoDrawStud(out, cx, cy, legoStudR, studBase, studHigh, dark)
			}
		}
	}

	return draw.EncodePNG(out)
}

// legoAverageColor returns the mean RGB of the (x, y, x+w, y+h)
// rectangle in src. Alpha is fixed at 255 — bricks are opaque.
func legoAverageColor(src image.Image, x, y, w, h int) color.RGBA {
	var rSum, gSum, bSum, n uint64
	for j := 0; j < h; j++ {
		for i := 0; i < w; i++ {
			r, g, b, _ := src.At(x+i, y+j).RGBA()
			rSum += uint64(r >> 8)
			gSum += uint64(g >> 8)
			bSum += uint64(b >> 8)
			n++
		}
	}
	if n == 0 {
		return color.RGBA{A: 255}
	}
	return color.RGBA{
		R: uint8(rSum / n),
		G: uint8(gSum / n),
		B: uint8(bSum / n),
		A: 255,
	}
}

// legoFillRect paints the (x, y, x+w, y+h) box in out with c.
func legoFillRect(out *image.RGBA, x, y, w, h int, c color.RGBA) {
	for j := 0; j < h; j++ {
		for i := 0; i < w; i++ {
			out.SetRGBA(x+i, y+j, c)
		}
	}
}

// legoDrawStud paints a filled circle of radius r centered at (cx, cy)
// with studBase, then re-paints the outer ring's upper-left quadrant
// with studHighlight and the lower-right quadrant with studShadow. The
// upper-right and lower-left quadrants of the ring stay studBase, which
// gives a soft transition between the two tones.
func legoDrawStud(out *image.RGBA, cx, cy, r int, studBase, studHighlight, studShadow color.RGBA) {
	r2 := r * r
	rInner2 := (r - 1) * (r - 1)
	for dy := -r; dy <= r; dy++ {
		for dx := -r; dx <= r; dx++ {
			d2 := dx*dx + dy*dy
			if d2 > r2 {
				continue
			}
			if d2 < rInner2 {
				out.SetRGBA(cx+dx, cy+dy, studBase)
				continue
			}
			// Outer ring pixel: tint by quadrant.
			switch {
			case dx <= 0 && dy <= 0:
				out.SetRGBA(cx+dx, cy+dy, studHighlight)
			case dx >= 0 && dy >= 0:
				out.SetRGBA(cx+dx, cy+dy, studShadow)
			default:
				out.SetRGBA(cx+dx, cy+dy, studBase)
			}
		}
	}
}

// legoScaleColor multiplies each RGB channel by factor and clamps to
// [0, 255]. Used to derive bevel highlights / shadows / stud tones from
// the brick's base color.
func legoScaleColor(c color.RGBA, factor float64) color.RGBA {
	return color.RGBA{
		R: legoClamp8(float64(c.R) * factor),
		G: legoClamp8(float64(c.G) * factor),
		B: legoClamp8(float64(c.B) * factor),
		A: 255,
	}
}

func legoClamp8(v float64) uint8 {
	if v < 0 {
		return 0
	}
	if v > 255 {
		return 255
	}
	return uint8(v)
}
