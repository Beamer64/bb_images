package signs

import (
	_ "embed"
	"errors"
	"image"
	"image/color"
	imgdraw "image/draw"
	"math/rand/v2"
	"strings"
	"time"

	"github.com/Beamer64/bb_images/internal/draw"
	"github.com/Beamer64/bb_images/internal/templates"
)

//go:embed res/discord.png
var discordBytes []byte

//go:embed res/markers/discord.marker.png
var discordMarkersBytes []byte

//go:embed res/fonts/roboto-regular.ttf
var discordRegularTTF []byte

//go:embed res/fonts/roboto-bold.ttf
var discordBoldTTF []byte

// Marker colors for the Discord message template.
//   - primary (magenta)  → message body, drawn in discordBodyColor
//   - secondary (cyan)   → display name, drawn in discordNameColor (bold)
//   - tertiary (yellow)  → "MM/DD/YYYY HH:MM AM/PM" timestamp, drawn in discordDateColor
//   - avatar (green)     → user's avatar, alpha-masked by the marker shape
var discordRoles = map[color.RGBA]string{
	{R: 255, G: 0, B: 255, A: 255}:  "primary",
	{R: 0, G: 255, B: 255, A: 255}:  "secondary",
	{R: 255, G: 255, B: 26, A: 255}: "tertiary",
	{R: 0, G: 255, B: 0, A: 255}:    "avatar",
}

var (
	discordBodyColor   = color.RGBA{R: 220, G: 221, B: 222, A: 255}
	discordNameColor   = color.RGBA{R: 235, G: 250, B: 236, A: 255}
	discordDateColor   = color.RGBA{R: 158, G: 166, B: 165, A: 255}
	discordAvatarColor = color.RGBA{G: 255, A: 255}
)

const (
	discordBodyFontPt      = 24
	discordNameFontPt      = 24
	discordDateFontPt      = 16
	discordMessageMaxChars = 90

	// Pixel gap between the end of the display name and the start of the
	// timestamp on the same baseline.
	discordNameDateGap = 12
)

// Discord renders a fake Discord message: user avatar, display name, an
// inline "MM/DD/YYYY HH:MM AM/PM" timestamp positioned right after the
// name, and the message body — composited onto the Discord template at
// the regions painted in the markers PNG. Avatar is fit to the marker's
// bounding box and alpha-masked by the marker shape.
func Discord(avatar image.Image, displayName, message string) ([]byte, error) {
	message = strings.TrimSpace(message)
	displayName = strings.TrimSpace(displayName)
	if message == "" || displayName == "" {
		return nil, errors.New("discord: empty message or display name")
	}
	if r := []rune(message); len(r) > discordMessageMaxChars {
		message = string(r[:discordMessageMaxChars])
	}

	template, err := draw.Decode(discordBytes)
	if err != nil {
		return nil, err
	}
	markersImg, err := draw.Decode(discordMarkersBytes)
	if err != nil {
		return nil, err
	}

	regions, err := templates.Detect(markersImg, discordRoles, 10)
	if err != nil {
		return nil, err
	}

	tb := template.Bounds()
	out := image.NewRGBA(image.Rect(0, 0, tb.Dx(), tb.Dy()))
	imgdraw.Draw(out, out.Bounds(), template, tb.Min, imgdraw.Src)

	timestamp := randomDiscordTimestamp()

	if err := renderTextIntoRegion(
		out, regions["primary"].Bounds, message,
		textOpts{Color: discordBodyColor, TTFBytes: discordRegularTTF, FixedPt: discordBodyFontPt, Align: alignLeft},
	); err != nil {
		return nil, err
	}
	if err := renderTextIntoRegion(
		out, regions["secondary"].Bounds, displayName,
		textOpts{Color: discordNameColor, TTFBytes: discordBoldTTF, FixedPt: discordNameFontPt, Align: alignLeft},
	); err != nil {
		return nil, err
	}

	// Position the timestamp on the same baseline as the display name,
	// right after it. Tertiary marker's Y/height are overridden; only its
	// width survives as a hint for how much horizontal room is available.
	nameW, err := measureTextWidth(displayName, discordBoldTTF, discordNameFontPt)
	if err != nil {
		return nil, err
	}
	secBounds := regions["secondary"].Bounds
	timeMinX := secBounds.Min.X + textPadding + nameW + discordNameDateGap
	timeBounds := image.Rect(
		timeMinX,
		secBounds.Min.Y,
		timeMinX+regions["tertiary"].Bounds.Dx(),
		secBounds.Max.Y,
	)
	if err := renderTextIntoRegion(
		out, timeBounds, timestamp,
		textOpts{Color: discordDateColor, TTFBytes: discordRegularTTF, FixedPt: discordDateFontPt, Align: alignLeft},
	); err != nil {
		return nil, err
	}

	placeMaskedImage(out, avatar, markersImg, discordAvatarColor, 10, regions["avatar"].Bounds)

	return draw.EncodePNG(out)
}

// randomDiscordTimestamp returns a Discord-style timestamp formatted as
// "MM/DD/YYYY HH:MM AM/PM" at a random instant within the last seven years.
func randomDiscordTimestamp() string {
	now := time.Now()
	earliest := now.AddDate(-7, 0, 0)
	span := now.Unix() - earliest.Unix()
	offset := rand.Int64N(span)
	return earliest.Add(time.Duration(offset) * time.Second).Format("01/02/2006 03:04 PM")
}
