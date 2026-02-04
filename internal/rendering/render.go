package rendering

import (
	"container/list"
	"fmt"
	"image"
	"image/color"
	"sync"
	"time"

	"calendar-wallpaper/internal/domain"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

const (
	BaseWidth  = 1179
	BaseHeight = 2556

	BaseMonthFont  = 38
	BaseFooterFont = 30
	BaseNumberFont = 22

	BaseDotRadius = 6
	BaseSpacing   = 32
)

const (
	MonthTitleScale = 1.00
	DayGridScale    = 1.00
	FooterScale     = 1.00
)

type FontSet struct {
	Month  font.Face
	Footer font.Face
	Number font.Face
}

const maxFontCacheEntries = 64

type fontCacheEntry struct {
	key   string
	faces FontSet
}

var fontCache struct {
	sync.Mutex
	items map[string]*list.Element
	order *list.List
}

func RenderCalendar(
	now time.Time,
	device domain.DeviceProfile,
	theme domain.Theme,
	mode string,
	lang string,
	weekends string,
	dayStyle domain.DayStyle,
	uiScale float64,
	bgStyle domain.BackgroundStyle,
	bgColor string,
) *image.RGBA {

	deviceScale := float64(device.Width) / float64(BaseWidth)
	scale := deviceScale * uiScale

	faces := getFontSet(scale)

	img := image.NewRGBA(image.Rect(0, 0, device.Width, device.Height))
	drawBackground(img, device, bgStyle, bgColor)

	if mode != "months" {
		return img
	}

	months := domain.BuildMonths(now, lang)

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
		faces,
	)

	footerY := device.ButtonsTop() + int(80*scale)
	drawFooterAtY(img, now, device, theme, lang, footerY, faces)

	return img
}

func getFontSet(scale float64) FontSet {
	key := fmt.Sprintf("%.4f", scale)

	fontCache.Lock()
	defer fontCache.Unlock()

	if fontCache.items == nil {
		fontCache.items = make(map[string]*list.Element)
		fontCache.order = list.New()
	}

	if el, ok := fontCache.items[key]; ok {
		fontCache.order.MoveToFront(el)
		return el.Value.(fontCacheEntry).faces
	}

	fontBytes := mustRead("fonts/SFPRODISPLAYBOLD.OTF")
	f := mustParseFont(fontBytes)

	faces := FontSet{
		Month:  mustFace(f, BaseMonthFont*scale*MonthTitleScale),
		Footer: mustFace(f, BaseFooterFont*scale*FooterScale),
		Number: mustFace(f, BaseNumberFont*scale*DayGridScale),
	}
	el := fontCache.order.PushFront(fontCacheEntry{key: key, faces: faces})
	fontCache.items[key] = el

	if fontCache.order.Len() > maxFontCacheEntries {
		oldest := fontCache.order.Back()
		if oldest != nil {
			fontCache.order.Remove(oldest)
			oldKey := oldest.Value.(fontCacheEntry).key
			delete(fontCache.items, oldKey)
		}
	}
	return faces
}

func drawMonths(
	img *image.RGBA,
	months []domain.MonthData,
	device domain.DeviceProfile,
	usableHeight int,
	offsetY int,
	theme domain.Theme,
	weekends string,
	dayStyle domain.DayStyle,
	scale float64,
	faces FontSet,
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
			faces,
		)
	}
}

func drawMonth(
	img *image.RGBA,
	cx, cy int,
	m domain.MonthData,
	theme domain.Theme,
	weekends string,
	style domain.DayStyle,
	cellW int,
	scale float64,
	faces FontSet,
) {
	titleOffset := int(52 * scale)

	titleColor := theme.Text
	if m.IsCurrent {
		titleColor = theme.Today
	}

	drawText(img, m.Name, cx, cy-titleOffset, titleColor, faces.Month)

	switch style {
	case domain.DayDots:
		drawMonthDots(img, cx, cy, m, theme, weekends, scale)
	case domain.DayBars:
		drawMonthBars(img, cx, cy, m, theme, weekends, scale)
	case domain.DayNumbers:
		drawMonthNumbers(img, cx, cy, m, theme, weekends, scale, faces)
	}
}

func drawMonthDots(
	img *image.RGBA,
	cx, cy int,
	m domain.MonthData,
	theme domain.Theme,
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
	m domain.MonthData,
	theme domain.Theme,
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
	m domain.MonthData,
	theme domain.Theme,
	weekends string,
	scale float64,
	faces FontSet,
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
			faces.Number,
		)
	}
}

func drawFooterAtY(
	img *image.RGBA,
	now time.Time,
	device domain.DeviceProfile,
	theme domain.Theme,
	lang string,
	y int,
	faces FontSet,
) {
	day := now.YearDay()
	total := domain.DaysInYear(now.Year())
	left := total - day
	percent := int(float64(day) / float64(total) * 100)

	drawText(
		img,
		footerText(left, percent, lang),
		device.Width/2,
		y,
		theme.Text,
		faces.Footer,
	)
}

func footerText(left, percent int, lang string) string {
	if lang == "ru" {
		return fmt.Sprintf("%d дн. осталось   %d%%", left, percent)
	}
	return fmt.Sprintf("%d d left   %d%%", left, percent)
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

func drawRect(img *image.RGBA, x, y, w, h int, col color.Color) {
	for dy := 0; dy < h; dy++ {
		for dx := 0; dx < w; dx++ {
			img.Set(x+dx, y+dy, col)
		}
	}
}

func resolveDayColor(day, col int, m domain.MonthData, theme domain.Theme, weekends string) color.Color {
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
