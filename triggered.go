package bb_images

import (
	"image"
	"image/color"
	"image/color/palette"
	imgdraw "image/draw"
	"image/gif"
	"math/rand/v2"

	"github.com/Beamer64/bb_images/internal/draw"
	"github.com/disintegration/imaging"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

const (
	triggeredFrames       = 10
	triggeredDelay        = 3 // hundredths of a second per frame
	triggeredJitter       = 8 // max per-frame pixel offset
	triggeredBannerHeight = 30
)

// Triggered returns the "TRIGGERED" meme animation (shake + red tint + banner) as GIF bytes.
func Triggered(src image.Image) ([]byte, error) {
	sb := src.Bounds()
	w, h := sb.Dx(), sb.Dy()

	// Enlarge once so each jittered crop never exposes edges.
	enlarged := imaging.Resize(src, w+triggeredJitter*2, h+triggeredJitter*2, imaging.Lanczos)
	redTint := image.NewUniform(color.RGBA{R: 255, A: 60})
	bannerBG := image.NewUniform(color.RGBA{R: 220, A: 230})
	textSrc := image.NewUniform(color.White)

	g := &gif.GIF{LoopCount: 0}
	for f := 0; f < triggeredFrames; f++ {
		offX := rand.IntN(triggeredJitter*2 + 1)
		offY := rand.IntN(triggeredJitter*2 + 1)

		frame := image.NewRGBA(image.Rect(0, 0, w, h))
		cropped := imaging.Crop(enlarged, image.Rect(offX, offY, offX+w, offY+h))
		imgdraw.Draw(frame, frame.Bounds(), cropped, image.Point{}, imgdraw.Src)
		imgdraw.Draw(frame, frame.Bounds(), redTint, image.Point{}, imgdraw.Over)

		bannerRect := image.Rect(0, h-triggeredBannerHeight, w, h)
		imgdraw.Draw(frame, bannerRect, bannerBG, image.Point{}, imgdraw.Over)

		d := &font.Drawer{Dst: frame, Src: textSrc, Face: basicfont.Face7x13}
		const label = "TRIGGERED"
		textWidth := d.MeasureString(label).Round()
		d.Dot = fixed.Point26_6{
			X: fixed.I((w - textWidth) / 2),
			Y: fixed.I(h - 10),
		}
		d.DrawString(label)

		paletted := image.NewPaletted(frame.Bounds(), palette.Plan9)
		imgdraw.FloydSteinberg.Draw(paletted, frame.Bounds(), frame, image.Point{})

		g.Image = append(g.Image, paletted)
		g.Delay = append(g.Delay, triggeredDelay)
	}

	return draw.EncodeGIF(g)
}
