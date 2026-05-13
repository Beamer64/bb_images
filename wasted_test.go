package bb_images

import (
	"bytes"
	"image"
	"image/png"
	"testing"
)

func TestWasted(t *testing.T) {
	src := image.NewRGBA(image.Rect(0, 0, 32, 32))
	out, err := Wasted(src)
	if err != nil {
		t.Fatalf("Wasted: %v", err)
	}
	got, err := png.Decode(bytes.NewReader(out))
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.Bounds().Dx() != 32 || got.Bounds().Dy() != 32 {
		t.Errorf("dims: got %v, want 32x32", got.Bounds())
	}
}
