package overlays

import (
	"bytes"
	"image"
	"image/png"
	"testing"
)

func TestPride(t *testing.T) {
	src := image.NewRGBA(image.Rect(0, 0, 32, 32))
	out, err := Pride(src, "gay")
	if err != nil {
		t.Fatalf("Pride: %v", err)
	}
	got, err := png.Decode(bytes.NewReader(out))
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.Bounds().Dx() != 32 || got.Bounds().Dy() != 32 {
		t.Errorf("dims: got %v, want 32x32", got.Bounds())
	}
}

func TestPrideUnknownFlag(t *testing.T) {
	src := image.NewRGBA(image.Rect(0, 0, 32, 32))
	if _, err := Pride(src, "definitely-not-a-real-flag"); err == nil {
		t.Fatalf("expected error for unknown flag")
	}
}

func TestPrideRejectsTraversal(t *testing.T) {
	src := image.NewRGBA(image.Rect(0, 0, 32, 32))
	// Path-traversal attempt — must be rejected before the FS lookup.
	if _, err := Pride(src, "../jail"); err == nil {
		t.Fatalf("expected error for path-traversal flag name")
	}
}

func TestPrideEmptyFlag(t *testing.T) {
	src := image.NewRGBA(image.Rect(0, 0, 32, 32))
	if _, err := Pride(src, ""); err == nil {
		t.Fatalf("expected error for empty flag name")
	}
}

func TestPrideFlagsListsAll(t *testing.T) {
	flags := PrideFlags()
	if len(flags) == 0 {
		t.Fatal("PrideFlags returned no entries; res/pride/*.png missing?")
	}
	// Must contain at least one known flag so a typo in the FS path
	// would fail the test loudly.
	found := false
	for _, f := range flags {
		if f == "gay" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected 'gay' in PrideFlags(), got %v", flags)
	}
}
