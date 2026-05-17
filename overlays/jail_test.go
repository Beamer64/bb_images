package overlays

import (
	"bytes"
	"image"
	"image/png"
	"testing"
)

func TestJail(t *testing.T) {
	src := image.NewRGBA(image.Rect(0, 0, 32, 32))
	out, err := Jail(src)
	if err != nil {
		t.Fatalf("Jail: %v", err)
	}
	got, err := png.Decode(bytes.NewReader(out))
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.Bounds().Dx() != 32 || got.Bounds().Dy() != 32 {
		t.Errorf("dims: got %v, want 32x32", got.Bounds())
	}
}
