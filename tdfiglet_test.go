package tdfiglet_test

import (
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/kellegous/tdfiglet"
)

const testFont = "fonts/brndamgx.tdf"

func TestLoadFont(t *testing.T) {
	f, err := tdfiglet.LoadFont(testFont)
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

func TestRender3DAsciiDevmode(t *testing.T) {
	want, err := os.ReadFile("3d-ascii.txt")
	if err != nil {
		t.Fatal(err)
	}

	f, err := tdfiglet.LoadFont("fonts/3d-ascii.tdf")
	if err != nil {
		t.Fatal(err)
	}

	out := f.RenderWith("devmode", tdfiglet.RenderOptions{
		Color:    tdfiglet.ColorANSI,
		Encoding: tdfiglet.EncodingUnicode,
	})
	got := ansiStrip(out)
	if got != string(want) {
		t.Fatalf("render mismatch:\n%s", diffLines(string(want), got))
	}
}

var ansiRE = regexp.MustCompile(`\x1b\[[0-9;]*m|\x03\d*,\d*|\x03`)

func ansiStrip(s string) string {
	return ansiRE.ReplaceAllString(s, "")
}

func diffLines(want, got string) string {
	wLines := strings.Split(strings.TrimRight(want, "\n"), "\n")
	gLines := strings.Split(strings.TrimRight(got, "\n"), "\n")
	var b strings.Builder
	n := len(wLines)
	if len(gLines) > n {
		n = len(gLines)
	}
	for i := 0; i < n; i++ {
		w := ""
		if i < len(wLines) {
			w = wLines[i]
		}
		g := ""
		if i < len(gLines) {
			g = gLines[i]
		}
		if w != g {
			b.WriteString("want: " + w + "\n")
			b.WriteString(" got: " + g + "\n")
		}
	}
	return b.String()
}

func TestRenderMIRC(t *testing.T) {
	f, err := tdfiglet.LoadFont(testFont)
	if err != nil {
		t.Fatal(err)
	}
	out := f.RenderWith("hi", tdfiglet.RenderOptions{Color: tdfiglet.ColorMIRC})
	if !strings.Contains(out, "\x03") {
		t.Fatal("expected mIRC color codes")
	}
}
