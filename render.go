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

// ===============================
// BASE (iPhone 15 / 15 Pro)
// ===============================

const (
	BaseWidth  = 1179
	BaseHeight = 2556

	BaseMonthFont  = 38
	BaseFooterFont = 30
	BaseNumberFont = 22

	BaseDotRadius = 6
	BaseSpacing   = 32
)

// ===============================
// FONTS (scaled per device)
// ===============================

var (
	monthFace  font.Face
	footerFace font.Face
	numberFace font.Face
)

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

	scale := device.Scale()

	initFonts(scale)

	img := image.NewRGBA(image.Rect(0, 0, device.Width, device.Height))
	draw.Draw(img, img.Bounds(), &image.Uniform{theme.Background}, image.Point{}, draw.Src)

	if mode != "months" {
		return img
	}

	months := BuildMonths(now, lang)

	safeTop := device.ClockBottom()
	safeBottom := device.ButtonsTop()

	// ===== FOOTER METRICS =====
	footerHeight := int(48 * scale)
	footerGap := int(32 * scale)

	gridTop := safeTop
	gridBottom := safeBottom - footerHeight - footerGap
	gridHeight := gridBottom - gridTop

	drawMonths(
		img,
		months,
		device,
		gridHeight,
		gridTop,
		theme,
		weekends,
		dayStyle,
	)

	footerY := device.Height - int(72*scale)
	drawFooterAtY(img, now, device, theme, lang, footerY)

	return img
}

// ===============================
// FONTS INIT
// ===============================

func initFonts(scale float64) {
	fontBytes := mustRead("fonts/Inter-Regular.ttf")
	f := mustParseFont(fontBytes)

	monthFace = mustFace(f, BaseMonthFont*scale)
	footerFace = mustFace(f, BaseFooterFont*scale)
	numberFace = mustFace(f, BaseNumberFont*scale)
}

// ===============================
// MONTH GRID
// ===============================

func drawMonths(
	img *image.RGBA,
	months []MonthData,
	device DeviceProfile,
	usableHeight int,
	offsetY int,
	theme Theme,
	weekends string,
	dayStyle DayStyle,
) {
	const cols = 3
	const rows = 4

	cellW := device.Width / cols
	cellH := usableHeight / rows

	for i, m := range months {
		c := i % cols
		r := i / cols

		cx := c*cellW + cellW/2
		cy := offsetY + r*cellH + cellH/2

		drawMonth(img, cx, cy, m, device, theme, weekends, dayStyle, cellW)
	}
}

// ===============================
// SINGLE MONTH
// ===============================

func drawMonth(
	img *image.RGBA,
	cx, cy int,
	m MonthData,
	device DeviceProfile,
	theme Theme,
	weekends string,
	style DayStyle,
	cellW int,
) {
	scale := device.Scale()

	titleOffset := int(40 * scale)

	titleColor := theme.Text
	if m.IsCurrent {
		titleColor = theme.Today
	}

	drawText(img, m.Name, cx, cy-titleOffset, titleColor, monthFace)

	switch style {
	case DayDots:
		drawMonthDots(img, cx, cy, m, device, theme, weekends, cellW)
	case DayBars:
		drawMonthBars(img, cx, cy, m, device, theme, weekends, cellW)
	case DayNumbers:
		drawMonthNumbers(img, cx, cy, m, device, theme, weekends, cellW)
	}
}

// ===============================
// DAY STYLES
// ===============================

func drawMonthDots(
	img *image.RGBA,
	cx, cy int,
	m MonthData,
	device DeviceProfile,
	theme Theme,
	weekends string,
	cellW int,
) {
	scale := device.Scale()

	cols := 7
	spacing := int(BaseSpacing * scale)
	radius := int(BaseDotRadius * scale)

	startX := cx - (cols-1)*spacing/2
	startY := cy

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
	device DeviceProfile,
	theme Theme,
	weekends string,
	cellW int,
) {
	scale := device.Scale()

	cols := 7
	spacing := int(BaseSpacing * scale)

	barW := int(16 * scale)
	barH := int(6 * scale)

	startX := cx - (cols-1)*spacing/2
	startY := cy

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
	device DeviceProfile,
	theme Theme,
	weekends string,
	cellW int,
) {
	scale := device.Scale()

	cols := 7
	spacing := int(30 * scale)

	startX := cx - (cols-1)*spacing/2
	startY := cy

	for day := 0; day < m.Days; day++ {
		col := (m.StartWeekday + day) % 7
		row := (m.StartWeekday + day) / 7

		x := startX + col*spacing
		y := startY + row*spacing

		color := resolveDayColor(day, col, m, theme, weekends)

		drawText(
			img,
			fmt.Sprintf("%d", day+1),
			x,
			y,
			color,
			numberFace,
		)
	}
}

// ===============================
// FOOTER
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

	drawText(img, text, device.Width/2, y, theme.Text, footerFace)
}

func footerText(left, percent int, lang string) string {
	if lang == "ru" {
		return fixedText("%d дн. осталось   %d%%", left, percent)
	}
	return fixedText("%d d left   %d%%", left, percent)
}

// ===============================
// HELPERS
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

func drawRect(img *image.RGBA, x, y, w, h int, col color.Color) {
	for dy := 0; dy < h; dy++ {
		for dx := 0; dx < w; dx++ {
			img.Set(x+dx, y+dy, col)
		}
	}
}

func resolveDayColor(
	day int,
	col int,
	m MonthData,
	theme Theme,
	weekends string,
) color.Color {

	// Сегодня — абсолютный приоритет
	if m.IsCurrent && day == m.PassedDays-1 {
		return theme.Today
	}

	// Прошедшие дни
	if day < m.PassedDays {
		return theme.Active
	}

	// Выходные
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

	// Будущие дни
	return theme.Future
}
