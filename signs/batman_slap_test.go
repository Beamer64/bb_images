package signs

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"testing"
)

func TestBatmanSlap(t *testing.T) {
	batman := image.NewRGBA(image.Rect(0, 0, 100, 100))
	robin := image.NewRGBA(image.Rect(0, 0, 100, 100))
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			batman.Set(x, y, color.RGBA{R: uint8(x * 2), G: uint8(y * 2), B: 64, A: 255})
			robin.Set(x, y, color.RGBA{R: uint8(y * 2), G: 64, B: uint8(x * 2), A: 255})
		}
	}
	out, err := BatmanSlap(batman, robin)
	if err != nil {
		t.Fatalf("BatmanSlap: %v", err)
	}
	if _, err := png.Decode(bytes.NewReader(out)); err != nil {
		t.Fatalf("decode: %v", err)
	}
}

func TestBatmanSlapNilBatman(t *testing.T) {
	robin := image.NewRGBA(image.Rect(0, 0, 10, 10))
	if _, err := BatmanSlap(nil, robin); err == nil {
		t.Fatalf("expected error for nil batman")
	}
}

func TestBatmanSlapNilRobin(t *testing.T) {
	batman := image.NewRGBA(image.Rect(0, 0, 10, 10))
	if _, err := BatmanSlap(batman, nil); err == nil {
		t.Fatalf("expected error for nil robin")
	}
}
