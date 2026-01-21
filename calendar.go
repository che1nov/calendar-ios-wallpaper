package main

import "time"

type MonthData struct {
	Name         string
	Days         int
	PassedDays   int
	StartWeekday int // Monday = 0
	IsCurrent    bool
}

var monthNames = map[string][]string{
	"en": {"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"},
	"ru": {"Янв", "Фев", "Мар", "Апр", "Май", "Июн", "Июл", "Авг", "Сен", "Окт", "Ноя", "Дек"},
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

		// сколько дней в месяце (ГО корректно учитывает високосные)
		daysInMonth := time.Date(year, time.Month(m)+1, 0, 0, 0, 0, 0, loc).Day()

		// Go: Sunday = 0 → ISO: Monday = 0
		startWeekday := (int(first.Weekday()) + 6) % 7

		isCurrent := int(now.Month()) == m

		passed := 0
		if isCurrent {
			passed = now.Day()
		} else if int(now.Month()) > m {
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
