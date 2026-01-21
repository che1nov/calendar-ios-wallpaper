package main

import "time"

type MonthData struct {
	Name         string
	Days         int
	StartWeekday int // Monday = 0
	PassedDays   int
	IsCurrent    bool
}

func BuildMonths(now time.Time, lang string) []MonthData {
	year := now.Year()
	loc := now.Location()

	names := monthNames[lang]
	if names == nil {
		names = monthNames["en"]
	}

	months := make([]MonthData, 12)

	for m := 1; m <= 12; m++ {
		first := time.Date(year, time.Month(m), 1, 0, 0, 0, 0, loc)
		daysInMonth := first.AddDate(0, 1, -1).Day()

		// Monday = 0
		startWeekday := (int(first.Weekday()) + 6) % 7

		isCurrent := int(now.Month()) == m

		passed := 0
		if isCurrent {
			passed = now.Day()
		} else if now.Month() > time.Month(m) {
			passed = daysInMonth
		}

		months[m-1] = MonthData{
			Name:         names[m-1],
			Days:         daysInMonth,
			StartWeekday: startWeekday,
			PassedDays:   passed,
			IsCurrent:    isCurrent,
		}
	}

	return months
}

var monthNames = map[string][]string{
	"en": {"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"},
	"ru": {"Янв", "Фев", "Мар", "Апр", "Май", "Июн", "Июл", "Авг", "Сен", "Окт", "Ноя", "Дек"},
}
