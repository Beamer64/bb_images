package signs

import (
	_ "embed"
	"errors"
	"image"
	"image/color"
	imgdraw "image/draw"
	"strings"

	"github.com/Beamer64/bb_images/internal/draw"
	"github.com/Beamer64/bb_images/internal/templates"
	"github.com/disintegration/imaging"
)

//go:embed res/tweet.png
var tweetBytes []byte

//go:embed res/markers/tweet.marker.png
var tweetMarkersBytes []byte

//go:embed res/fonts/chirp-regular.ttf
var tweetRegularTTF []byte

//go:embed res/fonts/chirp-bold.ttf
var tweetBoldTTF []byte

//go:embed res/fonts/chirp-medium.ttf
var tweetMediumTTF []byte

//go:embed res/fonts/chirp-heavy.ttf
var tweetHeavyTTF []byte

// Twitter's character limit. Body text longer than this is truncated with
// an ellipsis so the rendered tweet is faithful to the real product cap.
const tweetMaxChars = 250

// Marker colors for the tweet template.
//   - primary (magenta)  → tweet body text, drawn in tweetTextColor
//   - secondary (cyan)   → display name, drawn in tweetTextColor
//   - tertiary (yellow)  → @handle, drawn in tweetHandleColor
//   - avatar (green)     → user's avatar, alpha-masked by the marker shape
var tweetRoles = map[color.RGBA]string{
	{R: 255, G: 0, B: 255, A: 255}:  "primary",
	{R: 0, G: 255, B: 255, A: 255}:  "secondary",
	{R: 255, G: 255, B: 26, A: 255}: "tertiary",
	{R: 0, G: 255, B: 0, A: 255}:    "avatar",
}

var (
	tweetTextColor   = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	tweetHandleColor = color.RGBA{R: 166, G: 159, B: 147, A: 255}
	tweetAvatarColor = color.RGBA{G: 255, A: 255}
)

// Display name (secondary) and @handle (tertiary) share this fixed pt size
// so they line up visually. Tune up/down to taste; auto-fit is intentionally
// not used for these two regions.
const tweetNameFontPt = 18

// Body text font size — matches Twitter's ~14px tweet text. The opentype
// face DPI is 72 so pt and px are 1:1 here.
const tweetBodyFontPt = 18

// Vertical pixel gap inserted between the bottom of the display-name region
// and the top of the @handle region — keeps the two stacked tightly the way
// they sit in a real Twitter card.
const tweetStackGap = 4

// Tweet renders a fake tweet: avatar, the user's displayed name, their
// "@username" handle (stacked just below the displayed name), and body
// text — composited onto the tweet template at the regions painted in the
// markers PNG. The @handle is derived from username (lowercased, spaces
// stripped, prefixed with "@"). Avatar is fit to the marker's bounding
// box and alpha-masked by the marker shape, so non-rectangular marker
// outlines (circles, rounded squares, etc.) come through cleanly.
func Tweet(avatar image.Image, displayName, username, tweetText string) ([]byte, error) {
	tweetText = strings.TrimSpace(tweetText)
	displayName = strings.TrimSpace(displayName)
	username = strings.TrimSpace(username)
	if tweetText == "" || displayName == "" || username == "" {
		return nil, errors.New("tweet: empty text, display name, or username")
	}
	if r := []rune(tweetText); len(r) > tweetMaxChars {
		tweetText = string(r[:tweetMaxChars])
	}

	template, err := draw.Decode(tweetBytes)
	if err != nil {
		return nil, err
	}
	markersImg, err := draw.Decode(tweetMarkersBytes)
	if err != nil {
		return nil, err
	}

	regions, err := templates.Detect(markersImg, tweetRoles, 10)
	if err != nil {
		return nil, err
	}

	tb := template.Bounds()
	out := image.NewRGBA(image.Rect(0, 0, tb.Dx(), tb.Dy()))
	imgdraw.Draw(out, out.Bounds(), template, tb.Min, imgdraw.Src)

	handle := "@" + strings.ToLower(strings.ReplaceAll(username, " ", ""))
	if err := renderTextIntoRegion(
		out, regions["primary"].Bounds, tweetText, textOpts{Color: tweetTextColor, TTFBytes: tweetRegularTTF, FixedPt: tweetBodyFontPt, Align: alignLeft},
	); err != nil {
		return nil, err
	}
	if err := renderTextIntoRegion(
		out, regions["secondary"].Bounds, displayName, textOpts{Color: tweetTextColor, TTFBytes: tweetBoldTTF, FixedPt: tweetNameFontPt, Align: alignLeft},
	); err != nil {
		return nil, err
	}

	// Stack the @handle directly below the display name with a small gap,
	// overriding the tertiary marker's position so the gap is consistent
	// regardless of where the yellow marker happens to be painted.
	secBounds := regions["secondary"].Bounds
	handleBounds := image.Rect(
		secBounds.Min.X,
		secBounds.Max.Y+tweetStackGap,
		secBounds.Max.X,
		secBounds.Max.Y+tweetStackGap+secBounds.Dy(),
	)
	if err := renderTextIntoRegion(
		out, handleBounds, handle, textOpts{Color: tweetHandleColor, TTFBytes: tweetMediumTTF, FixedPt: tweetNameFontPt, Align: alignLeft},
	); err != nil {
		return nil, err
	}

	placeMaskedImage(out, avatar, markersImg, tweetAvatarColor, 10, regions["avatar"].Bounds)

	return draw.EncodePNG(out)
}

// renderTextIntoRegion renders text into the region's bounding box per opts
// (auto-fit or fixed size, regular or bold, in the requested color), then
// composites it onto out at the same position.
func renderTextIntoRegion(out *image.RGBA, bounds image.Rectangle, text string, opts textOpts) error {
	canvas, err := renderFittedText(text, bounds.Dx(), bounds.Dy(), opts)
	if err != nil {
		return err
	}
	dst := image.Rect(bounds.Min.X, bounds.Min.Y, bounds.Min.X+canvas.Bounds().Dx(), bounds.Min.Y+canvas.Bounds().Dy())
	imgdraw.Draw(out, dst, canvas, image.Point{}, imgdraw.Over)
	return nil
}

// placeMaskedImage resizes src to the bounding box, then copies src pixel-by-
// pixel onto out only where the markers image matches markerColor. This makes
// the marker's painted shape act as the avatar's alpha mask.
func placeMaskedImage(out *image.RGBA, src image.Image, markers image.Image, markerColor color.RGBA, tolerance int, bounds image.Rectangle) {
	sized := imaging.Resize(src, bounds.Dx(), bounds.Dy(), imaging.Lanczos)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := markers.At(x, y).RGBA()
			if absDiff(int(r>>8), int(markerColor.R)) <= tolerance &&
				absDiff(int(g>>8), int(markerColor.G)) <= tolerance &&
				absDiff(int(b>>8), int(markerColor.B)) <= tolerance {
				out.Set(x, y, sized.At(x-bounds.Min.X, y-bounds.Min.Y))
			}
		}
	}
}

func absDiff(a, b int) int {
	if d := a - b; d < 0 {
		return -d
	} else {
		return d
	}
}
