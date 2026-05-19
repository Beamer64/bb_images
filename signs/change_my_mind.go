package signs

import (
	_ "embed"
	"errors"
	"image"
	"image/color"
	imgdraw "image/draw"
	"math"
	"strings"

	"github.com/Beamer64/bb_images/internal/draw"
	"github.com/Beamer64/bb_images/internal/templates"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

//go:embed res/change-my-mind.png
var changeMyMindBytes []byte

//go:embed res/markers/change-my-mind.marker.png
var changeMyMindMarkersBytes []byte

// Marker colors used to detect editable regions in the markers file.
// Pure-channel values resist JPEG/PNG rounding noise; the templates
// package matches each with a small per-channel tolerance.
var changeMyMindRoles = map[color.RGBA]string{
	{R: 255, G: 0, B: 255, A: 255}: "text",
}

const (
	maxFontPt   = 72
	minFontPt   = 8
	textPadding = 6
)

// textAlign selects horizontal alignment for renderFittedText. The zero
// value is alignCenter to preserve existing call sites.
type textAlign int

const (
	alignCenter textAlign = iota
	alignLeft
)

// textOpts controls how renderFittedText draws — color, font (any TTF bytes;
// nil defaults to x/image's bundled Go Regular), an optional fixed point
// size that bypasses auto-fit, and horizontal alignment.
type textOpts struct {
	Color    color.Color
	TTFBytes []byte    // nil = goregular default
	FixedPt  int       // 0 = binary-search auto-fit; >0 = render at exactly this pt
	Align    textAlign // 0 = alignCenter (default)
}

// ChangeMyMind renders text onto the "Change My Mind" sign, auto-fit and
// word-wrapped, then perspective-warped onto the four corners detected
// from the magenta marker region in the parallel markers PNG. Font is
// currently x/image's bundled Go Regular; to swap to Roboto Medium drop
// the .ttf in res/fonts/ and replace the goregular import with an
// //go:embed-backed []byte, then parse the same way.
func ChangeMyMind(text string) ([]byte, error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return nil, errors.New("change-my-mind: empty text")
	}

	template, err := draw.Decode(changeMyMindBytes)
	if err != nil {
		return nil, err
	}
	markers, err := draw.Decode(changeMyMindMarkersBytes)
	if err != nil {
		return nil, err
	}

	regions, err := templates.Detect(markers, changeMyMindRoles, 10)
	if err != nil {
		return nil, err
	}
	r := regions["text"]

	// The sign tilts counterclockwise, so the leftmost extreme is the
	// sign's TL, topmost is TR, rightmost is BR, bottommost is BL.
	corners := [4][2]float64{
		{float64(r.Left.X), float64(r.Left.Y)},
		{float64(r.Top.X), float64(r.Top.Y)},
		{float64(r.Right.X), float64(r.Right.Y)},
		{float64(r.Bottom.X), float64(r.Bottom.Y)},
	}

	// Match the source canvas to the unrolled region's edge lengths so
	// the warp scales text ~1× rather than squashing or stretching it.
	topLen := math.Hypot(corners[1][0]-corners[0][0], corners[1][1]-corners[0][1])
	leftLen := math.Hypot(corners[3][0]-corners[0][0], corners[3][1]-corners[0][1])
	canvasW, canvasH := int(topLen), int(leftLen)
	if canvasW < 1 || canvasH < 1 {
		return nil, errors.New("change-my-mind: degenerate text region")
	}

	textCanvas, err := renderFittedText(text, canvasW, canvasH, textOpts{Color: color.Black})
	if err != nil {
		return nil, err
	}

	out := warpAndComposite(template, textCanvas, corners)
	return draw.EncodePNG(out)
}

// renderFittedText draws word-wrapped text centered on a w×h transparent
// RGBA canvas. With opts.FixedPt == 0 it binary-searches for the largest
// font size in [min, max] that fits in the box minus padding; with
// FixedPt > 0 it skips the search and renders at that exact size.
func renderFittedText(text string, w, h int, opts textOpts) (*image.RGBA, error) {
	ttfBytes := opts.TTFBytes
	if ttfBytes == nil {
		ttfBytes = goregular.TTF
	}
	ttf, err := opentype.Parse(ttfBytes)
	if err != nil {
		return nil, err
	}

	innerW := w - 2*textPadding
	innerH := h - 2*textPadding

	var (
		bestFace  font.Face
		bestLines []string
	)

	if opts.FixedPt > 0 {
		face, ferr := opentype.NewFace(
			ttf, &opentype.FaceOptions{
				Size: float64(opts.FixedPt),
				DPI:  72,
			},
		)
		if ferr != nil {
			return nil, ferr
		}
		bestFace = face
		bestLines = wrapText(text, face, innerW)
	} else {
		lo, hi := minFontPt, maxFontPt
		for lo <= hi {
			mid := (lo + hi) / 2
			face, ferr := opentype.NewFace(
				ttf, &opentype.FaceOptions{
					Size: float64(mid),
					DPI:  72,
				},
			)
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
			face, ferr := opentype.NewFace(
				ttf, &opentype.FaceOptions{
					Size: float64(minFontPt),
					DPI:  72,
				},
			)
			if ferr != nil {
				return nil, ferr
			}
			bestFace = face
			bestLines = wrapText(text, face, innerW)
		}
	}

	canvas := image.NewRGBA(image.Rect(0, 0, w, h))
	drawer := &font.Drawer{
		Dst:  canvas,
		Src:  image.NewUniform(opts.Color),
		Face: bestFace,
	}
	metrics := bestFace.Metrics()
	lineH := metrics.Height.Round()
	totalH := lineH * len(bestLines)
	baselineY := (h-totalH)/2 + metrics.Ascent.Round()

	for _, line := range bestLines {
		lineW := drawer.MeasureString(line).Round()
		var x int
		switch opts.Align {
		case alignLeft:
			x = textPadding
		default:
			x = (w - lineW) / 2
		}
		drawer.Dot = fixed.P(x, baselineY)
		drawer.DrawString(line)
		baselineY += lineH
	}
	return canvas, nil
}

// measureTextWidth returns the rendered pixel width of text at sizePt using
// the supplied TTF bytes (nil falls back to goregular). Uses DPI 72 to match
// renderFittedText, so the returned width lines up with what renderFittedText
// would actually draw.
func measureTextWidth(text string, ttfBytes []byte, sizePt int) (int, error) {
	if ttfBytes == nil {
		ttfBytes = goregular.TTF
	}
	ttf, err := opentype.Parse(ttfBytes)
	if err != nil {
		return 0, err
	}
	face, err := opentype.NewFace(ttf, &opentype.FaceOptions{Size: float64(sizePt), DPI: 72})
	if err != nil {
		return 0, err
	}
	d := &font.Drawer{Face: face}
	return d.MeasureString(text).Round(), nil
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
			out.SetRGBA(
				x, y, color.RGBA{
					R: uint8(float64(sc.R)*a + float64(dc.R)*(1-a)),
					G: uint8(float64(sc.G)*a + float64(dc.G)*(1-a)),
					B: uint8(float64(sc.B)*a + float64(dc.B)*(1-a)),
					A: 255,
				},
			)
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
