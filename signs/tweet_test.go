package signs

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"testing"
)

func TestTweet(t *testing.T) {
	avatar := image.NewRGBA(image.Rect(0, 0, 64, 64))
	for y := 0; y < 64; y++ {
		for x := 0; x < 64; x++ {
			avatar.Set(x, y, color.RGBA{R: uint8(x * 4), G: uint8(y * 4), B: 128, A: 255})
		}
	}
	out, err := Tweet(avatar, "Test User", "testuser", "Hello world from a unit test")
	if err != nil {
		t.Fatalf("Tweet: %v", err)
	}
	got, err := png.Decode(bytes.NewReader(out))
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.Bounds().Dx() == 0 || got.Bounds().Dy() == 0 {
		t.Errorf("zero dims: %v", got.Bounds())
	}
}

func TestTweetEmpty(t *testing.T) {
	avatar := image.NewRGBA(image.Rect(0, 0, 16, 16))
	if _, err := Tweet(avatar, "", "user", "hi"); err == nil {
		t.Error("expected error for empty display name, got nil")
	}
	if _, err := Tweet(avatar, "Test", "", "hi"); err == nil {
		t.Error("expected error for empty username, got nil")
	}
	if _, err := Tweet(avatar, "Test", "user", "  "); err == nil {
		t.Error("expected error for empty tweet text, got nil")
	}
}
