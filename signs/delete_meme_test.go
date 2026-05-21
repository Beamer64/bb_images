package signs

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"testing"
)

func TestDeleteMeme(t *testing.T) {
	meme := image.NewRGBA(image.Rect(0, 0, 100, 100))
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			meme.Set(x, y, color.RGBA{R: uint8(x * 2), G: uint8(y * 2), B: 128, A: 255})
		}
	}
	out, err := DeleteMeme(meme)
	if err != nil {
		t.Fatalf("DeleteMeme: %v", err)
	}
	if _, err := png.Decode(bytes.NewReader(out)); err != nil {
		t.Fatalf("decode: %v", err)
	}
}

func TestDeleteMemeNilImage(t *testing.T) {
	if _, err := DeleteMeme(nil); err == nil {
		t.Fatalf("expected error for nil image")
	}
}
