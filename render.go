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

	// iPhone 15 lock screen safe areas
	TopSafe    = 380
	BottomSafe = 200
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

	f, err := opentype.Parse(fontBytes)
	if err != nil {
		panic(err)
	}

	monthFace, _ = opentype.NewFace(f, &opentype.FaceOptions{
		Size: 40,
		DPI:  72,
	})

	footerFace, _ = opentype.NewFace(f, &opentype.FaceOptions{
		Size: 34,
		DPI:  72,
	})
}

// Главная функция рендера
func RenderCalendarWithLang(now time.Time, theme Theme, mode string, lang string) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, Width, Height))
	draw.Draw(img, img.Bounds(), &image.Uniform{theme.Background}, image.Point{}, draw.Src)

	_, left, percent := Progress(now)

	if mode == "months" {
		months := BuildMonths(now, lang)
		drawMonths(img, months, theme)

		passed := 100 - percent
		drawFooter(img, left, passed, theme)
	}

	return img
}

func drawMonths(img *image.RGBA, months []MonthData, theme Theme) {
	cols := 3

	startY := TopSafe + 320 // где начинается первый ряд
	rowGap := 220           // расстояние между рядами месяцев

	cellW := Width / cols

	for i, m := range months {
		c := i % cols
		r := i / cols

		cx := c*cellW + cellW/2
		cy := startY + r*rowGap

		drawMonth(img, cx, cy, m, theme)
	}
}

func drawMonth(img *image.RGBA, cx, cy int, m MonthData, theme Theme) {
	titleColor := theme.Text
	if m.IsCurrent {
		titleColor = theme.Today
	}

	// Название месяца — ближе к сетке
	drawText(img, m.Name, cx, cy-80, titleColor, monthFace)

	cols := 7
	spacing := 32 // ⬅️ было 36
	radius := 9   // ⬅️ было 10

	startX := cx - (cols-1)*spacing/2
	startY := cy - 10 // ⬅️ было -10, оставляем компактно

	for i := 0; i < m.Days; i++ {
		x := startX + (i%cols)*spacing
		y := startY + (i/cols)*spacing

		col := theme.Future
		if i < m.PassedDays {
			col = theme.Active
		}
		if m.IsCurrent && m.PassedDays > 0 && i == m.PassedDays-1 {
			col = theme.Today
		}

		drawCircle(img, x, y, radius, col)
	}
}

func drawFooter(img *image.RGBA, left, passed int, theme Theme) {
	text := fmt.Sprintf("%d d left   %d%%", left, passed)

	y := Height - BottomSafe + 60
	drawText(img, text, Width/2, y, theme.Text, footerFace)
}

func drawText(img *image.RGBA, text string, cx int, y int, col color.Color, face font.Face) {
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
