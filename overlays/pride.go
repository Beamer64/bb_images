package overlays

import (
	"embed"
	"errors"
	"fmt"
	"image"
	"sort"
	"strings"

	"github.com/Beamer64/bb_images/internal/draw"
	"github.com/disintegration/imaging"
)

//go:embed res/pride/*.png
var prideFlagsFS embed.FS

// prideOpacity controls how visible the flag is over the source. 0.5
// reads cleanly without obscuring the face.
const prideOpacity = 0.5

// PrideFlags returns the list of available pride flag names (the
// filename without the .png extension), sorted alphabetically. Useful
// for keeping spec-side choice lists in sync with what's actually
// embedded.
func PrideFlags() []string {
	entries, err := prideFlagsFS.ReadDir("res/pride")
	if err != nil {
		return nil
	}
	names := make([]string, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if !strings.HasSuffix(name, ".png") {
			continue
		}
		names = append(names, strings.TrimSuffix(name, ".png"))
	}
	sort.Strings(names)
	return names
}

// Pride composites the chosen pride flag onto src at prideOpacity.
// flagName matches the basename of one of the PNGs in res/pride/ (e.g.
// "gay", "trans", "progress"). Input is normalized to lowercase, and
// only [a-z0-9] characters are accepted — anything else returns an
// error to prevent the lookup from being tricked into resolving a path
// outside the embedded flag directory.
func Pride(src image.Image, flagName string) ([]byte, error) {
	name := strings.ToLower(strings.TrimSpace(flagName))
	if name == "" {
		return nil, errors.New("pride: empty flag name")
	}
	for _, r := range name {
		if !((r >= 'a' && r <= 'z') || (r >= '0' && r <= '9')) {
			return nil, fmt.Errorf("pride: invalid flag name %q", flagName)
		}
	}

	flagBytes, err := prideFlagsFS.ReadFile("res/pride/" + name + ".png")
	if err != nil {
		return nil, fmt.Errorf("pride: unknown flag %q: %w", flagName, err)
	}

	flag, err := draw.Decode(flagBytes)
	if err != nil {
		return nil, err
	}

	b := src.Bounds()
	resized := imaging.Fill(flag, b.Dx(), b.Dy(), imaging.Center, imaging.Linear)
	out := imaging.Overlay(src, resized, image.Point{}, prideOpacity)
	return draw.EncodePNG(out)
}
