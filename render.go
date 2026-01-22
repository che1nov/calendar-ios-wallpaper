package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"os"
	"time"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

const (
	Width  = 1179
	Height = 2556
)

var monthFace, footerFace font.Face

func init() {
	data, _ := os.ReadFile("fonts/Inter-Regular.ttf")
	f, _ := opentype.Parse(data)

	monthFace, _ = opentype.NewFace(f, &opentype.FaceOptions{Size: 40, DPI: 72})
	footerFace, _ = opentype.NewFace(f, &opentype.FaceOptions{Size: 34, DPI: 72})
}

func RenderCalendar(
	now time.Time,
	theme Theme,
	mode, lang, weekends string,
) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, Width, Height))
	draw.Draw(img, img.Bounds(), &image.Uniform{theme.Background}, image.Point{}, draw.Src)

	_, left, percent := Progress(now)

	months := BuildMonths(now, lang)
	drawMonths(img, months, theme, weekends)

	drawText(img, fmt.Sprintf("%d d left   %d%%", left, percent), Width/2, Height-180, theme.Text, footerFace)
	return img
}

func drawMonths(
	img *image.RGBA,
	months []MonthData,
	theme Theme,
	weekends string,
) {
	const cols, rows = 3, 4

	topSafe, bottomSafe := calcSafeAreas(Height)
	usableHeight := Height - topSafe - bottomSafe

	cellW := Width / cols
	cellH := usableHeight / rows

	for i, m := range months {
		col := i % cols
		row := i / cols

		cx := col*cellW + cellW/2
		cy := topSafe + row*cellH + cellH/2

		drawMonth(img, cx, cy, m, theme, weekends)
	}
}

func drawMonth(
	img *image.RGBA,
	cx, cy int,
	m MonthData,
	theme Theme,
	weekends string,
) {
	drawText(img, m.Name, cx, cy-80, theme.Text, monthFace)

	spacing := 32
	radius := 9
	startX := cx - 3*spacing
	startY := cy - 10

	for d := 0; d < m.Days; d++ {
		i := m.StartWeekday + d
		colIndex := i % 7
		row := i / 7

		x := startX + colIndex*spacing
		y := startY + row*spacing

		dot := theme.Future

		isWeekend := colIndex == 5 || colIndex == 6

		if isWeekend {
			switch weekends {
			case "gray":
				dot = theme.WeekendGray
			case "green":
				dot = theme.WeekendGreen
			case "blue":
				dot = theme.WeekendBlue
			case "red":
				dot = theme.WeekendRed
			case "off":
				// ничего
			}
		}

		// прошедшие дни
		if d < m.PassedDays {
			dot = theme.Active
		}

		// сегодня — приоритет
		if m.IsCurrent && d == m.PassedDays-1 {
			dot = theme.Today
		}

		drawCircle(img, x, y, radius, dot)
	}

}

func drawText(img *image.RGBA, text string, cx, y int, col color.Color, face font.Face) {
	d := &font.Drawer{Dst: img, Src: image.NewUniform(col), Face: face}
	w := d.MeasureString(text).Round()
	d.Dot = fixed.P(cx-w/2, y)
	d.DrawString(text)
}

func drawCircle(img *image.RGBA, cx, cy, r int, col color.Color) {
	for y := -r; y <= r; y++ {
		for x := -r; x <= r; x++ {
			if x*x+y*y <= r*r {
				img.Set(cx+x, cy+y, col)
			}
		}
	}
}

func calcSafeAreas(height int) (top, bottom int) {
	top = int(float64(height) * 0.22)    // зона часов / Dynamic Island
	bottom = int(float64(height) * 0.16) // зона кнопок
	return
}
