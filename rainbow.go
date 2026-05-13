package bb_images

import (
	"image"
	"image/color"
	"image/color/palette"
	imgdraw "image/draw"
	"image/gif"
	"math"

	"github.com/Beamer64/bb_images/internal/draw"
	"github.com/disintegration/imaging"
)

const (
	rainbowFrames = 12
	rainbowDelay  = 8
)

// Rainbow shifts the hue of every pixel by a full 360° across frames as an animated GIF.
func Rainbow(src image.Image) ([]byte, error) {
	g := &gif.GIF{LoopCount: 0}
	for f := 0; f < rainbowFrames; f++ {
		shift := float64(f) * 360.0 / float64(rainbowFrames)

		frame := imaging.AdjustFunc(src, func(c color.NRGBA) color.NRGBA {
			h, s, l := rgbToHSL(c.R, c.G, c.B)
			h = math.Mod(h+shift, 360)
			r, gg, b := hslToRGB(h, s, l)
			return color.NRGBA{R: r, G: gg, B: b, A: c.A}
		})

		paletted := image.NewPaletted(frame.Bounds(), palette.Plan9)
		imgdraw.FloydSteinberg.Draw(paletted, frame.Bounds(), frame, frame.Bounds().Min)

		g.Image = append(g.Image, paletted)
		g.Delay = append(g.Delay, rainbowDelay)
	}
	return draw.EncodeGIF(g)
}

func rgbToHSL(rb, gb, bb uint8) (h, s, l float64) {
	r := float64(rb) / 255
	g := float64(gb) / 255
	b := float64(bb) / 255

	maxC := math.Max(math.Max(r, g), b)
	minC := math.Min(math.Min(r, g), b)
	l = (maxC + minC) / 2

	if maxC == minC {
		return 0, 0, l
	}

	d := maxC - minC
	if l > 0.5 {
		s = d / (2 - maxC - minC)
	} else {
		s = d / (maxC + minC)
	}

	switch maxC {
	case r:
		h = (g - b) / d
		if g < b {
			h += 6
		}
	case g:
		h = (b-r)/d + 2
	case b:
		h = (r-g)/d + 4
	}
	h *= 60
	return
}

func hslToRGB(h, s, l float64) (uint8, uint8, uint8) {
	if s == 0 {
		v := clamp8(l * 255)
		return v, v, v
	}
	var q float64
	if l < 0.5 {
		q = l * (1 + s)
	} else {
		q = l + s - l*s
	}
	p := 2*l - q
	hh := h / 360
	return clamp8(hueToRGB(p, q, hh+1.0/3) * 255),
		clamp8(hueToRGB(p, q, hh) * 255),
		clamp8(hueToRGB(p, q, hh-1.0/3) * 255)
}

func hueToRGB(p, q, t float64) float64 {
	if t < 0 {
		t += 1
	}
	if t > 1 {
		t -= 1
	}
	if t < 1.0/6 {
		return p + (q-p)*6*t
	}
	if t < 0.5 {
		return q
	}
	if t < 2.0/3 {
		return p + (q-p)*(2.0/3-t)*6
	}
	return p
}
