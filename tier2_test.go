package bb_images

import (
	"bytes"
	"image"
	"image/gif"
	"testing"
)

func TestTier2Filters(t *testing.T) {
	tests := []struct {
		name string
		fn   func(image.Image) ([]byte, error)
	}{
		{"Shake", Shake},
		{"Glitch", Glitch},
		{"TvStatic", TvStatic},
	}

	src := image.NewRGBA(image.Rect(0, 0, 64, 64))
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
			if g.Image[0].Bounds().Dx() != 64 || g.Image[0].Bounds().Dy() != 64 {
				t.Errorf("%s frame dims: got %v, want 64x64", tt.name, g.Image[0].Bounds())
			}
		})
	}
}
