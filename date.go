package timex

import (
	"errors"
	"fmt"
	"time"
)

// Date represents a date with day precision.
type Date struct {
	year  int
	month int
	day   int
}

// NewDate returns the date corresponding to year, month, and day.
func NewDate(year, month, day int) (Date, error) {
	if month < 1 || month > 12 {
		return Date{}, errors.New("month is out of range [1,12]")
	}

	if days := daysInMonth(year, month); day < 1 || day > days {
		return Date{}, fmt.Errorf("day is out of range [1,%d]", days)
	}

	return Date{year: year, month: month, day: day}, nil
}

// MustNewDate is like NewDate but panics if the date cannot be created.
func MustNewDate(year, month, day int) Date {
	date, err := NewDate(year, month, day)
	if err != nil {
		panic(`timex: NewDate: ` + err.Error())
	}
	return date
}

// DateFromTime returns the date specified by t.
func DateFromTime(t time.Time) Date {
	year, month, day := t.Date()
	return Date{year: year, month: int(month), day: day}
}

// Today returns the current date in the given location.
func Today(location *time.Location) Date {
	t := time.Now().In(location)
	return DateFromTime(t)
}

// Time returns the time.Time specified by d in the given location.
func (d Date) Time(location *time.Location) time.Time {
	return time.Date(d.year, time.Month(d.month), d.day, 0, 0, 0, 0, location)
}

// Year returns the year specified by d.
func (d Date) Year() int { return d.year }

// Month returns the month specified by d.
func (d Date) Month() int { return d.month }

// Day returns the day of month specified by d.
func (d Date) Day() int { return d.day }

// Quarter returns the quarter specified by d.
func (d Date) Quarter() int {
	return (d.month-1)/3 + 1
}

// DayOfYear returns the day of year specified by d.
func (d Date) DayOfYear() int {
	return daysBeforeMonth(d.year, d.month) + d.day
}

// Weekday returns the day of week specified by d.
func (d Date) Weekday() time.Weekday {
	weekday := ordinalDay(d.year, d.month, d.day) % 7
	if weekday < 0 {
		weekday += 7
	}
	return time.Weekday(weekday)
}

// ISOWeek returns the ISO 8601 year and week number specified by d.
func (d Date) ISOWeek() (year, week int) {
	delta := int(time.Thursday - d.Weekday())
	if delta == 4 { // Sunday
		delta = -3
	}

	thursday := d.Add(0, 0, delta)
	return thursday.year, (thursday.DayOfYear()-1)/7 + 1
}

// norm normalize the hi and lo into [1, base].
func norm(hi, lo, base int) (int, int) {
	if lo < 1 {
		n := -(lo/base - 1)
		lo += n * base
		hi -= n
	}
	if lo > base {
		n := (lo - 1) / base
		lo -= n * base
		hi += n
	}
	return hi, lo
}

// Add returns the date corresponding to adding the given number of years, months, and days to d.
func (d Date) Add(years, months, days int) Date {
	var (
		year  = d.year + years
		month = d.month + months
		day   = d.day + days
	)

	year, month = norm(year, month, 12)

	for day < 1 {
		year, month = norm(year, month-1, 12)
		day += daysInMonth(year, month)
	}

	n := daysInMonth(year, month)
	for day > n {
		day -= n
		year, month = norm(year, month+1, 12)

		n = daysInMonth(year, month)
	}

	return Date{year: year, month: month, day: day}
}

// Before reports whether the date d is before dd.
func (d Date) Before(dd Date) bool {
	return d.year < dd.year ||
		d.year == dd.year && d.month < dd.month ||
		d.year == dd.year && d.month == dd.month && d.day < dd.day
}

// After reports whether the date d is after dd.
func (d Date) After(dd Date) bool {
	return d.year > dd.year ||
		d.year == dd.year && d.month > dd.month ||
		d.year == dd.year && d.month == dd.month && d.day > dd.day
}

// Equal reports whether the date d and dd is the same date.
func (d Date) Equal(dd Date) bool {
	return d.year == dd.year && d.month == dd.month && d.day == dd.day
}

// daysInYear[month] counts the number of days in a non-leap year when the month ends.
var daysInYear = [...]int32{
	0,
	31,
	31 + 28,
	31 + 28 + 31,
	31 + 28 + 31 + 30,
	31 + 28 + 31 + 30 + 31,
	31 + 28 + 31 + 30 + 31 + 30,
	31 + 28 + 31 + 30 + 31 + 30 + 31,
	31 + 28 + 31 + 30 + 31 + 30 + 31 + 31,
	31 + 28 + 31 + 30 + 31 + 30 + 31 + 31 + 30,
	31 + 28 + 31 + 30 + 31 + 30 + 31 + 31 + 30 + 31,
	31 + 28 + 31 + 30 + 31 + 30 + 31 + 31 + 30 + 31 + 30,
	31 + 28 + 31 + 30 + 31 + 30 + 31 + 31 + 30 + 31 + 30 + 31,
}

func isLeap(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

func daysInMonth(year, month int) int {
	if month == 2 && isLeap(year) {
		return 29
	}
	return int(daysInYear[month] - daysInYear[month-1])
}

// ordinalDayBeforeYear returns the ordinal day of last year's last day.
// Day 1 is Jan 1 of year 1.
func ordinalDayBeforeYear(year int) int {
	// If year is 5, delta reach 1; year is -4, delta reach -1.
	y := year
	if year > 0 {
		y--
	}

	delta := y/4 - y/100 + y/400
	// Handle year 0, it is a leap year.
	if year <= 0 {
		delta--
	}
	return (year-1)*365 + delta
}

// daysBeforeMonth returns the number of days in the year before month.
func daysBeforeMonth(year, month int) int {
	days := int(daysInYear[month-1])
	if month > 2 && isLeap(year) {
		days++
	}
	return days
}

// ordinalDay returns the ordinal day of the specified date.
// Day 1 is Jan 1 of year 1.
func ordinalDay(year, month, day int) int {
	n := ordinalDayBeforeYear(year)
	n += daysBeforeMonth(year, month)
	n += day
	return n
}
