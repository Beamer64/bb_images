package bb_images

import (
	"bytes"
	"image"
	"image/png"
	"testing"
)

func TestTier1Filters(t *testing.T) {
	tests := []struct {
		name string
		fn   func(image.Image) ([]byte, error)
	}{
		{"Invert", Invert},
		{"Blur", Blur},
		{"Sepia", Sepia},
		{"Posterize", Posterize},
		{"Earth", Earth},
		{"Ground", Ground},
		{"Freeze", Freeze},
		{"Night", Night},
		{"Deepfry", Deepfry},
	}

	src := image.NewRGBA(image.Rect(0, 0, 16, 16))
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
			if got.Bounds().Dx() != 16 || got.Bounds().Dy() != 16 {
				t.Errorf("%s dims: got %v, want 16x16", tt.name, got.Bounds())
			}
		})
	}
}
