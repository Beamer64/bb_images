package bb_images

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"testing"
)

func TestMirror(t *testing.T) {
	src := image.NewRGBA(image.Rect(0, 0, 2, 1))
	src.Set(0, 0, color.RGBA{R: 255, A: 255})
	src.Set(1, 0, color.RGBA{B: 255, A: 255})

	out, err := Mirror(src)
	if err != nil {
		t.Fatalf("Mirror: %v", err)
	}
	got, err := png.Decode(bytes.NewReader(out))
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.Bounds().Dx() != 2 || got.Bounds().Dy() != 1 {
		t.Fatalf("dims: got %v, want 2x1", got.Bounds())
	}
	r, _, _, _ := got.At(1, 0).RGBA()
	if r>>8 != 255 {
		t.Errorf("after mirror, red should be at (1,0); got R=%d", r>>8)
	}
}
