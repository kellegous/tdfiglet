package tdfiglet

import (
	"encoding/binary"
	"errors"
	"fmt"
	"os"
)

const (
	NumChars    = 94
	colorFont   = 2
	dataOffset  = 233
	charListOff = 45
)

var (
	fontMagic    = []byte("\x13TheDraw FONTS file\x1a")
	asciiCharSet = []byte("!\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~")
)

// Font is a loaded TheDraw color (.tdf) font.
type Font struct {
	Name    string
	Spacing uint8
	Height  uint8
	glyphs  [NumChars]*Glyph
}

// Glyph is one character's bitmap from a font file.
type Glyph struct {
	Width  uint8
	Height uint8
	Cells  []Cell
}

// Cell is one character position in a glyph row.
type Cell struct {
	Ch    byte
	Color uint8
}

// LoadFont reads a TheDraw FONTS (.tdf) file.
func LoadFont(path string) (*Font, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return parseFont(raw)
}

func parseFont(raw []byte) (*Font, error) {
	if len(raw) < dataOffset {
		return nil, errors.New("tdfiglet: font file too small")
	}
	if !bytesHasPrefix(raw, fontMagic) {
		return nil, errors.New("tdfiglet: invalid font magic")
	}

	if raw[41] != colorFont {
		return nil, fmt.Errorf("tdfiglet: unsupported font type %d", raw[41])
	}

	nameLen := int(raw[24])
	if 25+nameLen > len(raw) {
		return nil, errors.New("tdfiglet: invalid font name length")
	}

	charOffsets := make([]uint16, NumChars)
	for i := 0; i < NumChars; i++ {
		charOffsets[i] = binary.LittleEndian.Uint16(raw[charListOff+i*2:])
	}

	data := raw[dataOffset:]
	f := &Font{
		Name:    string(raw[25 : 25+nameLen]),
		Spacing: raw[42],
	}

	for i := 0; i < NumChars; i++ {
		if charOffsets[i] == 0xffff {
			continue
		}
		off := int(charOffsets[i])
		if off+2 > len(data) {
			return nil, fmt.Errorf("tdfiglet: glyph %q out of range", asciiCharSet[i])
		}
		if int(data[off+1]) > int(f.Height) {
			f.Height = data[off+1]
		}
	}

	for i := 0; i < NumChars; i++ {
		if charOffsets[i] == 0xffff {
			continue
		}
		g, err := readGlyph(data, charOffsets[i], f.Height)
		if err != nil {
			return nil, err
		}
		f.glyphs[i] = g
	}

	return f, nil
}

func readGlyph(data []byte, offset uint16, fontHeight uint8) (*Glyph, error) {
	base := int(offset)
	if base+2 > len(data) {
		return nil, errors.New("tdfiglet: truncated glyph header")
	}

	p := data[base:]
	width := int(p[0])
	height := int(p[1])
	p = p[2:]

	if height > int(fontHeight) {
		fontHeight = uint8(height)
	}

	cells := make([]Cell, width*int(fontHeight))
	row, col := 0, 0

	for len(p) > 0 && p[0] != 0 {
		ch := p[0]
		p = p[1:]

		if ch == '\r' {
			row++
			col = 0
			continue
		}
		if len(p) == 0 {
			break
		}
		color := p[0]
		p = p[1:]

		if ch < 0x20 {
			ch = ' '
		}
		if row < int(fontHeight) && col < width {
			cells[row*width+col] = Cell{Ch: ch, Color: color}
		}
		col++
	}

	return &Glyph{
		Width:  uint8(width),
		Height: uint8(height),
		Cells:  cells,
	}, nil
}

func (f *Font) glyphAt(c byte) *Glyph {
	for i := 0; i < NumChars; i++ {
		if asciiCharSet[i] == c {
			return f.glyphs[i]
		}
	}
	return nil
}

func bytesHasPrefix(b, prefix []byte) bool {
	return len(b) >= len(prefix) && string(b[:len(prefix)]) == string(prefix)
}
