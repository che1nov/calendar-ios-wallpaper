package main

import "time"

type MonthData struct {
	Name       string
	Days       int
	PassedDays int
	IsCurrent  bool
}

var monthNames = map[string][]string{
	"en": {
		"Jan", "Feb", "Mar",
		"Apr", "May", "Jun",
		"Jul", "Aug", "Sep",
		"Oct", "Nov", "Dec",
	},
	"ru": {
		"Янв", "Фев", "Мар",
		"Апр", "Май", "Июн",
		"Июл", "Авг", "Сен",
		"Окт", "Ноя", "Дек",
	},
}

func DaysInYear(year int) int {
	if time.Date(year, 12, 31, 0, 0, 0, 0, time.UTC).YearDay() == 366 {
		return 366
	}
	return 365
}

func Progress(t time.Time) (day, left, percent int) {
	day = t.YearDay()
	total := DaysInYear(t.Year())
	left = total - day
	percent = int(float64(day) / float64(total) * 100)
	return
}

func BuildMonths(now time.Time, lang string) []MonthData {
	year := now.Year()
	loc := now.Location()

	names, ok := monthNames[lang]
	if !ok {
		names = monthNames["en"]
	}

	months := make([]MonthData, 12)

	for m := 1; m <= 12; m++ {
		first := time.Date(year, time.Month(m), 1, 0, 0, 0, 0, loc)
		days := first.AddDate(0, 1, -1).Day()

		passed := 0
		if int(now.Month()) > m {
			passed = days
		} else if int(now.Month()) == m {
			passed = now.Day()
		}

		months[m-1] = MonthData{
			Name:       names[m-1],
			Days:       days,
			PassedDays: passed,
			IsCurrent:  int(now.Month()) == m,
		}
	}

	return months
}
