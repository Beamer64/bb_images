package bb_images

import (
	_ "embed"
	"errors"
	"image"
	"image/color"
	imgdraw "image/draw"
	"math"
	"strings"

	"github.com/Beamer64/bb_images/internal/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

//go:embed res/sign_templates/change-my-mind.jpg
var changeMyMindBytes []byte

// Sign corners in template-image space, ordered TL, TR, BR, BL.
// Defined as a parallelogram (rotated rectangle) so text scale stays uniform
// across the warp — no perspective stretching on the far side. Center is
// roughly (489, 316), W ≈ 360, H ≈ 100, angle ≈ -25°.
//
// Tuning notes: shift all 4 corners by the same dx/dy to translate. To
// rotate, keep opposite-corner offsets paired (TL/BR and TR/BL). To resize,
// scale all corners away from / toward the center together.
var changeMyMindCorners = [4][2]float64{
	{305, 347},
	{631, 195},
	{673, 285},
	{347, 437},
}

const (
	changeMyMindCanvasW   = 349
	changeMyMindCanvasH   = 130
	changeMyMindMaxFontPt = 72
	changeMyMindMinFontPt = 8
	changeMyMindPadding   = 6
)

// ChangeMyMind renders text onto the "Change My Mind" sign, auto-fit and
// word-wrapped, then perspective-warped to align with the sign's four
// corners. Font is currently x/image's bundled Go Regular; to swap to
// Roboto-Medium, drop the .ttf in res/fonts/ and replace the goregular
// import with an //go:embed-backed []byte, then parse the same way.
func ChangeMyMind(text string) ([]byte, error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return nil, errors.New("change-my-mind: empty text")
	}

	template, err := draw.Decode(changeMyMindBytes)
	if err != nil {
		return nil, err
	}

	textCanvas, err := renderFittedText(text, changeMyMindCanvasW, changeMyMindCanvasH)
	if err != nil {
		return nil, err
	}

	out := warpAndComposite(template, textCanvas, changeMyMindCorners)
	return draw.EncodePNG(out)
}

// renderFittedText draws word-wrapped text centered on a w×h transparent RGBA
// canvas, picking the largest font size in [min, max] that fits in the box
// minus padding via binary search.
func renderFittedText(text string, w, h int) (*image.RGBA, error) {
	ttf, err := opentype.Parse(goregular.TTF)
	if err != nil {
		return nil, err
	}

	innerW := w - 2*changeMyMindPadding
	innerH := h - 2*changeMyMindPadding

	var (
		bestFace  font.Face
		bestLines []string
	)

	lo, hi := changeMyMindMinFontPt, changeMyMindMaxFontPt
	for lo <= hi {
		mid := (lo + hi) / 2
		face, ferr := opentype.NewFace(ttf, &opentype.FaceOptions{
			Size: float64(mid),
			DPI:  72,
		})
		if ferr != nil {
			return nil, ferr
		}
		lines, ok := wrapAndCheck(text, face, innerW, innerH)
		if ok {
			bestFace = face
			bestLines = lines
			lo = mid + 1
		} else {
			hi = mid - 1
		}
	}

	if bestFace == nil {
		face, ferr := opentype.NewFace(ttf, &opentype.FaceOptions{
			Size: float64(changeMyMindMinFontPt),
			DPI:  72,
		})
		if ferr != nil {
			return nil, ferr
		}
		bestFace = face
		bestLines = wrapText(text, face, innerW)
	}

	canvas := image.NewRGBA(image.Rect(0, 0, w, h))
	drawer := &font.Drawer{
		Dst:  canvas,
		Src:  image.NewUniform(color.Black),
		Face: bestFace,
	}
	metrics := bestFace.Metrics()
	lineH := metrics.Height.Round()
	totalH := lineH * len(bestLines)
	baselineY := (h-totalH)/2 + metrics.Ascent.Round()

	for _, line := range bestLines {
		lineW := drawer.MeasureString(line).Round()
		drawer.Dot = fixed.P((w-lineW)/2, baselineY)
		drawer.DrawString(line)
		baselineY += lineH
	}
	return canvas, nil
}

// wrapText greedy-fits words into lines that each measure at most maxW pixels.
func wrapText(text string, face font.Face, maxW int) []string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return nil
	}
	d := &font.Drawer{Face: face}
	lines := make([]string, 0)
	cur := words[0]
	for _, w := range words[1:] {
		candidate := cur + " " + w
		if d.MeasureString(candidate).Round() <= maxW {
			cur = candidate
		} else {
			lines = append(lines, cur)
			cur = w
		}
	}
	return append(lines, cur)
}

// wrapAndCheck wraps text and reports whether every line fits in maxW and the
// total stacked height fits in maxH.
func wrapAndCheck(text string, face font.Face, maxW, maxH int) ([]string, bool) {
	lines := wrapText(text, face, maxW)
	if len(lines) == 0 {
		return lines, false
	}
	if face.Metrics().Height.Round()*len(lines) > maxH {
		return lines, false
	}
	d := &font.Drawer{Face: face}
	for _, line := range lines {
		if d.MeasureString(line).Round() > maxW {
			return lines, false
		}
	}
	return lines, true
}

// warpAndComposite perspective-maps src so its corners land on dstCorners
// of bg, alpha-blending over the background using inverse mapping (each dst
// pixel samples its source counterpart).
func warpAndComposite(bg image.Image, src *image.RGBA, dstCorners [4][2]float64) *image.RGBA {
	sb := src.Bounds()
	srcCorners := [4][2]float64{
		{0, 0},
		{float64(sb.Dx()), 0},
		{float64(sb.Dx()), float64(sb.Dy())},
		{0, float64(sb.Dy())},
	}
	H := computeHomography(dstCorners, srcCorners)

	minX, maxX := math.Inf(1), math.Inf(-1)
	minY, maxY := math.Inf(1), math.Inf(-1)
	for _, c := range dstCorners {
		if c[0] < minX {
			minX = c[0]
		}
		if c[0] > maxX {
			maxX = c[0]
		}
		if c[1] < minY {
			minY = c[1]
		}
		if c[1] > maxY {
			maxY = c[1]
		}
	}

	bb := bg.Bounds()
	out := image.NewRGBA(image.Rect(0, 0, bb.Dx(), bb.Dy()))
	imgdraw.Draw(out, out.Bounds(), bg, bb.Min, imgdraw.Src)

	for y := int(math.Floor(minY)); y <= int(math.Ceil(maxY)); y++ {
		if y < 0 || y >= bb.Dy() {
			continue
		}
		for x := int(math.Floor(minX)); x <= int(math.Ceil(maxX)); x++ {
			if x < 0 || x >= bb.Dx() {
				continue
			}
			u, v := applyHomography(H, float64(x), float64(y))
			iu, iv := int(u), int(v)
			if iu < 0 || iu >= sb.Dx() || iv < 0 || iv >= sb.Dy() {
				continue
			}
			sc := src.RGBAAt(iu, iv)
			if sc.A == 0 {
				continue
			}
			dc := out.RGBAAt(x, y)
			a := float64(sc.A) / 255
			out.SetRGBA(x, y, color.RGBA{
				R: uint8(float64(sc.R)*a + float64(dc.R)*(1-a)),
				G: uint8(float64(sc.G)*a + float64(dc.G)*(1-a)),
				B: uint8(float64(sc.B)*a + float64(dc.B)*(1-a)),
				A: 255,
			})
		}
	}
	return out
}

// computeHomography returns the 3×3 matrix mapping src→dst given four
// matching corner pairs (TL, TR, BR, BL). H[8] is normalized to 1.
func computeHomography(src, dst [4][2]float64) [9]float64 {
	var A [8][8]float64
	var b [8]float64
	for i := 0; i < 4; i++ {
		u, v := src[i][0], src[i][1]
		x, y := dst[i][0], dst[i][1]
		A[2*i] = [8]float64{u, v, 1, 0, 0, 0, -u * x, -v * x}
		b[2*i] = x
		A[2*i+1] = [8]float64{0, 0, 0, u, v, 1, -u * y, -v * y}
		b[2*i+1] = y
	}
	h := solve8(A, b)
	return [9]float64{h[0], h[1], h[2], h[3], h[4], h[5], h[6], h[7], 1}
}

// applyHomography transforms (x, y) by the 3×3 matrix h.
func applyHomography(h [9]float64, x, y float64) (float64, float64) {
	w := h[6]*x + h[7]*y + h[8]
	return (h[0]*x + h[1]*y + h[2]) / w, (h[3]*x + h[4]*y + h[5]) / w
}

// solve8 solves Ax = b for an 8×8 system via Gauss–Jordan with partial pivoting.
func solve8(A [8][8]float64, b [8]float64) [8]float64 {
	for i := 0; i < 8; i++ {
		maxRow := i
		for k := i + 1; k < 8; k++ {
			if math.Abs(A[k][i]) > math.Abs(A[maxRow][i]) {
				maxRow = k
			}
		}
		A[i], A[maxRow] = A[maxRow], A[i]
		b[i], b[maxRow] = b[maxRow], b[i]
		piv := A[i][i]
		for j := 0; j < 8; j++ {
			A[i][j] /= piv
		}
		b[i] /= piv
		for k := 0; k < 8; k++ {
			if k != i && A[k][i] != 0 {
				f := A[k][i]
				for j := 0; j < 8; j++ {
					A[k][j] -= f * A[i][j]
				}
				b[k] -= f * b[i]
			}
		}
	}
	return b
}
