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

var (
	monthFace  font.Face
	footerFace font.Face
)

func init() {
	fontBytes, err := os.ReadFile("fonts/Inter-Regular.ttf")
	if err != nil {
		panic("fonts/Inter-Regular.ttf not found")
	}

	f, _ := opentype.Parse(fontBytes)

	monthFace, _ = opentype.NewFace(f, &opentype.FaceOptions{
		Size: 40,
		DPI:  72,
	})

	footerFace, _ = opentype.NewFace(f, &opentype.FaceOptions{
		Size: 34,
		DPI:  72,
	})
}

func RenderCalendar(
	now time.Time,
	theme Theme,
	mode, lang string,
	weekends string,
) *image.RGBA {

	img := image.NewRGBA(image.Rect(0, 0, Width, Height))
	draw.Draw(img, img.Bounds(), &image.Uniform{theme.Background}, image.Point{}, draw.Src)

	_, left, percent := Progress(now)

	if mode == "months" {
		months := BuildMonths(now, lang)
		drawMonths(img, months, theme, weekends)

		drawFooter(img, left, 100-percent, theme, lang)
	}

	return img
}

func drawMonths(img *image.RGBA, months []MonthData, theme Theme, weekendMode string) {
	const cols, rows = 3, 4

	topSafe, bottomSafe := calcSafeAreas(Height)
	usableHeight := Height - topSafe - bottomSafe

	cellW := Width / cols
	cellH := usableHeight / rows

	for i, m := range months {
		c := i % cols
		r := i / cols

		cx := c*cellW + cellW/2
		cy := topSafe + r*cellH + cellH/2

		drawMonth(img, cx, cy, m, theme, weekendMode)
	}
}

func drawMonth(
	img *image.RGBA,
	cx, cy int,
	m MonthData,
	theme Theme,
	weekends string,
) {
	titleColor := theme.Text
	if m.IsCurrent {
		titleColor = theme.Today
	}

	drawText(img, m.Name, cx, cy-80, titleColor, monthFace)

	cols := 7
	spacing := 32
	radius := 9

	startX := cx - (cols-1)*spacing/2
	startY := cy - 10

	for day := 0; day < m.Days; day++ {

		// позиция в сетке
		visualIndex := m.StartWeekday + day
		col := visualIndex % 7
		row := visualIndex / 7

		x := startX + col*spacing
		y := startY + row*spacing

		dotColor := theme.Future

		// прошедшие дни
		if day+1 < m.PassedDays {
			dotColor = theme.Active
		}

		// выходные (Сб = 5, Вс = 6)
		if col == 5 || col == 6 {
			switch weekends {
			case "gray":
				dotColor = theme.WeekendGray
			case "green":
				dotColor = theme.WeekendGreen
			case "blue":
				dotColor = theme.WeekendBlue
			case "red":
				dotColor = theme.WeekendRed
			}
		}

		// СЕГОДНЯ — абсолютный приоритет
		if m.IsCurrent && day+1 == m.PassedDays {
			dotColor = theme.Today
		}

		drawCircle(img, x, y, radius, dotColor)
	}
}

func drawFooter(img *image.RGBA, left, passed int, theme Theme, lang string) {
	text := footerText(left, passed, lang)
	drawText(img, text, Width/2, footerY(), theme.Text, footerFace)
}

func footerText(left, passed int, lang string) string {
	if lang == "ru" {
		return fmt.Sprintf("%d дн. осталось   %d%%", left, passed)
	}
	return fmt.Sprintf("%d d left   %d%%", left, passed)
}

func footerY() int {
	_, bottomSafe := calcSafeAreas(Height)
	return Height - bottomSafe + 60
}

func drawText(img *image.RGBA, text string, cx, y int, col color.Color, face font.Face) {
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: face,
	}

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
	top = int(float64(height) * 0.22)
	bottom = int(float64(height) * 0.16)
	return
}
