package signs

import (
	"bytes"
	"image/png"
	"testing"
)

func TestChangeMyMind(t *testing.T) {
	out, err := ChangeMyMind("Climate change is real")
	if err != nil {
		t.Fatalf("ChangeMyMind: %v", err)
	}
	got, err := png.Decode(bytes.NewReader(out))
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.Bounds().Dx() == 0 || got.Bounds().Dy() == 0 {
		t.Errorf("zero dims: %v", got.Bounds())
	}
}

func TestChangeMyMindEmpty(t *testing.T) {
	if _, err := ChangeMyMind("   "); err == nil {
		t.Error("expected error for empty text, got nil")
	}
}

func TestChangeMyMindLong(t *testing.T) {
	long := "The implementation of perspective-corrected text rendering on a tilted sign requires solving an eight-parameter homography from four corner correspondences."
	out, err := ChangeMyMind(long)
	if err != nil {
		t.Fatalf("ChangeMyMind long: %v", err)
	}
	if _, err := png.Decode(bytes.NewReader(out)); err != nil {
		t.Fatalf("decode long: %v", err)
	}
}
