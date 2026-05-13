package bb_images

import (
	"bytes"
	"image"
	"image/png"
	"testing"
)

func TestPixelate(t *testing.T) {
	src := image.NewRGBA(image.Rect(0, 0, 16, 16))
	out, err := Pixelate(src, 4)
	if err != nil {
		t.Fatalf("Pixelate: %v", err)
	}
	got, err := png.Decode(bytes.NewReader(out))
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.Bounds().Dx() != 16 || got.Bounds().Dy() != 16 {
		t.Errorf("dims: got %v, want 16x16", got.Bounds())
	}
}

func TestPixelateInvalidFactor(t *testing.T) {
	src := image.NewRGBA(image.Rect(0, 0, 16, 16))
	if _, err := Pixelate(src, 1); err == nil {
		t.Error("expected error for factor < 2, got nil")
	}
}

func TestPixelateImageSmallerThanFactor(t *testing.T) {
	src := image.NewRGBA(image.Rect(0, 0, 4, 4))
	if _, err := Pixelate(src, 8); err == nil {
		t.Error("expected error when image smaller than factor, got nil")
	}
}
