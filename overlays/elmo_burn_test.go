package overlays

import (
	"bytes"
	"image"
	"image/gif"
	"testing"
)

func TestElmoBurn(t *testing.T) {
	src := image.NewRGBA(image.Rect(0, 0, 32, 32))
	out, err := ElmoBurn(src)
	if err != nil {
		t.Fatalf("ElmoBurn: %v", err)
	}
	got, err := gif.DecodeAll(bytes.NewReader(out))
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(got.Image) == 0 {
		t.Fatalf("expected at least one frame, got 0")
	}
	if b := got.Image[0].Bounds(); b.Dx() != 32 || b.Dy() != 32 {
		t.Errorf("frame 0 dims: got %v, want 32x32", b)
	}
}
