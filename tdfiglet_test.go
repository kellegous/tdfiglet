package tdfiglet_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kellegous/tdfiglet"
)

const (
	testFont    = "fonts/brndamgx.tdf"
	cFigletBin  = "tdfiglet/tdfiglet"
	compareText = "hello world"
)

// normalizeCOutput trims the leading blank line the C binary prints before
// rendered text, and any trailing newlines after the last row.
func normalizeCOutput(s string) string {
	s = strings.TrimPrefix(s, "\n")
	return strings.TrimRight(s, "\n")
}

func TestLoadFont(t *testing.T) {
	f, err := tdfiglet.LoadFontFile(testFont)
	if err != nil {
		t.Fatal(err)
	}
	if f.Name == "" {
		t.Fatal("expected font name")
	}
	if f.Height == 0 {
		t.Fatal("expected non-zero height")
	}
}

func TestCompareCImplementation(t *testing.T) {
	if _, err := os.Stat(cFigletBin); err != nil {
		t.Skipf("%s not present, skipping C comparison", cFigletBin)
	}

	entries, err := os.ReadDir("fonts")
	if err != nil {
		t.Fatal(err)
	}

	opts := tdfiglet.RenderOptions{
		Color:    tdfiglet.ColorANSI,
		Encoding: tdfiglet.EncodingUnicode,
	}

	for _, e := range entries {
		if filepath.Ext(e.Name()) != ".tdf" {
			continue
		}
		font := filepath.Join("fonts", e.Name())
		t.Run(e.Name(), func(t *testing.T) {
			cOut, err := exec.Command(cFigletBin, "-f", font, compareText).Output()
			if err != nil {
				t.Fatalf("c tdfiglet: %v", err)
			}

			f, err := tdfiglet.LoadFontFile(font)
			if err != nil {
				t.Fatal(err)
			}

			got := normalizeCOutput(f.RenderWith(compareText, opts))
			want := normalizeCOutput(string(cOut))
			if got != want {
				t.Fatal("render output differs from C implementation")
			}
		})
	}
}

func TestRenderMIRC(t *testing.T) {
	f, err := tdfiglet.LoadFontFile(testFont)
	if err != nil {
		t.Fatal(err)
	}
	out := f.RenderWith("hi", tdfiglet.RenderOptions{Color: tdfiglet.ColorMIRC})
	if !strings.Contains(out, "\x03") {
		t.Fatal("expected mIRC color codes")
	}
}
