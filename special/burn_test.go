package special

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"testing"
)

func TestBurn(t *testing.T) {
	src := image.NewRGBA(image.Rect(0, 0, 64, 64))
	for y := 0; y < 64; y++ {
		for x := 0; x < 64; x++ {
			src.Set(x, y, color.RGBA{R: uint8(x * 4), G: uint8(y * 4), B: 128, A: 255})
		}
	}
	out, err := Burn(src)
	if err != nil {
		t.Fatalf("Burn: %v", err)
	}
	got, err := png.Decode(bytes.NewReader(out))
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.Bounds().Dx() != 64 || got.Bounds().Dy() != 64 {
		t.Errorf("dims: got %v, want 64x64", got.Bounds())
	}
}
