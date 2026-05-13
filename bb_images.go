// Package bb_images provides in-process image manipulation filters.
// Each filter is a top-level function (Pixelate, Mirror, Flip, etc.)
// taking an image.Image and returning the encoded result (PNG for
// static filters, GIF for animated). No external API calls.
package bb_images

// (Tier 1 filters live in invert.go, blur.go, sepia.go, posterize.go,
// earth.go, ground.go, freeze.go, night.go, deepfry.go — all share
// helpers in tint.go.)
