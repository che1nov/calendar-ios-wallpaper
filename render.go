package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"time"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

const (
	monthTitleOffsetY = 40
)

var (
	monthFace  font.Face
	footerFace font.Face
	numberFace font.Face
)

func init() {
	fontBytes := mustRead("fonts/Inter-Regular.ttf")
	f := mustParseFont(fontBytes)

	monthFace = mustFace(f, 38)
	footerFace = mustFace(f, 30)
	numberFace = mustFace(f, 28)
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
	dayStyle DayStyle,
) *image.RGBA {

	img := image.NewRGBA(image.Rect(0, 0, device.Width, device.Height))
	draw.Draw(img, img.Bounds(), &image.Uniform{theme.Background}, image.Point{}, draw.Src)

	if mode != "months" {
		return img
	}

	months := BuildMonths(now, lang)

	safeTop := device.ClockBottom()
	safeBottom := device.ButtonsTop()

	// Сетка занимает ВСЁ пространство между зонами
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
		dayStyle, // ← ВАЖНО
	)

	// Footer рисуем ПОСЛЕ сетки
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
	dayStyle DayStyle,
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

		drawMonth(img, cx, cy, m, theme, weekends, cellW, dayStyle)

	}
}

func drawMonth(
	img *image.RGBA,
	cx, cy int,
	m MonthData,
	theme Theme,
	weekends string,
	cellW int,
	style DayStyle,
) {
	titleColor := theme.Text
	if m.IsCurrent {
		titleColor = theme.Today
	}

	drawText(img, m.Name, cx, cy-monthTitleOffsetY, titleColor, monthFace)

	switch style {
	case DayDots:
		drawMonthDots(img, cx, cy, m, theme, weekends, cellW)
	case DayBars:
		drawMonthBars(img, cx, cy, m, theme, weekends, cellW)
	case DayNumbers:
		drawMonthNumbers(img, cx, cy, m, theme, weekends, cellW)
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

func drawRect(
	img *image.RGBA,
	x, y int,
	w, h int,
	col color.Color,
) {
	for dy := 0; dy < h; dy++ {
		for dx := 0; dx < w; dx++ {
			img.Set(x+dx, y+dy, col)
		}
	}
}

func drawMonthDots(
	img *image.RGBA,
	cx, cy int,
	m MonthData,
	theme Theme,
	weekends string,
	cellW int,
) {
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

		color := resolveDayColor(day, col, m, theme, weekends)
		drawCircle(img, x, y, radius, color)
	}
}

func drawMonthBars(
	img *image.RGBA,
	cx, cy int,
	m MonthData,
	theme Theme,
	weekends string,
	cellW int,
) {
	cols := 7
	spacing := cellW / 10
	barW := spacing * 2
	barH := spacing / 4

	startX := cx - (cols-1)*spacing/2
	startY := cy - barH/2

	for day := 0; day < m.Days; day++ {
		col := (m.StartWeekday + day) % 7
		row := (m.StartWeekday + day) / 7

		x := startX + col*spacing
		y := startY + row*spacing

		color := resolveDayColor(day, col, m, theme, weekends)

		drawRect(img, x-barW/2, y-barH/2, barW, barH, color)
	}
}

func drawMonthNumbers(
	img *image.RGBA,
	cx, cy int,
	m MonthData,
	theme Theme,
	weekends string,
	cellW int,
) {
	cols := 7
	spacing := cellW / 9

	startX := cx - (cols-1)*spacing/2
	startY := cy - 5

	for day := 0; day < m.Days; day++ {
		col := (m.StartWeekday + day) % 7
		row := (m.StartWeekday + day) / 7

		x := startX + col*spacing
		y := startY + row*spacing + numberFace.Metrics().Ascent.Round()/2

		color := resolveDayColor(day, col, m, theme, weekends)

		drawText(
			img,
			fmt.Sprintf("%d", day+1),
			x,
			y,
			color,
			numberFace, // размер ~18–20
		)
	}
}

func resolveDayColor(
	day int,
	col int,
	m MonthData,
	theme Theme,
	weekends string,
) color.Color {

	if m.IsCurrent && day == m.PassedDays-1 {
		return theme.Today
	}

	if day < m.PassedDays {
		return theme.Active
	}

	if weekends != "off" && (col == 5 || col == 6) {
		switch weekends {
		case "gray":
			return theme.Future
		case "green":
			return theme.WeekendGreen
		case "blue":
			return theme.WeekendBlue
		case "red":
			return theme.WeekendRed
		}
	}

	return theme.Future
}
