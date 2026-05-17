// Package bb_images is the top-level module doc for the image-manipulation
// library. The actual filters live in topic sub-packages — import only the
// one(s) you need:
//
//	overlays  — drop an embedded PNG over the source (america, communism, jail, wasted)
//	signs     — text-on-template composites with perspective warp (change-my-mind, ...)
//	color     — pixel-level transforms (invert, blur, sepia, posterize, tints, deepfry)
//	edges     — edge-based effects (sobel, sketch, charcoal, hog)
//	animated  — GIF outputs (triggered, expand, shake, glitch, tv-static, rain, rainbow, spin, glitch-static)
//	spatial   — geometric transforms (mirror, flip, pixelate, colors, mosaic, swirl, triangle, magik, stringify)
//	special   — misc one-offs (ascii, paint, rgb, burn)
//
// All filters take an image.Image and return encoded bytes (PNG for static
// filters, GIF for animated). No external API calls; assets are embedded
// at build time via //go:embed.
package bb_images
