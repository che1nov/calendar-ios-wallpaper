package main

import "image/color"

type Theme struct {
	Background color.RGBA
	Active     color.RGBA
	Future     color.RGBA
	Text       color.RGBA
	Today      color.RGBA

	WeekendGreen color.RGBA
	WeekendBlue  color.RGBA
	WeekendRed   color.RGBA
}

func IOSTheme() Theme {
	return Theme{
		Background: color.RGBA{0, 0, 0, 255},
		Active:     color.RGBA{220, 220, 220, 255},
		Future:     color.RGBA{90, 90, 90, 255},
		Text:       color.RGBA{200, 200, 200, 255},
		Today:      color.RGBA{255, 140, 0, 255},

		WeekendGreen: color.RGBA{110, 140, 110, 255},
		WeekendBlue:  color.RGBA{110, 125, 150, 255},
		WeekendRed:   color.RGBA{150, 110, 110, 255},
	}
}
