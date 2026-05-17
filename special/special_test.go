package special

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"testing"
)

func TestSpecialFilters(t *testing.T) {
	tests := []struct {
		name string
		fn   func(image.Image) ([]byte, error)
	}{
		{"Ascii", Ascii},
		{"Paint", Paint},
		{"RGB", RGB},
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
