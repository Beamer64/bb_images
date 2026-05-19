package signs

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"testing"
)

func TestThanksObama(t *testing.T) {
	avatar := image.NewRGBA(image.Rect(0, 0, 100, 100))
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			avatar.Set(x, y, color.RGBA{R: uint8(x * 2), G: uint8(y * 2), B: 128, A: 255})
		}
	}
	out, err := ThanksObama(avatar)
	if err != nil {
		t.Fatalf("ThanksObama: %v", err)
	}
	if _, err := png.Decode(bytes.NewReader(out)); err != nil {
		t.Fatalf("decode: %v", err)
	}
}

func TestThanksObamaNilAvatar(t *testing.T) {
	if _, err := ThanksObama(nil); err == nil {
		t.Fatalf("expected error for nil avatar")
	}
}
