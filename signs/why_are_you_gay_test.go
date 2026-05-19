package signs

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"testing"
)

func TestWhyAreYouGay(t *testing.T) {
	interviewee := image.NewRGBA(image.Rect(0, 0, 100, 100))
	interviewer := image.NewRGBA(image.Rect(0, 0, 100, 100))
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			interviewee.Set(x, y, color.RGBA{R: uint8(x * 2), G: uint8(y * 2), B: 64, A: 255})
			interviewer.Set(x, y, color.RGBA{R: uint8(y * 2), G: 64, B: uint8(x * 2), A: 255})
		}
	}
	out, err := WhyAreYouGay(interviewee, interviewer)
	if err != nil {
		t.Fatalf("WhyAreYouGay: %v", err)
	}
	if _, err := png.Decode(bytes.NewReader(out)); err != nil {
		t.Fatalf("decode: %v", err)
	}
}

func TestWhyAreYouGayNilInterviewee(t *testing.T) {
	interviewer := image.NewRGBA(image.Rect(0, 0, 10, 10))
	if _, err := WhyAreYouGay(nil, interviewer); err == nil {
		t.Fatalf("expected error for nil interviewee")
	}
}

func TestWhyAreYouGayNilInterviewer(t *testing.T) {
	interviewee := image.NewRGBA(image.Rect(0, 0, 10, 10))
	if _, err := WhyAreYouGay(interviewee, nil); err == nil {
		t.Fatalf("expected error for nil interviewer")
	}
}
