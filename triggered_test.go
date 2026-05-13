package bb_images

import (
	"bytes"
	"image"
	"image/gif"
	"testing"
)

func TestTriggered(t *testing.T) {
	src := image.NewRGBA(image.Rect(0, 0, 64, 64))
	out, err := Triggered(src)
	if err != nil {
		t.Fatalf("Triggered: %v", err)
	}
	g, err := gif.DecodeAll(bytes.NewReader(out))
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(g.Image) < 2 {
		t.Errorf("expected multi-frame GIF, got %d frames", len(g.Image))
	}
	if g.Image[0].Bounds().Dx() != 64 || g.Image[0].Bounds().Dy() != 64 {
		t.Errorf("frame dims: got %v, want 64x64", g.Image[0].Bounds())
	}
}
