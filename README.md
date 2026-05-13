~40 image filters and counting in pure Go — color transforms, edge detection, perspective warps, and animated GIFs — no third-party API or native dependencies.

bb_images is a pure-Go image manipulation library with about 40 filters, built to replace a third-party image-processing API entirely. It covers color transforms (sepia, posterize, color tints, deep-fry), edge-based effects (sketch, charcoal), spatial transforms (swirl, mosaic, low-poly triangulation, fitting text onto tilted-template signs), animated GIF effects (triggered, expand, glitch, rainbow, string-art), and overlay-based composites that drop an embedded image on top of the source. 

Every filter is a single top-level function in one flat package, so consumers get the whole library with one import. Overlay images and fonts are embedded into the binary via //go:embed; the tunable knobs — frame counts, opacity, jitter amplitude, etc. — live at the top of each filter's file.
