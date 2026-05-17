package signs

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"regexp"
	"testing"
)

func TestDiscord(t *testing.T) {
	avatar := image.NewRGBA(image.Rect(0, 0, 64, 64))
	for y := 0; y < 64; y++ {
		for x := 0; x < 64; x++ {
			avatar.Set(x, y, color.RGBA{R: uint8(x * 4), G: uint8(y * 4), B: 128, A: 255})
		}
	}
	out, err := Discord(avatar, "Test User", "Hello from a unit test")
	if err != nil {
		t.Fatalf("Discord: %v", err)
	}
	got, err := png.Decode(bytes.NewReader(out))
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.Bounds().Dx() == 0 || got.Bounds().Dy() == 0 {
		t.Errorf("zero dims: %v", got.Bounds())
	}
}

func TestDiscordEmpty(t *testing.T) {
	avatar := image.NewRGBA(image.Rect(0, 0, 16, 16))
	if _, err := Discord(avatar, "", "hi"); err == nil {
		t.Error("expected error for empty display name, got nil")
	}
	if _, err := Discord(avatar, "user", "  "); err == nil {
		t.Error("expected error for empty message, got nil")
	}
}

func TestRandomDiscordTimestamp(t *testing.T) {
	pattern := regexp.MustCompile(`^\d{2}/\d{2}/\d{4} \d{2}:\d{2} (AM|PM)$`)
	for i := 0; i < 20; i++ {
		s := randomDiscordTimestamp()
		if !pattern.MatchString(s) {
			t.Errorf("timestamp %q did not match MM/DD/YYYY HH:MM AM/PM", s)
		}
	}
}
