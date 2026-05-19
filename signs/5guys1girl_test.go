package signs

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"testing"
)

func TestFiveGuysOneGirl(t *testing.T) {
	guys := image.NewRGBA(image.Rect(0, 0, 100, 100))
	girl := image.NewRGBA(image.Rect(0, 0, 100, 100))
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			guys.Set(x, y, color.RGBA{R: uint8(x * 2), G: uint8(y * 2), B: 64, A: 255})
			girl.Set(x, y, color.RGBA{R: uint8(y * 2), G: 64, B: uint8(x * 2), A: 255})
		}
	}
	out, err := FiveGuysOneGirl(guys, girl)
	if err != nil {
		t.Fatalf("FiveGuysOneGirl: %v", err)
	}
	if _, err := png.Decode(bytes.NewReader(out)); err != nil {
		t.Fatalf("decode: %v", err)
	}
}

func TestFiveGuysOneGirlNilGuys(t *testing.T) {
	girl := image.NewRGBA(image.Rect(0, 0, 10, 10))
	if _, err := FiveGuysOneGirl(nil, girl); err == nil {
		t.Fatalf("expected error for nil guys")
	}
}

func TestFiveGuysOneGirlNilGirl(t *testing.T) {
	guys := image.NewRGBA(image.Rect(0, 0, 10, 10))
	if _, err := FiveGuysOneGirl(guys, nil); err == nil {
		t.Fatalf("expected error for nil girl")
	}
}
