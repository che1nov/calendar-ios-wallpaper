package rendering

import (
	"os"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

func mustRead(path string) []byte {
	b, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return b
}

func mustParseFont(b []byte) *opentype.Font {
	f, err := opentype.Parse(b)
	if err != nil {
		panic(err)
	}
	return f
}

func mustFace(f *opentype.Font, size float64) font.Face {
	face, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size: size,
		DPI:  72,
	})
	if err != nil {
		panic(err)
	}
	return face
}
