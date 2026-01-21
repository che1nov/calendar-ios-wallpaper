package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"time"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

func RenderCalendar(
	now time.Time,
	device DeviceProfile,
	theme Theme,
	mode, lang, weekends string,
) *image.RGBA {

	img := image.NewRGBA(image.Rect(0, 0, device.Width, device.Height))
	draw.Draw(img, img.Bounds(), &image.Uniform{theme.Background}, image.Point{}, draw.Src)

	if mode != "months" {
		return img
	}

	months := BuildMonths(now, lang)

	top := device.TopInset + device.ClockInset
	bottom := device.BottomInset

	gridHeight := device.Height - top - bottom

	cellW := device.Width / 3
	cellH := gridHeight / 4

	for i, m := range months {
		c := i % 3
		r := i / 3

		cx := c*cellW + cellW/2
		cy := top + r*cellH + cellH/2

		drawMonth(img, cx, cy, m, theme, weekends)
	}

	drawFooter(img, device, theme, lang, now)

	return img
}

func drawMonth(img *image.RGBA, cx, cy int, m MonthData, theme Theme, weekends string) {
	drawText(img, m.Name, cx, cy-70, theme.Text)

	spacing := 30
	radius := 8
	startX := cx - 3*spacing
	startY := cy - 10

	for d := 0; d < m.Days; d++ {
		i := m.StartWeekday + d
		col := i % 7
		row := i / 7

		x := startX + col*spacing
		y := startY + row*spacing

		colorDot := theme.Future

		if d < m.PassedDays {
			colorDot = theme.Active
		}

		if col >= 5 {
			switch weekends {
			case "gray":
				colorDot = theme.WeekendGray
			case "green":
				colorDot = theme.WeekendGreen
			case "blue":
				colorDot = theme.WeekendBlue
			case "red":
				colorDot = theme.WeekendRed
			}
		}

		if m.IsCurrent && d == m.PassedDays-1 {
			colorDot = theme.Today
		}

		drawCircle(img, x, y, radius, colorDot)
	}
}

func drawFooter(img *image.RGBA, device DeviceProfile, theme Theme, lang string, now time.Time) {
	left := 365 - now.YearDay()
	txt := fmt.Sprintf("%d d left", left)
	if lang == "ru" {
		txt = fmt.Sprintf("%d дн. осталось", left)
	}
	drawText(img, txt, device.Width/2, device.Height-device.BottomInset+40, theme.Text)
}

func drawText(img *image.RGBA, text string, cx, y int, col color.Color) {
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
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
