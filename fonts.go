package main

import (
	"os"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

var (
	monthFace  font.Face
	footerFace font.Face
)

func init() {
	data, err := os.ReadFile("fonts/Inter-Regular.ttf")
	if err != nil {
		panic(err)
	}

	f, err := opentype.Parse(data)
	if err != nil {
		panic(err)
	}

	monthFace, _ = opentype.NewFace(f, &opentype.FaceOptions{Size: 36, DPI: 72})
	footerFace, _ = opentype.NewFace(f, &opentype.FaceOptions{Size: 30, DPI: 72})
}
