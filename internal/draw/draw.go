// Package draw holds shared helpers for the bb_images filter packages:
// decoding incoming bytes and encoding output as PNG or GIF.
package draw

import (
	"bytes"
	"image"
	"image/gif"
	_ "image/jpeg" // register JPEG decoder
	"image/png"
)

// Decode parses src as an image (PNG, JPEG, or any registered format).
func Decode(src []byte) (image.Image, error) {
	img, _, err := image.Decode(bytes.NewReader(src))
	return img, err
}

// EncodePNG returns img as PNG bytes.
func EncodePNG(img image.Image) ([]byte, error) {
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// EncodeGIF returns g as GIF bytes.
func EncodeGIF(g *gif.GIF) ([]byte, error) {
	var buf bytes.Buffer
	if err := gif.EncodeAll(&buf, g); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
