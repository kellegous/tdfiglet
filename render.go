package tdfiglet

import (
	"fmt"
	"strings"

	"golang.org/x/text/encoding/charmap"
)

const defaultWidth = 80

// Justify controls horizontal alignment of rendered text.
type Justify int

const (
	JustifyLeft Justify = iota
	JustifyRight
	JustifyCenter
)

// ColorFormat selects ANSI escape sequences or mIRC color codes.
type ColorFormat int

const (
	ColorANSI ColorFormat = iota
	ColorMIRC
)

// Encoding selects Unicode (CP437 → UTF-8) or raw IBM bytes for glyphs.
type Encoding int

const (
	EncodingUnicode Encoding = iota
	EncodingASCII
)

// RenderOptions configures text output. Zero values use library defaults
// (left justify, width 80, ANSI color, Unicode encoding).
type RenderOptions struct {
	Justify  Justify
	Width    int
	Color    ColorFormat
	Encoding Encoding
}

func (o *RenderOptions) applyDefaults() {
	if o.Width == 0 {
		o.Width = defaultWidth
	}
}

// Render draws text using default options (left, width 80, ANSI, Unicode).
func (f *Font) Render(text string) string {
	return f.RenderWith(text, RenderOptions{})
}

// RenderWith draws text using the given options.
func (f *Font) RenderWith(text string, opt RenderOptions) string {
	opt.applyDefaults()
	return renderString(text, f, opt)
}

var (
	fgANSI = [16]int{30, 34, 32, 36, 31, 35, 33, 37, 90, 94, 92, 96, 91, 95, 93, 97}
	bgANSI = [8]int{40, 44, 42, 46, 41, 45, 43, 47}
	fgMIRC = [16]int{1, 2, 3, 10, 5, 6, 7, 15, 14, 12, 9, 11, 4, 13, 8, 0}
	bgMIRC = [16]int{1, 2, 3, 10, 5, 6, 7, 15, 14, 12, 9, 11, 4, 13, 8, 0}

	cp437 = charmap.CodePage437.NewDecoder()
)

func ibmToUTF8(ch byte) string {
	r, _ := cp437.Bytes([]byte{ch})
	if len(r) == 0 {
		return " "
	}
	return string(r)
}

func cellChar(ch byte, enc Encoding) string {
	if ch < 0x20 {
		ch = ' '
	}
	if enc == EncodingUnicode {
		return ibmToUTF8(ch)
	}
	return string(ch)
}

func colorSeq(color uint8, format ColorFormat) string {
	fg := int(color & 0x0f)
	bg := int((color & 0xf0) >> 4)
	if format == ColorANSI {
		return fmt.Sprintf("\x1b[%d;%dm", fgANSI[fg], bgANSI[bg])
	}
	return fmt.Sprintf("\x03%d,%d", fgMIRC[fg], bgMIRC[bg])
}

func colorReset(format ColorFormat) string {
	if format == ColorANSI {
		return "\x1b[0m"
	}
	return "\x03"
}

func renderRow(g *Glyph, row int, opt RenderOptions, b *strings.Builder) {
	if g == nil || g.Width == 0 {
		return
	}
	rows := len(g.Cells) / int(g.Width)
	if row >= rows {
		return
	}
	var lastColor uint8
	for col := 0; col < int(g.Width); col++ {
		cell := g.Cells[int(g.Width)*row+col]
		if col == 0 || cell.Color != lastColor {
			b.WriteString(colorSeq(cell.Color, opt.Color))
			lastColor = cell.Color
		}
		b.WriteString(cellChar(cell.Ch, opt.Encoding))
	}
	b.WriteString(colorReset(opt.Color))
}

func renderString(text string, font *Font, opt RenderOptions) string {
	maxHeight := 0
	lineWidth := 0
	bytes := []byte(text)
	length := len(bytes)

	for i := 0; i < length; i++ {
		g := font.glyphAt(bytes[i])
		if g == nil {
			continue
		}
		if int(g.Height) > maxHeight {
			maxHeight = int(g.Height)
		}
		lineWidth += int(g.Width)
		if lineWidth+1 < length {
			lineWidth += int(font.Spacing)
		}
	}

	padding := 0
	switch opt.Justify {
	case JustifyCenter:
		padding = (opt.Width - lineWidth) / 2
	case JustifyRight:
		padding = opt.Width - lineWidth
	}
	if padding < 0 {
		padding = 0
	}

	var b strings.Builder
	for row := 0; row < maxHeight; row++ {
		b.WriteString(strings.Repeat(" ", padding))
		for i := 0; i < length; i++ {
			g := font.glyphAt(bytes[i])
			if g == nil {
				continue
			}
			renderRow(g, row, opt, &b)
			b.WriteString(colorReset(opt.Color))
			b.WriteString(strings.Repeat(" ", int(font.Spacing)))
		}
		b.WriteString(colorReset(opt.Color))
		b.WriteByte('\n')
	}
	return b.String()
}
