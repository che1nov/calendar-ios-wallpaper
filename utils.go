package main

import (
	"fmt"
	"os"
	"time"

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

func fixedText(format string, a ...any) string {
	return fmt.Sprintf(format, a...)
}

func DaysInYear(year int) int {
	if time.Date(year, time.December, 31, 0, 0, 0, 0, time.UTC).YearDay() == 366 {
		return 366
	}
	return 365
}
