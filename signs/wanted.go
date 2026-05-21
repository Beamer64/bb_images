package signs

import (
	_ "embed"
	"errors"
	"image"
	"image/color"
	imgdraw "image/draw"

	"github.com/Beamer64/bb_images/internal/draw"
	"github.com/Beamer64/bb_images/internal/templates"
	"github.com/disintegration/imaging"
)

//go:embed res/wanted.png
var wantedBytes []byte

//go:embed res/markers/wanted.marker.png
var wantedMarkersBytes []byte

var wantedAvatarColor = color.RGBA{G: 255, A: 255}

// wantedAvatarOpacity controls how much of the template's paper texture
// (stains, wear, etc.) shows through the placed avatar.
//   - 1.0 = fully opaque avatar, no stains visible
//   - 0.0 = avatar invisible
//
// ~0.78 leaves the face clearly readable while still letting the weathered
// poster come through.
const wantedAvatarOpacity = 0.65

// Wanted composites a sepia-toned, slightly transparent copy of the avatar
// onto each green marker region of the Wild West wanted-poster template.
// The sepia warms the avatar to match the parchment, and the partial
// transparency lets the template's stains and ink wear show through for a
// faded "found pinned to a saloon wall" look.
func Wanted(avatar image.Image) ([]byte, error) {
	if avatar == nil {
		return nil, errors.New("wanted: nil avatar")
	}

	template, err := draw.Decode(wantedBytes)
	if err != nil {
		return nil, err
	}
	markersImg, err := draw.Decode(wantedMarkersBytes)
	if err != nil {
		return nil, err
	}

	tb := template.Bounds()
	out := image.NewRGBA(image.Rect(0, 0, tb.Dx(), tb.Dy()))
	imgdraw.Draw(out, out.Bounds(), template, tb.Min, imgdraw.Src)

	regions := templates.ConnectedRegions(markersImg, wantedAvatarColor, 10)
	if len(regions) == 0 {
		return nil, errors.New("wanted: no green marker regions found")
	}
	for _, r := range regions {
		placeSepiaMaskedImage(out, avatar, markersImg, wantedAvatarColor, 10, r, wantedAvatarOpacity)
	}

	return draw.EncodePNG(out)
}

// placeSepiaMaskedImage resizes src to the bounding box, then for every
// pixel inside bounds where the markers image matches markerColor it:
//  1. reads the resized src pixel,
//  2. transforms it through the standard sepia matrix,
//  3. alpha-blends the result with the existing pixel in out using α (the
//     template content that was already drawn shows through at 1-α).
//
// This is the placeMaskedImage variant used by the weathered/parchment-
// style templates (currently just `wanted`). If a second command needs the
// same look, this helper can be hoisted into a shared file.
func placeSepiaMaskedImage(
	out *image.RGBA,
	src image.Image,
	markers image.Image,
	markerColor color.RGBA,
	tolerance int,
	bounds image.Rectangle,
	alpha float64,
) {
	sized := imaging.Resize(src, bounds.Dx(), bounds.Dy(), imaging.Lanczos)
	invA := 1.0 - alpha
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			mr, mg, mb, _ := markers.At(x, y).RGBA()
			if absDiff(int(mr>>8), int(markerColor.R)) > tolerance ||
				absDiff(int(mg>>8), int(markerColor.G)) > tolerance ||
				absDiff(int(mb>>8), int(markerColor.B)) > tolerance {
				continue
			}

			sr, sg, sb, _ := sized.At(x-bounds.Min.X, y-bounds.Min.Y).RGBA()
			srf, sgf, sbf := float64(sr>>8), float64(sg>>8), float64(sb>>8)

			// Classic sepia matrix.
			sepR := 0.393*srf + 0.769*sgf + 0.189*sbf
			sepG := 0.349*srf + 0.686*sgf + 0.168*sbf
			sepB := 0.272*srf + 0.534*sgf + 0.131*sbf
			if sepR > 255 {
				sepR = 255
			}
			if sepG > 255 {
				sepG = 255
			}
			if sepB > 255 {
				sepB = 255
			}

			// Blend with the template content currently in out at this pixel.
			tr, tg, tb, _ := out.At(x, y).RGBA()
			trf, tgf, tbf := float64(tr>>8), float64(tg>>8), float64(tb>>8)

			out.Set(
				x, y, color.RGBA{
					R: uint8(alpha*sepR + invA*trf),
					G: uint8(alpha*sepG + invA*tgf),
					B: uint8(alpha*sepB + invA*tbf),
					A: 255,
				},
			)
		}
	}
}
