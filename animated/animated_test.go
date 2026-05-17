package animated

import (
	"bytes"
	"image"
	"image/color"
	"image/gif"
	"testing"
)

func TestAnimatedFilters(t *testing.T) {
	tests := []struct {
		name string
		fn   func(image.Image) ([]byte, error)
	}{
		{"Shake", Shake},
		{"Glitch", Glitch},
		{"TvStatic", TvStatic},
		{"Spin", Spin},
		{"Rainbow", Rainbow},
		{"Rain", Rain},
		{"GlitchStatic", GlitchStatic},
	}

	src := image.NewRGBA(image.Rect(0, 0, 64, 64))
	for y := 0; y < 64; y++ {
		for x := 0; x < 64; x++ {
			src.Set(x, y, color.RGBA{R: uint8(x * 4), G: uint8(y * 4), B: 128, A: 255})
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := tt.fn(src)
			if err != nil {
				t.Fatalf("%s: %v", tt.name, err)
			}
			g, err := gif.DecodeAll(bytes.NewReader(out))
			if err != nil {
				t.Fatalf("decode: %v", err)
			}
			if len(g.Image) < 2 {
				t.Errorf("%s: expected multi-frame GIF, got %d frames", tt.name, len(g.Image))
			}
		})
	}
}
