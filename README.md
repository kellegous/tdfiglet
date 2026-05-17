# tdfiglet

A Go library for rendering large ASCII art text using [TheDraw](https://en.wikipedia.org/wiki/TheDraw) color font (`.tdf`) files.

## Features

- Load TheDraw FONTS (`.tdf`) color bitmap fonts
- Per-glyph foreground and background colors
- ANSI terminal escape sequences or mIRC color codes
- CP437 to UTF-8 conversion for Unicode terminals, or raw IBM bytes
- Left, center, and right justification within a configurable width

## Installation

```bash
go get github.com/kellegous/tdfiglet
```

## Usage

Load a font by path and render text with default options (left-aligned, 80 columns, ANSI colors, Unicode encoding):

```go
package main

import (
	"fmt"
	"os"

	"github.com/kellegous/tdfiglet"
)

func main() {
	font, err := tdfiglet.LoadFont("fonts/brndamgx.tdf")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Print(font.Render("hello"))
}
```

### Render options

Use `RenderWith` to control layout, colors, and encoding:

```go
out := font.RenderWith("hello", tdfiglet.RenderOptions{
	Justify:  tdfiglet.JustifyCenter,
	Width:    80,
	Color:    tdfiglet.ColorANSI,   // or tdfiglet.ColorMIRC
	Encoding: tdfiglet.EncodingUnicode, // or tdfiglet.EncodingASCII
})
```

Zero values in `RenderOptions` use library defaults: left justify, width 80, ANSI color, Unicode encoding.

## Fonts

This library reads TheDraw color font files (type 2). Fonts are not embedded in the library, but you can copy those you want to use from the [fonts](fonts) directory.

Supported characters are the 94 printable ASCII glyphs (`!` through `~`).

## Acknowledgments

This library is an approximate port of the C program [tdfiglet](https://github.com/tat3r/tdfiglet). In fact, this library is currently tested against the C implementation.

## License

MIT — see [LICENSE](LICENSE).
