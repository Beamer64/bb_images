// Package templates detects colored marker regions in template images.
// Each template ships a parallel "markers" image where editable regions
// are painted in distinct colors. Detect scans the markers image and
// returns, per registered color, the bounding box plus the four
// directional extreme points — for a tilted quadrilateral those are the
// four corners, and the caller maps them onto TL/TR/BR/BL based on the
// template's orientation.
package templates

import (
	"fmt"
	"image"
	"image/color"
)

// Region describes a detected color marker. Bounds is the axis-aligned
// bounding rectangle of all matching pixels. Top/Right/Bottom/Left are
// the four directional extremes (smallest Y / largest X / largest Y /
// smallest X). Count is the number of matching pixels — useful for
// sanity-checking that a region was actually found.
type Region struct {
	Bounds image.Rectangle
	Top    image.Point
	Right  image.Point
	Bottom image.Point
	Left   image.Point
	Count  int
}

// Detect scans markers for pixels matching each color in roles (within
// tolerance per RGB channel) and returns one Region per role. Returns
// an error if any role's color does not appear in the image.
func Detect(markers image.Image, roles map[color.RGBA]string, tolerance int) (map[string]Region, error) {
	type acc struct {
		minX, maxX, minY, maxY           int
		topPt, rightPt, bottomPt, leftPt image.Point
		count                            int
	}
	accs := make(map[string]*acc, len(roles))
	for _, role := range roles {
		accs[role] = &acc{
			minX: int(^uint(0) >> 1),
			minY: int(^uint(0) >> 1),
			maxX: -1 << 31,
			maxY: -1 << 31,
		}
	}

	b := markers.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			r16, g16, b16, _ := markers.At(x, y).RGBA()
			r8, g8, b8 := int(r16>>8), int(g16>>8), int(b16>>8)
			for col, role := range roles {
				if absInt(r8-int(col.R)) <= tolerance &&
					absInt(g8-int(col.G)) <= tolerance &&
					absInt(b8-int(col.B)) <= tolerance {
					a := accs[role]
					if x < a.minX {
						a.minX = x
						a.leftPt = image.Point{X: x, Y: y}
					}
					if x > a.maxX {
						a.maxX = x
						a.rightPt = image.Point{X: x, Y: y}
					}
					if y < a.minY {
						a.minY = y
						a.topPt = image.Point{X: x, Y: y}
					}
					if y > a.maxY {
						a.maxY = y
						a.bottomPt = image.Point{X: x, Y: y}
					}
					a.count++
					break // a pixel matches at most one role
				}
			}
		}
	}

	regions := make(map[string]Region, len(roles))
	for _, role := range roles {
		a := accs[role]
		if a.count == 0 {
			return nil, fmt.Errorf("templates: no pixels matching role %q", role)
		}
		regions[role] = Region{
			Bounds: image.Rect(a.minX, a.minY, a.maxX+1, a.maxY+1),
			Top:    a.topPt,
			Right:  a.rightPt,
			Bottom: a.bottomPt,
			Left:   a.leftPt,
			Count:  a.count,
		}
	}
	return regions, nil
}

func absInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
