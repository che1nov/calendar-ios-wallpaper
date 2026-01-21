package main

import "time"

type MonthData struct {
	Name         string
	Days         int
	PassedDays   int
	IsCurrent    bool
	StartWeekday int // 0=Mon
}

var monthNames = map[string][]string{
	"en": {"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"},
	"ru": {"Янв", "Фев", "Мар", "Апр", "Май", "Июн", "Июл", "Авг", "Сен", "Окт", "Ноя", "Дек"},
}

func BuildMonths(now time.Time, lang string) []MonthData {
	names := monthNames[lang]
	if names == nil {
		names = monthNames["en"]
	}

	year := now.Year()
	loc := now.Location()
	out := make([]MonthData, 12)

	for m := 1; m <= 12; m++ {
		first := time.Date(year, time.Month(m), 1, 0, 0, 0, 0, loc)
		days := time.Date(year, time.Month(m)+1, 0, 0, 0, 0, 0, loc).Day()

		startWeekday := (int(first.Weekday()) + 6) % 7 // Mon=0

		passed := 0
		if int(now.Month()) > m {
			passed = days
		} else if int(now.Month()) == m {
			passed = now.Day()
		}

		out[m-1] = MonthData{
			Name:         names[m-1],
			Days:         days,
			PassedDays:   passed,
			IsCurrent:    int(now.Month()) == m,
			StartWeekday: startWeekday,
		}
	}

	return out
}
