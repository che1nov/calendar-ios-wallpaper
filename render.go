package main

import (
	"image"
	"image/color"
	"image/draw"
	"time"

	"golang.org/x/image/font"
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

	if mode == "months" {
		months := BuildMonths(now, lang)
		drawMonths(img, months, device, theme, weekends)
		drawFooter(img, device, theme)
	}

	return img
}

func drawMonths(
	img *image.RGBA,
	months []MonthData,
	device DeviceProfile,
	theme Theme,
	weekends string,
) {
	const cols, rows = 3, 4

	usableTop := device.ClockInset
	usableHeight := device.Height - device.ClockInset - device.BottomInset

	cellW := device.Width / cols
	cellH := usableHeight / rows

	for i, m := range months {
		c := i % cols
		r := i / cols

		cx := c*cellW + cellW/2
		cy := usableTop + r*cellH + cellH/2

		drawMonth(img, cx, cy, m, device, theme, weekends)
	}
}

func drawMonth(
	img *image.RGBA,
	cx, cy int,
	m MonthData,
	device DeviceProfile,
	theme Theme,
	weekends string,
) {
	titleColor := theme.Text
	if m.IsCurrent {
		titleColor = theme.Today
	}

	drawText(img, m.Name, cx, cy-70, titleColor, monthFace)

	spacing := 32
	radius := 9
	startX := cx - (6 * spacing / 2)
	startY := cy

	for d := 0; d < m.Days; d++ {
		i := m.StartWeekday + d
		col := i % 7
		row := i / 7

		x := startX + col*spacing
		y := startY + row*spacing

		colr := theme.Future
		if d < m.PassedDays {
			colr = theme.Active
		}

		if m.IsCurrent && d == m.PassedDays-1 {
			colr = theme.Today
		}

		drawCircle(img, x, y, radius, colr)
	}
}

func drawFooter(img *image.RGBA, device DeviceProfile, theme Theme) {
	drawText(img, "Year progress", device.Width/2, device.Height-device.BottomInset/2, theme.Text, footerFace)
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
