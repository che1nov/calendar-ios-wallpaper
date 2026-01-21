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

func RenderCalendar(now time.Time, theme Theme, mode, lang string) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, Width, Height))
	draw.Draw(img, img.Bounds(), &image.Uniform{theme.Background}, image.Point{}, draw.Src)

	_, left, percent := Progress(now)

	months := BuildMonths(now, lang)
	drawMonths(img, months, theme)

	drawText(img, fmt.Sprintf("%d d left   %d%%", left, percent), Width/2, Height-180, theme.Text, footerFace)
	return img
}

func drawMonths(img *image.RGBA, months []MonthData, theme Theme) {
	cellW := Width / 3
	cellH := (Height - 500) / 4
	offsetY := 320

	for i, m := range months {
		c := i % 3
		r := i / 3
		cx := c*cellW + cellW/2
		cy := offsetY + r*cellH + cellH/2
		drawMonth(img, cx, cy, m, theme)
	}
}

func drawMonth(img *image.RGBA, cx, cy int, m MonthData, theme Theme) {
	drawText(img, m.Name, cx, cy-80, theme.Text, monthFace)

	spacing := 32
	radius := 9
	startX := cx - 3*spacing
	startY := cy - 10

	for d := 0; d < m.Days; d++ {
		i := m.StartWeekday + d
		x := startX + (i%7)*spacing
		y := startY + (i/7)*spacing

		col := theme.Future
		if d < m.PassedDays {
			col = theme.Active
		}
		if m.IsCurrent && d == m.PassedDays-1 {
			col = theme.Today
		}
		drawCircle(img, x, y, radius, col)
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
