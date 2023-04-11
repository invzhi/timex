package timex

import (
	"errors"
	"fmt"
	"time"
)

// Date represents a date.
//
// The zero value of type Date is December 31 of year 0.
//
// swagger:strfmt date
type Date struct {
	ordinal int
}

// NewDate returns the date corresponding to year, month, and day.
func NewDate(year, month, day int) (Date, error) {
	if month < 1 || month > 12 {
		return Date{}, errors.New("month is out of range [1,12]")
	}

	if days := daysInMonth(year, month); day < 1 || day > days {
		return Date{}, fmt.Errorf("day is out of range [1,%d]", days)
	}

	n := toOrdinal(year, month, day)
	return Date{ordinal: n}, nil
}

// MustNewDate is like NewDate but panics if the date cannot be created.
func MustNewDate(year, month, day int) Date {
	date, err := NewDate(year, month, day)
	if err != nil {
		panic(`timex: NewDate: ` + err.Error())
	}
	return date
}

// DateFromOrdinal returns the date specified by proleptic Gregorian ordinal.
func DateFromOrdinal(n int) Date {
	return Date{ordinal: n}
}

// DateFromTime returns the date specified by t.
func DateFromTime(t time.Time) Date {
	year, month, day := t.Date()
	n := toOrdinal(year, int(month), day)
	return Date{ordinal: n}
}

// Today returns the current date in the given location.
func Today(location *time.Location) Date {
	t := time.Now().In(location)
	return DateFromTime(t)
}

// Time returns the time.Time specified by d in the given location.
func (d Date) Time(location *time.Location) time.Time {
	year, month, day := fromOrdinal(d.ordinal)
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, location)
}

// Ordinal returns the proleptic Gregorian ordinal specified by d.
func (d Date) Ordinal() int {
	return d.ordinal
}

// OrdinalDate returns the ordinal date specified by d.
func (d Date) OrdinalDate() (year, dayOfYear int) {
	year, month, day := fromOrdinal(d.ordinal)
	return year, daysBeforeMonth(year, month) + day
}

// Date returns the date specified by d.
func (d Date) Date() (year, month, day int) {
	year, month, day = fromOrdinal(d.ordinal)
	return
}

// Year returns the year specified by d.
func (d Date) Year() int {
	year, _, _ := fromOrdinal(d.ordinal)
	return year
}

// Month returns the month specified by d.
func (d Date) Month() int {
	_, month, _ := fromOrdinal(d.ordinal)
	return month
}

// Day returns the day of month specified by d.
func (d Date) Day() int {
	_, _, day := fromOrdinal(d.ordinal)
	return day
}

// DayOfYear returns the day of year specified by d.
func (d Date) DayOfYear() int {
	year, month, day := fromOrdinal(d.ordinal)
	return daysBeforeMonth(year, month) + day
}

// Weekday returns the day of week specified by d.
func (d Date) Weekday() time.Weekday {
	weekday := d.ordinal % 7 // Day 1 is monday.
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
	year, dayOfYear := thursday.OrdinalDate()
	return year, (dayOfYear-1)/7 + 1
}

// Quarter returns the quarter specified by d.
func (d Date) Quarter() int {
	_, month, _ := fromOrdinal(d.ordinal)
	return (month-1)/3 + 1
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
	year, month, day := fromOrdinal(d.ordinal)

	year += years
	month += months
	day += days

	year, month = norm(year, month, 12)
	n := ordinalBeforeYear(year)
	n += daysBeforeMonth(year, month)
	n += day

	return Date{ordinal: n}
}

// IsZero reports whether the date d is the zero of proleptic Gregorian ordinal, December 31 of year 0.
func (d Date) IsZero() bool {
	return d.ordinal == 0
}

// Before reports whether the date d is before dd.
func (d Date) Before(dd Date) bool {
	return d.ordinal < dd.ordinal
}

// After reports whether the date d is after dd.
func (d Date) After(dd Date) bool {
	return d.ordinal > dd.ordinal
}

// Equal reports whether the date d and dd is the same date.
func (d Date) Equal(dd Date) bool {
	return d.ordinal == dd.ordinal
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

// daysInMonth returns the number of days in the specified month.
func daysInMonth(year, month int) int {
	if month == 2 && isLeap(year) {
		return 29
	}
	return int(daysInYear[month] - daysInYear[month-1])
}

// daysBeforeMonth returns the number of days in the year before specified month.
func daysBeforeMonth(year, month int) int {
	days := int(daysInYear[month-1])
	if month > 2 && isLeap(year) {
		days++
	}
	return days
}

var (
	daysEvery400Years = ordinalBeforeYear(401)
	daysEvery100Years = ordinalBeforeYear(101)
	daysEvery4Years   = ordinalBeforeYear(5)
)

// ordinalBeforeYear returns the proleptic Gregorian ordinal of last year's last day.
// Ordinal day 1 is January 1 of year 1.
func ordinalBeforeYear(year int) int {
	var delta int
	if year > 0 {
		y := year - 1 // If year is 5, delta reach 1.
		delta = y/4 - y/100 + y/400
	} else {
		y := year                       // If year is -4, delta reach -1.
		delta = y/4 - y/100 + y/400 - 1 // Handle year 0, it is a leap year.
	}
	return (year-1)*365 + delta
}

// fromOrdinal returns the date of the specified proleptic Gregorian ordinal.
// Ordinal day 1 is January 1 of year 1.
func fromOrdinal(n int) (year, month, day int) {
	n400, n := norm(0, n, daysEvery400Years)
	year += n400 * 400
	// A leap day is added every 400 years, n100 will be increased incorrectly unless make a judgement here.
	if n == daysEvery400Years {
		year, month, day = year+400, 12, 31
		return
	}

	n100, n := norm(0, n, daysEvery100Years)
	year += n100 * 100

	n4, n := norm(0, n, daysEvery4Years)
	year += n4 * 4
	// A leap day is added every 4 years, n1 will be increased incorrectly unless make a judgement here.
	if n == daysEvery4Years {
		year, month, day = year+4, 12, 31
		return
	}

	n1, n := norm(0, n, 365)
	year += n1 + 1 // Start from year 1.
	return fromOrdinalDate(year, n)
}

// fromOrdinalDate returns the date of the specified ordinal date.
func fromOrdinalDate(year, dayOfYear int) (int, int, int) {
	min, max := 1, 12
	for min < max {
		m := max - (max-min)/2
		n := daysBeforeMonth(year, m)
		switch {
		case n < dayOfYear:
			min = m
		case n > dayOfYear:
			max = m - 1
		default:
			min, max = m-1, m-1
		}
	}
	month := min
	day := dayOfYear - daysBeforeMonth(year, month)
	return year, month, day
}

// toOrdinal returns the proleptic Gregorian ordinal of the specified date.
// January 1 of year 1 is ordinal day 1.
func toOrdinal(year, month, day int) int {
	n := ordinalBeforeYear(year)
	n += daysBeforeMonth(year, month)
	n += day
	return n
}
