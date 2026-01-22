package main

import (
	"image"
	"image/color"
	"image/draw"
	"time"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

const (
	monthTitleOffsetY = 40 // было ~70, сдвигаем ниже
)

var (
	monthFace  font.Face
	footerFace font.Face
)

func init() {
	fontBytes := mustRead("fonts/Inter-Regular.ttf")
	f := mustParseFont(fontBytes)

	monthFace = mustFace(f, 38)
	footerFace = mustFace(f, 30)
}

// ===============================
// MAIN ENTRY
// ===============================

func RenderCalendar(
	now time.Time,
	device DeviceProfile,
	theme Theme,
	mode string,
	lang string,
	weekends string,
) *image.RGBA {

	img := image.NewRGBA(image.Rect(0, 0, device.Width, device.Height))
	draw.Draw(img, img.Bounds(), &image.Uniform{theme.Background}, image.Point{}, draw.Src)

	if mode != "months" {
		return img
	}

	months := BuildMonths(now, lang)

	safeTop := device.ClockBottom()
	safeBottom := device.ButtonsTop()

	// 1️⃣ Сетка занимает ВСЁ пространство между зонами
	gridTop := safeTop
	gridBottom := safeBottom
	gridHeight := gridBottom - gridTop

	drawMonths(
		img,
		months,
		device.Width,
		gridHeight,
		gridTop,
		theme,
		weekends,
	)

	// 2️⃣ Footer рисуем ПОСЛЕ сетки
	footerY := safeBottom + 172

	drawFooterAtY(
		img,
		now,
		device,
		theme,
		lang,
		footerY,
	)

	return img
}

// ===============================
// MONTH GRID
// ===============================

func drawMonths(
	img *image.RGBA,
	months []MonthData,
	width int,
	usableHeight int,
	offsetY int,
	theme Theme,
	weekends string,
) {
	const cols = 3
	const rows = 4

	cellW := width / cols
	cellH := usableHeight / rows

	for i, m := range months {
		c := i % cols
		r := i / cols

		cx := c*cellW + cellW/2
		cy := offsetY + r*cellH + cellH/2

		drawMonth(img, cx, cy, m, theme, weekends, cellW)
	}
}

func drawMonth(
	img *image.RGBA,
	cx, cy int,
	m MonthData,
	theme Theme,
	weekends string,
	cellW int,
) {
	titleColor := theme.Text
	if m.IsCurrent {
		titleColor = theme.Today
	}

	drawText(img, m.Name, cx, cy-monthTitleOffsetY, titleColor, monthFace)

	cols := 7
	spacing := cellW / 10
	radius := spacing / 4

	startX := cx - (cols-1)*spacing/2
	startY := cy - 5

	for day := 0; day < m.Days; day++ {
		col := (m.StartWeekday + day) % 7
		row := (m.StartWeekday + day) / 7

		x := startX + col*spacing
		y := startY + row*spacing

		dotColor := theme.Future

		if day < m.PassedDays {
			dotColor = theme.Active
		}

		// Weekend coloring
		if weekends != "off" && (col == 5 || col == 6) {
			switch weekends {
			case "gray":
				dotColor = theme.Future
			case "green":
				dotColor = theme.WeekendGreen
			case "blue":
				dotColor = theme.WeekendBlue
			case "red":
				dotColor = theme.WeekendRed
			}
		}

		if m.IsCurrent && day == m.PassedDays-1 {
			dotColor = theme.Today
		}

		drawCircle(img, x, y, radius, dotColor)
	}
}

// ===============================
// FOOTER (YEAR PROGRESS)
// ===============================

func drawFooterAtY(
	img *image.RGBA,
	now time.Time,
	device DeviceProfile,
	theme Theme,
	lang string,
	y int,
) {
	day := now.YearDay()
	total := DaysInYear(now.Year())
	left := total - day
	percent := int(float64(day) / float64(total) * 100)

	text := footerText(left, percent, lang)

	drawText(
		img,
		text,
		device.Width/2,
		y,
		theme.Text,
		footerFace,
	)
}

func footerText(left, percent int, lang string) string {
	if lang == "ru" {
		return fixedText("%dдн. осталось   %d%%", left, percent)
	}
	return fixedText("%dd left   %d%%", left, percent)
}

// ===============================
// DRAW HELPERS
// ===============================

func drawText(
	img *image.RGBA,
	text string,
	cx, y int,
	col color.Color,
	face font.Face,
) {
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
