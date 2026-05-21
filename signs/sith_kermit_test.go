package signs

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"testing"
)

func TestSithKermit(t *testing.T) {
	sith := image.NewRGBA(image.Rect(0, 0, 100, 100))
	kermit := image.NewRGBA(image.Rect(0, 0, 100, 100))
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			sith.Set(x, y, color.RGBA{R: uint8(x * 2), G: uint8(y * 2), B: 64, A: 255})
			kermit.Set(x, y, color.RGBA{R: uint8(y * 2), G: 64, B: uint8(x * 2), A: 255})
		}
	}
	out, err := SithKermit(sith, kermit)
	if err != nil {
		t.Fatalf("SithKermit: %v", err)
	}
	if _, err := png.Decode(bytes.NewReader(out)); err != nil {
		t.Fatalf("decode: %v", err)
	}
}

func TestSithKermitNilSith(t *testing.T) {
	kermit := image.NewRGBA(image.Rect(0, 0, 10, 10))
	if _, err := SithKermit(nil, kermit); err == nil {
		t.Fatalf("expected error for nil sith")
	}
}

func TestSithKermitNilKermit(t *testing.T) {
	sith := image.NewRGBA(image.Rect(0, 0, 10, 10))
	if _, err := SithKermit(sith, nil); err == nil {
		t.Fatalf("expected error for nil kermit")
	}
}
