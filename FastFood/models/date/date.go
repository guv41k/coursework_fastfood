package date

import (
	"fmt"
	"strconv"
	"strings"
)

type Date struct {
	Day   int
	Month int
	Year  int
}

func Parse(s string) (Date, bool) {
	parts := strings.Split(s, ".")
	if len(parts) != 3 {
		return Date{}, false
	}

	day, err1 := strconv.Atoi(parts[0])
	month, err2 := strconv.Atoi(parts[1])
	year, err3 := strconv.Atoi(parts[2])
	if err1 != nil || err2 != nil || err3 != nil {
		return Date{}, false
	}

	d := Date{
		Day:   day,
		Month: month,
		Year:  year,
	}
	if !d.IsValid() {
		return Date{}, false
	}
	return d, true
}

func (d Date) String() string {
	return fmt.Sprintf("%02d.%02d.%04d", d.Day, d.Month, d.Year)
}

func (d Date) Before(o Date) bool {
	if d.Year != o.Year {
		return d.Year < o.Year
	}
	if d.Month != o.Month {
		return d.Month < o.Month
	}
	return d.Day < o.Day
}

func (d Date) IsValid() bool {
	if d.Year < 1 || d.Month < 1 || d.Month > 12 || d.Day < 1 {
		return false
	}
	return d.Day <= daysInMonth(d.Month, d.Year)
}

func daysInMonth(month, year int) int {
	switch month {
	case 1, 3, 5, 7, 8, 10, 12:
		return 31
	case 4, 6, 9, 11:
		return 30
	case 2:
		if year%4 == 0 && (year%100 != 0 || year%400 == 0) {
			return 29
		}
		return 28
	}
	return 0
}
