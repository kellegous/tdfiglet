package tdfiglet_test

import (
	"fmt"
	"log"

	"github.com/kellegous/tdfiglet"
)

func ExampleLoadFontFile() {
	font, err := tdfiglet.LoadFontFile("fonts/brndamgx.tdf")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(font.Render("hello"))
}

func ExampleFont_RenderWith() {
	font, err := tdfiglet.LoadFontFile("fonts/brndamgx.tdf")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(font.RenderWith("hello", tdfiglet.RenderOptions{
		Justify: tdfiglet.JustifyCenter,
		Width:   80,
		Color:   tdfiglet.ColorANSI,
	}))
}
