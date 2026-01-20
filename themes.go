package main

import "image/color"

type Theme struct {
	Background color.RGBA
	Active     color.RGBA
	Future     color.RGBA
	Text       color.RGBA
	Today      color.RGBA
}

var Themes = map[string]Theme{
	"ios": {
		Background: color.RGBA{0, 0, 0, 255},       // ЧЁРНЫЙ фон
		Active:     color.RGBA{220, 220, 220, 255}, // прошедшие дни
		Future:     color.RGBA{90, 90, 90, 255},    // будущие дни
		Text:       color.RGBA{200, 200, 200, 255}, // текст
		Today:      color.RGBA{255, 140, 0, 255},   // сегодня
	},
}
