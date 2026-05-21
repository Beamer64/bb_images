// Package draw holds shared helpers for the bb_images filter packages:
// decoding incoming bytes and encoding output as PNG or GIF.
package draw

import (
	"bytes"
	"image"
	"image/color"
	imgdraw "image/draw"
	"image/gif"
	_ "image/jpeg" // register JPEG decoder
	"image/png"
	"sync"
)

// Decode parses src as an image (PNG, JPEG, or any registered format).
func Decode(src []byte) (image.Image, error) {
	img, _, err := image.Decode(bytes.NewReader(src))
	return img, err
}

// DecodeGIF parses src as an animated GIF, returning all frames + per-
// frame delays/disposal info.
func DecodeGIF(src []byte) (*gif.GIF, error) {
	return gif.DecodeAll(bytes.NewReader(src))
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

// RenderFrames decodes each frame of g as a fully-rendered RGBA image —
// i.e. the actual visible content the viewer would see at that point in
// the animation, not the optimized delta stored in g.Image. Handles all
// three disposal methods (None/Background/Previous) so partial-frame
// deltas don't accidentally "stick" or "smear" into subsequent frames.
//
// Returns the rendered frames in playback order alongside the original
// per-frame delays (hundredths of a second), suitable for re-encoding
// into a new gif.GIF after compositing onto each frame.
func RenderFrames(g *gif.GIF) ([]*image.RGBA, []int) {
	bounds := image.Rect(0, 0, g.Config.Width, g.Config.Height)
	canvas := image.NewRGBA(bounds)
	transparent := image.NewUniform(color.RGBA{})
	var savedCanvas *image.RGBA // snapshot for DisposalPrevious restore

	frames := make([]*image.RGBA, 0, len(g.Image))
	for i, frame := range g.Image {
		fb := frame.Bounds()

		// If this frame uses DisposalPrevious, snapshot the canvas now so
		// we can restore it after drawing the frame, before moving on.
		if g.Disposal[i] == gif.DisposalPrevious {
			if savedCanvas == nil {
				savedCanvas = image.NewRGBA(canvas.Bounds())
			}
			imgdraw.Draw(savedCanvas, canvas.Bounds(), canvas, image.Point{}, imgdraw.Src)
		}

		// Draw the frame's delta onto the running canvas. Paletted GIFs
		// can have a transparent index, so Over respects that.
		imgdraw.Draw(canvas, fb, frame, fb.Min, imgdraw.Over)

		// Snapshot the canvas as the rendered frame for the caller.
		rendered := image.NewRGBA(canvas.Bounds())
		imgdraw.Draw(rendered, canvas.Bounds(), canvas, image.Point{}, imgdraw.Src)
		frames = append(frames, rendered)

		// Apply this frame's disposal to set up the canvas for the next
		// iteration.
		switch g.Disposal[i] {
		case gif.DisposalBackground:
			imgdraw.Draw(canvas, fb, transparent, image.Point{}, imgdraw.Src)
		case gif.DisposalPrevious:
			if savedCanvas != nil {
				imgdraw.Draw(canvas, canvas.Bounds(), savedCanvas, image.Point{}, imgdraw.Src)
			}
		}
		// DisposalNone (and the unspecified default, 0): leave canvas
		// as-is — the next frame is drawn on top of this one.
	}

	return frames, g.Delay
}

// LazyFrames wraps a static GIF byte slice (typically an //go:embed) and
// decodes it into rendered RGBA frames on first access. Subsequent calls
// return the cached frames, so a bot command that hits the same template
// on every invocation pays the decode cost exactly once per process.
//
// The cached frames are immutable from the cache's perspective — callers
// must not mutate them. (AnimateOverGIF below treats them read-only.)
type LazyFrames struct {
	bytes []byte

	once    sync.Once
	frames  []*image.RGBA
	delays  []int
	loop    int
	loadErr error
}

// NewLazyFrames returns a cache wrapper around src. The decode itself is
// deferred until the first Get() call.
func NewLazyFrames(src []byte) *LazyFrames {
	return &LazyFrames{bytes: src}
}

// Get returns the decoded frames, per-frame delays, GIF loop count, and
// any decode error. Safe to call concurrently; the underlying decode
// runs exactly once.
func (l *LazyFrames) Get() ([]*image.RGBA, []int, int, error) {
	l.once.Do(func() {
		g, err := DecodeGIF(l.bytes)
		if err != nil {
			l.loadErr = err
			return
		}
		l.frames, l.delays = RenderFrames(g)
		l.loop = g.LoopCount
	})
	return l.frames, l.delays, l.loop, l.loadErr
}

// AnimateOverGIF runs fn over every frame of tmpl in parallel (one
// goroutine per frame) and returns the assembled GIF bytes. The caller's
// fn receives a fully-rendered RGBA frame and must return the paletted
// frame to emit. Per-frame delays and the loop count are copied from the
// template; output disposal is set to DisposalBackground on every frame.
//
// Parallelizing here matters because the dominant per-frame costs
// (resize + Floyd-Steinberg dither) are CPU-bound and embarrassingly
// parallel — sequential processing leaves cores idle for no reason.
func AnimateOverGIF(tmpl *LazyFrames, fn func(*image.RGBA) *image.Paletted) ([]byte, error) {
	frames, delays, loopCount, err := tmpl.Get()
	if err != nil {
		return nil, err
	}

	paletted := make([]*image.Paletted, len(frames))
	var wg sync.WaitGroup
	wg.Add(len(frames))
	for i, frame := range frames {
		i, frame := i, frame
		go func() {
			defer wg.Done()
			paletted[i] = fn(frame)
		}()
	}
	wg.Wait()

	disposal := make([]byte, len(paletted))
	for i := range disposal {
		disposal[i] = gif.DisposalBackground
	}

	return EncodeGIF(&gif.GIF{
		LoopCount: loopCount,
		Image:     paletted,
		Delay:     delays,
		Disposal:  disposal,
	})
}
