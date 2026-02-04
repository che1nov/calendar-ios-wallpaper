package rendering

import (
	"image"
	"time"

	"calendar-wallpaper/internal/domain"
)

type Renderer struct{}

func (Renderer) RenderCalendar(
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
	return RenderCalendar(
		now,
		device,
		theme,
		mode,
		lang,
		weekends,
		dayStyle,
		uiScale,
		bgStyle,
		bgColor,
	)
}
