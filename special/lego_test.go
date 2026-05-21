package special

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"testing"
)

func TestLego(t *testing.T) {
	src := image.NewRGBA(image.Rect(0, 0, 100, 100))
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			src.Set(x, y, color.RGBA{R: uint8(x * 2), G: uint8(y * 2), B: 128, A: 255})
		}
	}
	out, err := Lego(src)
	if err != nil {
		t.Fatalf("Lego: %v", err)
	}
	img, err := png.Decode(bytes.NewReader(out))
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	if img.Bounds().Dx() != 100 || img.Bounds().Dy() != 100 {
		t.Fatalf("expected 100x100, got %dx%d", img.Bounds().Dx(), img.Bounds().Dy())
	}
}

// TestLegoEdgeBricks exercises the partial-brick code path with
// dimensions that aren't a clean multiple of legoBrickSize. The function
// should produce a valid PNG and not panic on the clipped edge cells.
func TestLegoEdgeBricks(t *testing.T) {
	src := image.NewRGBA(image.Rect(0, 0, 73, 91))
	for y := 0; y < 91; y++ {
		for x := 0; x < 73; x++ {
			src.Set(x, y, color.RGBA{R: 200, G: 100, B: 50, A: 255})
		}
	}
	out, err := Lego(src)
	if err != nil {
		t.Fatalf("Lego: %v", err)
	}
	if _, err := png.Decode(bytes.NewReader(out)); err != nil {
		t.Fatalf("decode: %v", err)
	}
}
