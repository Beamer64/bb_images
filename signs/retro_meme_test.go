package signs

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"testing"
)

func newTestAvatar() *image.RGBA {
	avatar := image.NewRGBA(image.Rect(0, 0, 100, 100))
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			avatar.Set(x, y, color.RGBA{R: uint8(x * 2), G: uint8(y * 2), B: 128, A: 255})
		}
	}
	return avatar
}

func TestRetroMemeBoth(t *testing.T) {
	out, err := RetroMeme(newTestAvatar(), "ONE DOES NOT SIMPLY", "WRITE A UNIT TEST")
	if err != nil {
		t.Fatalf("RetroMeme: %v", err)
	}
	if _, err := png.Decode(bytes.NewReader(out)); err != nil {
		t.Fatalf("decode: %v", err)
	}
}

func TestRetroMemeTopOnly(t *testing.T) {
	out, err := RetroMeme(newTestAvatar(), "JUST TOP", "")
	if err != nil {
		t.Fatalf("RetroMeme top-only: %v", err)
	}
	if _, err := png.Decode(bytes.NewReader(out)); err != nil {
		t.Fatalf("decode: %v", err)
	}
}

func TestRetroMemeBottomOnly(t *testing.T) {
	out, err := RetroMeme(newTestAvatar(), "", "JUST BOTTOM")
	if err != nil {
		t.Fatalf("RetroMeme bottom-only: %v", err)
	}
	if _, err := png.Decode(bytes.NewReader(out)); err != nil {
		t.Fatalf("decode: %v", err)
	}
}

func TestRetroMemeNeither(t *testing.T) {
	out, err := RetroMeme(newTestAvatar(), "", "  ")
	if err != nil {
		t.Fatalf("RetroMeme empty: %v", err)
	}
	if _, err := png.Decode(bytes.NewReader(out)); err != nil {
		t.Fatalf("decode: %v", err)
	}
}
