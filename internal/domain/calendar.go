package domain

import "time"

type MonthData struct {
	Name         string
	Days         int
	PassedDays   int
	IsCurrent    bool
	StartWeekday int
}

var monthNames = map[string][]string{
	"en": {"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"},
	"ru": {"Янв", "Фев", "Мар", "Апр", "Май", "Июн", "Июл", "Авг", "Сен", "Окт", "Ноя", "Дек"},
}

func NormalizeLang(lang string) string {
	if _, ok := monthNames[lang]; ok {
		return lang
	}
	return "en"
}

func Progress(t time.Time) (day, left, percent int) {
	day = t.YearDay()
	total := 365
	if isLeap(t.Year()) {
		total = 366
	}
	left = total - day
	percent = int(float64(day) / float64(total) * 100)
	return
}

func isLeap(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

func BuildMonths(now time.Time, lang string) []MonthData {
	year := now.Year()
	loc := now.Location()

	names := monthNames[NormalizeLang(lang)]
	months := make([]MonthData, 12)

	for m := 1; m <= 12; m++ {
		first := time.Date(year, time.Month(m), 1, 0, 0, 0, 0, loc)
		days := first.AddDate(0, 1, -1).Day()

		weekday := (int(first.Weekday()) + 6) % 7

		passed := 0
		if int(now.Month()) > m {
			passed = days
		} else if int(now.Month()) == m {
			passed = now.Day()
		}

		months[m-1] = MonthData{
			Name:         names[m-1],
			Days:         days,
			PassedDays:   passed,
			IsCurrent:    int(now.Month()) == m,
			StartWeekday: weekday,
		}
	}
	return months
}

func DaysInYear(year int) int {
	if time.Date(year, time.December, 31, 0, 0, 0, 0, time.UTC).YearDay() == 366 {
		return 366
	}
	return 365
}
