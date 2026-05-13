package bb_images

import (
	"bytes"
	"image"
	"image/color"
	"image/gif"
	"image/png"
	"testing"
)

func TestTier4StaticFilters(t *testing.T) {
	tests := []struct {
		name string
		fn   func(image.Image) ([]byte, error)
	}{
		{"Ascii", Ascii},
		{"Paint", Paint},
	}

	src := image.NewRGBA(image.Rect(0, 0, 32, 32))
	for y := 0; y < 32; y++ {
		for x := 0; x < 32; x++ {
			src.Set(x, y, color.RGBA{R: uint8(x * 8), G: uint8(y * 8), B: 128, A: 255})
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := tt.fn(src)
			if err != nil {
				t.Fatalf("%s: %v", tt.name, err)
			}
			got, err := png.Decode(bytes.NewReader(out))
			if err != nil {
				t.Fatalf("decode: %v", err)
			}
			if got.Bounds().Dx() == 0 || got.Bounds().Dy() == 0 {
				t.Errorf("%s: zero dims %v", tt.name, got.Bounds())
			}
		})
	}
}

func TestTier4AnimatedFilters(t *testing.T) {
	tests := []struct {
		name string
		fn   func(image.Image) ([]byte, error)
	}{
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
