package main

import (
	"fmt"
	"image"
	"image/color"
	"time"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

/*
BASE — дизайн под iPhone 15 / 15 Pro
*/

const (
	BaseWidth  = 1179
	BaseHeight = 2556

	BaseMonthFont  = 38
	BaseFooterFont = 30
	BaseNumberFont = 22

	BaseDotRadius = 6
	BaseSpacing   = 32
)

/*
DESIGN TUNING
*/
const (
	MonthTitleScale = 1.00 // названия месяцев
	DayGridScale    = 1.00 // точки / бары / числа
	FooterScale     = 1.00 // нижний текст
)

/*
FONTS
*/
var (
	monthFace  font.Face
	footerFace font.Face
	numberFace font.Face
)

/*
MAIN ENTRY
*/
func RenderCalendar(
	now time.Time,
	device DeviceProfile,
	theme Theme,
	mode string,
	lang string,
	weekends string,
	dayStyle DayStyle,
	uiScale float64,
	bgStyle BackgroundStyle,
	bgColor BackgroundColor,
) *image.RGBA {

	deviceScale := device.Scale()
	scale := deviceScale * uiScale

	initFonts(scale)

	img := image.NewRGBA(image.Rect(0, 0, device.Width, device.Height))
	drawBackground(img, device, bgStyle, bgColor)

	if mode != "months" {
		return img
	}

	months := BuildMonths(now, lang)

	safeTop := device.ClockBottom()
	safeBottom := device.ButtonsTop()

	footerHeight := int(52 * scale)
	footerGap := int(28 * scale)

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
		scale,
	)

	footerY := device.ButtonsTop() + int(80*scale)
	drawFooterAtY(img, now, device, theme, lang, footerY)

	return img
}

/*
FONTS INIT
*/
func initFonts(scale float64) {
	fontBytes := mustRead("fonts/SFPRODISPLAYBOLD.OTF")
	f := mustParseFont(fontBytes)

	monthFace = mustFace(f, BaseMonthFont*scale*MonthTitleScale)
	footerFace = mustFace(f, BaseFooterFont*scale*FooterScale)
	numberFace = mustFace(f, BaseNumberFont*scale*DayGridScale)
}

/*
MONTH GRID
*/
func drawMonths(
	img *image.RGBA,
	months []MonthData,
	device DeviceProfile,
	usableHeight int,
	offsetY int,
	theme Theme,
	weekends string,
	dayStyle DayStyle,
	scale float64,
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

		drawMonth(
			img,
			cx,
			cy,
			m,
			theme,
			weekends,
			dayStyle,
			cellW,
			scale,
		)
	}
}

/*
SINGLE MONTH
*/
func drawMonth(
	img *image.RGBA,
	cx, cy int,
	m MonthData,
	theme Theme,
	weekends string,
	style DayStyle,
	cellW int,
	scale float64,
) {
	titleOffset := int(52 * scale)

	titleColor := theme.Text
	if m.IsCurrent {
		titleColor = theme.Today
	}

	drawText(img, m.Name, cx, cy-titleOffset, titleColor, monthFace)

	switch style {
	case DayDots:
		drawMonthDots(img, cx, cy, m, theme, weekends, scale)
	case DayBars:
		drawMonthBars(img, cx, cy, m, theme, weekends, scale)
	case DayNumbers:
		drawMonthNumbers(img, cx, cy, m, theme, weekends, scale)
	}
}

/*
DAY STYLES
*/
func drawMonthDots(
	img *image.RGBA,
	cx, cy int,
	m MonthData,
	theme Theme,
	weekends string,
	scale float64,
) {
	gridScale := scale * DayGridScale

	cols := 7
	spacing := int(BaseSpacing * gridScale)
	radius := int(BaseDotRadius * gridScale)

	startX := cx - (cols-1)*spacing/2
	startY := cy

	for day := 0; day < m.Days; day++ {
		col := (m.StartWeekday + day) % 7
		row := (m.StartWeekday + day) / 7

		x := startX + col*spacing
		y := startY + row*spacing

		drawCircle(img, x, y, radius, resolveDayColor(day, col, m, theme, weekends))
	}
}

func drawMonthBars(
	img *image.RGBA,
	cx, cy int,
	m MonthData,
	theme Theme,
	weekends string,
	scale float64,
) {
	gridScale := scale * DayGridScale

	cols := 7
	spacing := int(BaseSpacing * gridScale)

	barW := int(20 * gridScale)
	barH := int(6 * gridScale)

	startX := cx - (cols-1)*spacing/2
	startY := cy

	for day := 0; day < m.Days; day++ {
		col := (m.StartWeekday + day) % 7
		row := (m.StartWeekday + day) / 7

		x := startX + col*spacing
		y := startY + row*spacing

		drawRect(img, x-barW/2, y-barH/2, barW, barH,
			resolveDayColor(day, col, m, theme, weekends))
	}
}

func drawMonthNumbers(
	img *image.RGBA,
	cx, cy int,
	m MonthData,
	theme Theme,
	weekends string,
	scale float64,
) {
	gridScale := scale * DayGridScale

	cols := 7
	spacing := int(30 * gridScale)

	startX := cx - (cols-1)*spacing/2
	startY := cy

	for day := 0; day < m.Days; day++ {
		col := (m.StartWeekday + day) % 7
		row := (m.StartWeekday + day) / 7

		x := startX + col*spacing
		y := startY + row*spacing

		drawText(
			img,
			fmt.Sprintf("%d", day+1),
			x,
			y,
			resolveDayColor(day, col, m, theme, weekends),
			numberFace,
		)
	}
}

/*
FOOTER
*/
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

	drawText(
		img,
		footerText(left, percent, lang),
		device.Width/2,
		y,
		theme.Text,
		footerFace,
	)
}

func footerText(left, percent int, lang string) string {
	if lang == "ru" {
		return fmt.Sprintf("%d дн. осталось   %d%%", left, percent)
	}
	return fmt.Sprintf("%d d left   %d%%", left, percent)
}

/*
DRAW HELPERS
*/
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

func drawRect(img *image.RGBA, x, y, w, h int, col color.Color) {
	for dy := 0; dy < h; dy++ {
		for dx := 0; dx < w; dx++ {
			img.Set(x+dx, y+dy, col)
		}
	}
}

func resolveDayColor(day, col int, m MonthData, theme Theme, weekends string) color.Color {
	if m.IsCurrent && day == m.PassedDays-1 {
		return theme.Today
	}
	if day < m.PassedDays {
		return theme.Active
	}
	if weekends != "off" && (col == 5 || col == 6) {
		switch weekends {
		case "gray":
			return theme.WeekendGray
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
