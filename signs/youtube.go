package signs

import (
	_ "embed"
	"errors"
	"fmt"
	"image"
	"image/color"
	imgdraw "image/draw"
	"math/rand/v2"
	"strings"

	"github.com/Beamer64/bb_images/internal/draw"
	"github.com/Beamer64/bb_images/internal/templates"
)

//go:embed res/youtube.png
var youtubeBytes []byte

//go:embed res/markers/youtube.marker.png
var youtubeMarkersBytes []byte

//go:embed res/fonts/roboto-regular.ttf
var youtubeRegularTTF []byte

//go:embed res/fonts/roboto-bold.ttf
var youtubeBoldTTF []byte

// Marker colors for the YouTube comment template.
//   - primary (magenta)  → comment body, drawn in youtubeTextColor
//   - secondary (cyan)   → "@username" in bold, drawn in youtubeTextColor
//   - tertiary (yellow)  → random "X time ago" timestamp, in youtubeTimeColor
//   - avatar (green)     → user's avatar, alpha-masked by the marker shape
var youtubeRoles = map[color.RGBA]string{
	{R: 255, G: 0, B: 255, A: 255}:  "primary",
	{R: 0, G: 255, B: 255, A: 255}:  "secondary",
	{R: 255, G: 255, B: 26, A: 255}: "tertiary",
	{R: 0, G: 255, B: 0, A: 255}:    "avatar",
}

var (
	youtubeTextColor   = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	youtubeTimeColor   = color.RGBA{R: 171, G: 164, B: 153, A: 255}
	youtubeAvatarColor = color.RGBA{G: 255, A: 255}
)

const (
	// @username (secondary) and timestamp (tertiary) share this size so
	// they sit on the same baseline.
	youtubeNameFontPt = 12

	// Comment body — currently double the username pt for visual emphasis.
	youtubeBodyFontPt = 15 // youtubeNameFontPt * 2

	// Horizontal pixels between the end of the username and the start of
	// the timestamp on the same line.
	youtubeNameTimeGap = 8

	// character limit to stop text from entering on a new line
	commentMaxChars = 75
)

// YouTube renders a fake YouTube comment: user avatar, "@username" in bold,
// a random "X time ago" timestamp, and the comment body, composited onto
// the youtube template at the regions painted in the markers PNG. The
// handle is derived from username (lowercased, spaces stripped, prefixed
// with "@"). Avatar is fit to the marker's bounding box and alpha-masked
// by the marker shape.
func YouTube(avatar image.Image, username, comment string) ([]byte, error) {
	comment = strings.TrimSpace(comment)
	username = strings.TrimSpace(username)
	if comment == "" || username == "" {
		return nil, errors.New("youtube: empty comment or username")
	}
	if r := []rune(comment); len(r) > commentMaxChars {
		comment = string(r[:commentMaxChars])
	}

	template, err := draw.Decode(youtubeBytes)
	if err != nil {
		return nil, err
	}
	markersImg, err := draw.Decode(youtubeMarkersBytes)
	if err != nil {
		return nil, err
	}

	regions, err := templates.Detect(markersImg, youtubeRoles, 10)
	if err != nil {
		return nil, err
	}

	tb := template.Bounds()
	out := image.NewRGBA(image.Rect(0, 0, tb.Dx(), tb.Dy()))
	imgdraw.Draw(out, out.Bounds(), template, tb.Min, imgdraw.Src)

	handle := "@" + strings.ToLower(strings.ReplaceAll(username, " ", ""))
	timeAgo := randomYouTubeTime()

	if err := renderTextIntoRegion(
		out, regions["primary"].Bounds, comment,
		textOpts{Color: youtubeTextColor, TTFBytes: youtubeRegularTTF, FixedPt: youtubeBodyFontPt, Align: alignLeft},
	); err != nil {
		return nil, err
	}
	if err := renderTextIntoRegion(
		out, regions["secondary"].Bounds, handle,
		textOpts{Color: youtubeTextColor, TTFBytes: youtubeBoldTTF, FixedPt: youtubeNameFontPt, Align: alignLeft},
	); err != nil {
		return nil, err
	}

	// Position the timestamp on the same baseline as the username, starting
	// right after it. Tertiary marker's Y/height are overridden; only its
	// width survives (so the user's painted region acts as a width hint).
	handleW, err := measureTextWidth(handle, youtubeBoldTTF, youtubeNameFontPt)
	if err != nil {
		return nil, err
	}
	secBounds := regions["secondary"].Bounds
	timeMinX := secBounds.Min.X + textPadding + handleW + youtubeNameTimeGap
	timeBounds := image.Rect(
		timeMinX,
		secBounds.Min.Y,
		timeMinX+regions["tertiary"].Bounds.Dx(),
		secBounds.Max.Y,
	)
	if err := renderTextIntoRegion(
		out, timeBounds, timeAgo,
		textOpts{Color: youtubeTimeColor, TTFBytes: youtubeRegularTTF, FixedPt: youtubeNameFontPt, Align: alignLeft},
	); err != nil {
		return nil, err
	}

	placeMaskedImage(out, avatar, markersImg, youtubeAvatarColor, 10, regions["avatar"].Bounds)

	return draw.EncodePNG(out)
}

// randomYouTubeTime returns a YouTube-style relative timestamp like
// "57 seconds ago" or "3 years ago" — random unit, random count within
// that unit's typical max.
func randomYouTubeTime() string {
	units := []struct {
		singular, plural string
		max              int
	}{
		{"second", "seconds", 59},
		{"minute", "minutes", 59},
		{"hour", "hours", 23},
		{"day", "days", 30},
		{"week", "weeks", 3},
		{"month", "months", 11},
		{"year", "years", 15},
	}
	u := units[rand.IntN(len(units))]
	n := 1 + rand.IntN(u.max)
	if n == 1 {
		return "1 " + u.singular + " ago"
	}
	return fmt.Sprintf("%d %s ago", n, u.plural)
}
