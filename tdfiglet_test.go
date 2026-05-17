package tdfiglet_test

import (
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
