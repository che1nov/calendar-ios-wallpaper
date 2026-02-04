package domain

type DayStyle string

const (
	DayDots    DayStyle = "dots"
	DayBars    DayStyle = "bars"
	DayNumbers DayStyle = "numbers"
)

func ParseDayStyle(v string) DayStyle {
	switch DayStyle(v) {
	case DayBars:
		return DayBars
	case DayNumbers:
		return DayNumbers
	default:
		return DayDots
	}
}
