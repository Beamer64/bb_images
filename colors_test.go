package bb_images

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"testing"
)

func TestColors(t *testing.T) {
	src := image.NewRGBA(image.Rect(0, 0, 4, 4))
	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			if x < 3 {
				src.Set(x, y, color.RGBA{R: 255, A: 255})
			} else {
				src.Set(x, y, color.RGBA{B: 255, A: 255})
			}
		}
	}

	out, err := Colors(src, 1)
	if err != nil {
		t.Fatalf("Colors: %v", err)
	}
	got, err := png.Decode(bytes.NewReader(out))
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.Bounds().Dx() != colorsSwatchWidth || got.Bounds().Dy() != colorsSwatchHeight {
		t.Fatalf("dims: got %v, want %dx%d", got.Bounds(), colorsSwatchWidth, colorsSwatchHeight)
	}
	r, _, _, _ := got.At(colorsSwatchWidth/2, colorsSwatchHeight/2).RGBA()
	if r>>8 < 200 {
		t.Errorf("expected dominant stripe to be red; got R=%d", r>>8)
	}
}

func TestColorsInvalidK(t *testing.T) {
	src := image.NewRGBA(image.Rect(0, 0, 4, 4))
	src.Set(0, 0, color.RGBA{R: 255, A: 255})
	if _, err := Colors(src, 0); err == nil {
		t.Error("expected error for k < 1, got nil")
	}
}

func TestColorsAllTransparent(t *testing.T) {
	src := image.NewRGBA(image.Rect(0, 0, 4, 4))
	if _, err := Colors(src, 1); err == nil {
		t.Error("expected error when image has no opaque pixels, got nil")
	}
}
