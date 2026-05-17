package signs

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"strings"
	"testing"
)

func TestYouTube(t *testing.T) {
	avatar := image.NewRGBA(image.Rect(0, 0, 64, 64))
	for y := 0; y < 64; y++ {
		for x := 0; x < 64; x++ {
			avatar.Set(x, y, color.RGBA{R: uint8(x * 4), G: uint8(y * 4), B: 128, A: 255})
		}
	}
	out, err := YouTube(avatar, "Test User", "First!")
	if err != nil {
		t.Fatalf("YouTube: %v", err)
	}
	got, err := png.Decode(bytes.NewReader(out))
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.Bounds().Dx() == 0 || got.Bounds().Dy() == 0 {
		t.Errorf("zero dims: %v", got.Bounds())
	}
}

func TestYouTubeEmpty(t *testing.T) {
	avatar := image.NewRGBA(image.Rect(0, 0, 16, 16))
	if _, err := YouTube(avatar, "", "comment"); err == nil {
		t.Error("expected error for empty username, got nil")
	}
	if _, err := YouTube(avatar, "user", "  "); err == nil {
		t.Error("expected error for empty comment, got nil")
	}
}

func TestRandomYouTubeTime(t *testing.T) {
	for i := 0; i < 20; i++ {
		s := randomYouTubeTime()
		if s == "" || !strings.HasSuffix(s, "ago") {
			t.Errorf("expected non-empty string ending with 'ago', got %q", s)
		}
	}
}
