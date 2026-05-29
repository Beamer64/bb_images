# bb_images

**~50 image filters in pure Go** — color transforms, edge detection, perspective warps, and animated GIFs — no third-party API, no native dependencies, no CGO.

```
go get github.com/Beamer64/bb_images
```

---

## What it is

bb_images is a pure-Go image manipulation library built to replace a third-party image-processing API entirely. It covers:

- **Color transforms** — sepia, posterize, inversion, tints (earth, freeze, ground, night), deep-fry, blur
- **Edge effects** — Sobel, sketch, charcoal, HOG
- **Spatial transforms** — swirl, mosaic, low-poly triangulation, "magik" liquid-rescale, pixelation, mirror/flip, string-art rendering
- **Animated GIFs** — triggered, expand, glitch, rainbow, rain, shake, spin, TV static
- **Overlays** — drop an embedded PNG over the source (jail, wasted, etc.)
- **Sign / meme composites** — fit text into the tilted text region of templates like change-my-mind, retro meme, fake tweet, fake YouTube card, batman slap, wanted poster
- **Misc** — ASCII art, paint, RGB split, burn

Every effect is one top-level function in a topic subpackage (`color.Sepia`, `animated.Triggered`, `signs.ChangeMyMind`, …), all sharing the same signature: take an `image.Image`, return encoded bytes (PNG for static effects, GIF for animated) and an error. Template PNGs and fonts are embedded into the binary at build time via `//go:embed` — there's nothing to ship alongside the library, no runtime asset loading.

The tunable knobs — frame counts, opacity, jitter amplitude, palette choice — live as `const`s at the top of each filter's file so they're easy to find and tweak.

---

## Showcase

A few showcase images can be found in the [showcase](https://github.com/Beamer64/bb_images/tree/264ea9fcb981b35e3baca5a986308c6195af92dc/showcase) folder.

---

## Quick start

```go
package main

import (
    "image"
    _ "image/jpeg"
    "os"

    "github.com/Beamer64/bb_images/color"
)

func main() {
    f, err := os.Open("input.jpg")
    if err != nil {
        panic(err)
    }
    defer f.Close()

    src, _, err := image.Decode(f)
    if err != nil {
        panic(err)
    }

    out, err := color.Sepia(src)
    if err != nil {
        panic(err)
    }

    if err := os.WriteFile("output.png", out, 0o644); err != nil {
        panic(err)
    }
}
```

That's the whole shape: decode → call effect → write bytes. Same pattern for every subpackage; animated effects just return a GIF instead of a PNG.

```go
import "github.com/Beamer64/bb_images/animated"

gifBytes, err := animated.Triggered(src) // returns GIF bytes
```

## The catalogue

Each subpackage owns a topic and exports its effects as plain top-level functions. Import only the ones you need.

| Package | What it does | Output | Examples |
|---|---|---|---|
| [`color`](color/) | Per-pixel color transforms | PNG | `Sepia`, `Invert`, `Posterize`, `Deepfry`, `Blur`, `Earth`, `Freeze`, `Ground`, `Night`, `Tint` |
| [`edges`](edges/) | Edge detection / linework | PNG | `Sobel`, `Sketch`, `Charcoal`, `Hog` |
| [`spatial`](spatial/) | Geometric / structural transforms | PNG | `Swirl`, `Mirror`, `Flip`, `Pixelate`, `Mosaic`, `Triangle`, `Magik`, `Stringify`, `Colors` |
| [`overlays`](overlays/) | Embedded PNG layered over the source | PNG | `Jail`, `Wasted` |
| [`signs`](signs/) | Avatar + fitted text on tilted templates | PNG | `ChangeMyMind`, `Tweet`, `Discord`, `YouTube`, `RetroMeme`, `BatmanSlap`, `ThanksObama`, `Wanted`, `WhyAreYouGay`, `FiveGuysOneGirl` |
| [`animated`](animated/) | Procedural GIF outputs | GIF | `Triggered`, `Expand`, `Glitch`, `GlitchStatic`, `Rain`, `Shake`, `Spin`, `TVStatic`, `Rainbow` |
| [`special`](special/) | Algorithmic one-offs | PNG | `ASCII`, `Paint`, `RGB`, `Burn` |

Helpers shared across packages live under [`internal/draw`](internal/draw/) (`Decode`, `EncodePNG`, `RenderFrames`, `LazyFrames`, `AnimateOverGIF`) and [`internal/templates`](internal/templates/) (`Detect`, `ConnectedRegions` — marker-color region detection for memes).

## Design notes

A few things to know if you're poking under the hood or contributing:

### Uniform API

Every effect takes `image.Image` and returns `([]byte, error)`. There are no flags, options structs, or config arguments — the public surface is intentionally minimal so consumers can swap effects by name without code changes. Tunable parameters live as file-local `const`s; the way to "customize" a filter is to fork it.

### Sign / meme templates use marker images

Templates like `change-my-mind.png` ship with a sibling `change-my-mind.marker.png` painted in pure RGB primaries (magenta = text region, green/blue = avatar slots, cyan/yellow = secondary text). The [`templates`](internal/templates/) package detects those regions and returns rectangles for the renderer to fill. This lets a single template support tilted text, perspective warps, and multi-avatar layouts without per-template detection code.

### GIF pipeline

Animated effects go through `internal/draw`'s shared pipeline:

- **`LazyFrames`** — decodes a GIF template once per process and caches the frames
- **`RenderFrames`** — runs the per-frame closure in parallel
- **`AnimateOverGIF`** — composites the source over each frame, quantizes to the Plan 9 palette with Floyd-Steinberg dithering, and assembles the output GIF

Resizes use `imaging.Linear` rather than `imaging.Lanczos`: palette quantization eats most of the quality difference and Linear is roughly 2× faster.

### No external services

No HTTP calls, no native libraries, no CGO. The whole library is `go build`able on any platform Go supports. Template PNGs and the Go basic font are embedded via `//go:embed`, so the only runtime input is the source image you hand in.

## Dependencies

| Module | Purpose |
|---|---|
| [`github.com/disintegration/imaging`](https://github.com/disintegration/imaging) | Resize, convolution, NRGBA helpers — vendored |
| [`golang.org/x/image`](https://pkg.go.dev/golang.org/x/image) | Font rendering (`opentype`, `basicfont`), extra image formats — vendored |
| Go standard library (`image`, `image/png`, `image/gif`, `image/color`, `image/draw`, …) | Everything else |

That's it. No CGO, no system libs, no Docker, no API keys.

## Used by

- **[BuddieBot](https://github.com/Beamer64/BuddieBot)** — Discord bot; all `/image *` commands are powered by this library.

## Contributing

Adding a new effect is mechanical:

1. Pick the right subpackage (or open an issue if you think it needs a new one).
2. Drop a new file: `package <name>` declaration, `//go:embed` any assets next to it, tunables as file-local `const`s, one exported function `Foo(src image.Image) ([]byte, error)`.
3. Encode through `internal/draw.EncodePNG` (static) or `internal/draw.AnimateOverGIF` / build a `gif.GIF` directly (animated).
4. Add a test that decodes a sample input and asserts the output is non-empty + decodes cleanly. The existing `*_test.go` files in each package are templates worth copying.

PRs welcome — keep the API shape (`image.Image` in, `([]byte, error)` out) so consumers can keep swapping effects by name.

## License

MIT — see [LICENSE](LICENSE).
