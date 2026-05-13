package bb_images

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"testing"
)

func TestFlip(t *testing.T) {
	src := image.NewRGBA(image.Rect(0, 0, 1, 2))
	src.Set(0, 0, color.RGBA{R: 255, A: 255})
	src.Set(0, 1, color.RGBA{B: 255, A: 255})

	out, err := Flip(src)
	if err != nil {
		t.Fatalf("Flip: %v", err)
	}
	got, err := png.Decode(bytes.NewReader(out))
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.Bounds().Dx() != 1 || got.Bounds().Dy() != 2 {
		t.Fatalf("dims: got %v, want 1x2", got.Bounds())
	}
	r, _, _, _ := got.At(0, 1).RGBA()
	if r>>8 != 255 {
		t.Errorf("after flip, red should be at (0,1); got R=%d", r>>8)
	}
}
