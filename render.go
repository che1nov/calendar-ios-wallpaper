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

var fontBase *opentype.Font

func init() {
	b, err := os.ReadFile("fonts/Inter-Regular.ttf")
	if err != nil {
		panic(err)
	}
	fontBase, _ = opentype.Parse(b)
}

func makeFont(size float64) font.Face {
	f, _ := opentype.NewFace(fontBase, &opentype.FaceOptions{
		Size: size,
		DPI:  72,
	})
	return f
}

func RenderCalendar(
	now time.Time,
	device DeviceProfile,
	theme Theme,
	mode, lang, weekends string,
) *image.RGBA {

	img := image.NewRGBA(image.Rect(0, 0, device.Width, device.Height))
	draw.Draw(img, img.Bounds(), &image.Uniform{theme.Background}, image.Point{}, draw.Src)

	months := BuildMonths(now, lang)

	drawMonths(img, months, device, theme, weekends)

	_, left, percent := Progress(now)
	drawFooter(img, left, 100-percent, device, theme, lang)

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

	// контент строго между safe-зонами
	contentTop := device.TopSafe
	contentBottom := device.Height - device.BottomSafe - device.FooterOffset - int(device.FooterFont) - 20

	usableHeight := contentBottom - contentTop

	cellW := device.Width / cols
	cellH := usableHeight / rows

	for i, m := range months {
		c := i % cols
		r := i / cols

		cx := c*cellW + cellW/2
		cy := contentTop + r*cellH + cellH/2

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
	monthFace := makeFont(device.MonthFont)

	drawText(img, m.Name, cx, cy-70, theme.Text, monthFace)

	startX := cx - (6 * device.DotSpacing / 2)
	startY := cy

	for day := 0; day < m.Days; day++ {
		idx := m.StartWeekday + day
		col := idx % 7
		row := idx / 7

		x := startX + col*device.DotSpacing
		y := startY + row*device.DotSpacing

		dot := theme.Future
		if day < m.PassedDays {
			dot = theme.Active
		}

		if col >= 5 {
			switch weekends {
			case "gray":
				dot = theme.WeekendGray
			case "green":
				dot = theme.WeekendGreen
			case "blue":
				dot = theme.WeekendBlue
			case "red":
				dot = theme.WeekendRed
			}
		}

		if m.IsCurrent && day == m.PassedDays-1 {
			dot = theme.Today
		}

		drawCircle(img, x, y, device.DotRadius, dot)
	}
}

func drawFooter(
	img *image.RGBA,
	left, percent int,
	device DeviceProfile,
	theme Theme,
	lang string,
) {
	footerFace := makeFont(device.FooterFont)

	text := footerText(left, percent, lang)

	y := device.Height -
		device.BottomSafe -
		device.FooterOffset

	drawText(img, text, device.Width/2, y, theme.Text, footerFace)
}

func footerText(left, percent int, lang string) string {
	if lang == "ru" {
		return fmt.Sprintf("%dдн. осталось   %d%%", left, percent)
	}
	return fmt.Sprintf("%dd left   %d%%", left, percent)
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
