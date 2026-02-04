package usecase

import (
	"errors"
	"image"
	"time"

	"calendar-wallpaper/internal/domain"
)

type Clock interface {
	Now() time.Time
}

type SystemClock struct{}

func (SystemClock) Now() time.Time {
	return time.Now()
}

type Renderer interface {
	RenderCalendar(
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
	) *image.RGBA
}

type Service struct {
	Clock    Clock
	Renderer Renderer
	Theme    domain.Theme
}

type RenderParams struct {
	DeviceKey   string
	Lang        string
	Weekends    string
	DayStyle    string
	Timezone    int
	SizePercent int
	BgStyle     string
	BgColor     string
}

func (s Service) RenderWallpaper(p RenderParams) (*image.RGBA, error) {
	if s.Clock == nil || s.Renderer == nil {
		return nil, errors.New("service dependencies are not configured")
	}

	device, ok := domain.Devices[p.DeviceKey]
	if !ok {
		device = domain.Devices["iphone-15"]
	}

	lang := domain.NormalizeLang(p.Lang)
	weekends := normalizeWeekends(p.Weekends)
	dayStyle := domain.ParseDayStyle(p.DayStyle)
	bgStyle := domain.ParseBackgroundStyle(p.BgStyle)
	bgColor := p.BgColor
	if bgColor == "" {
		bgColor = "black"
	}

	size := p.SizePercent
	if size == 0 {
		size = 100
	}
	if size < 80 {
		size = 80
	} else if size > 130 {
		size = 130
	}
	uiScale := float64(size) / 100.0

	loc := time.FixedZone("user", p.Timezone*3600)
	now := s.Clock.Now().In(loc)

	img := s.Renderer.RenderCalendar(
		now,
		device,
		s.Theme,
		"months",
		lang,
		weekends,
		dayStyle,
		uiScale,
		bgStyle,
		bgColor,
	)
	return img, nil
}

func normalizeWeekends(v string) string {
	switch v {
	case "gray", "green", "blue", "red":
		return v
	default:
		return "off"
	}
}
