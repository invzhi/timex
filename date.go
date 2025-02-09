package timex

import (
	"errors"
	"fmt"
	"math"
	"time"
)

// Date represents a date.
//
// The zero value of type Date is January 1 of year 1.
//
// swagger:strfmt date
type Date struct {
	ordinal int // ordinal represents days since January 1 of year 1.
}

// NewDate returns the date corresponding to year, month, and day.
func NewDate(year, month, day int) (Date, error) {
	if month < 1 || month > 12 {
		return Date{}, errors.New("month is out of range [1,12]")
	}

	if days := daysInMonth(year, month); day < 1 || day > days {
		return Date{}, fmt.Errorf("day is out of range [1,%d]", days)
	}

	n := calendarToOrdinal(year, month, day)
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

// DateFromOrdinalDate returns the date corresponding to year, day of year.
func DateFromOrdinalDate(year, dayOfYear int) (Date, error) {
	days := 365
	if isLeap(year) {
		days++
	}

	if dayOfYear < 1 || dayOfYear > days {
		return Date{}, fmt.Errorf("day of year is out of range [1,%d]", days)
	}

	n := ordinalDateToOrdinal(year, dayOfYear)
	return Date{ordinal: n}, nil
}

// MustDateFromOrdinalDate is like DateFromOrdinalDate but panics if the date cannot be created.
func MustDateFromOrdinalDate(year, dayOfYear int) Date {
	date, err := DateFromOrdinalDate(year, dayOfYear)
	if err != nil {
		panic(`timex: DateFromOrdinalDate: ` + err.Error())
	}
	return date
}

// DateFromTime returns the date specified by t.
func DateFromTime(t time.Time) Date {
	year, month, day := t.Date()
	n := calendarToOrdinal(year, int(month), day)
	return Date{ordinal: n}
}

// Today returns the current date in the given location.
func Today(location *time.Location) Date {
	t := time.Now().In(location)
	return DateFromTime(t)
}

// Time returns the time.Time specified by d in the given location.
func (d Date) Time(location *time.Location) time.Time {
	year, month, day := ordinalToCalendar(d.ordinal)
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, location)
}

// OrdinalDate returns the ordinal date specified by d.
func (d Date) OrdinalDate() (year, dayOfYear int) {
	return ordinalToOrdinalDate(d.ordinal)
}

// Date returns the year, month and day specified by d.
func (d Date) Date() (year, month, day int) {
	year, month, day = ordinalToCalendar(d.ordinal)
	return
}

// Year returns the year specified by d.
func (d Date) Year() int {
	year, _ := ordinalToOrdinalDate(d.ordinal)
	return year
}

// Quarter returns the quarter specified by d.
func (d Date) Quarter() int {
	_, month, _ := ordinalToCalendar(d.ordinal)
	return (month-1)/3 + 1
}

// Month returns the month specified by d.
func (d Date) Month() int {
	_, month, _ := ordinalToCalendar(d.ordinal)
	return month
}

// Day returns the day of month specified by d.
func (d Date) Day() int {
	_, _, day := ordinalToCalendar(d.ordinal)
	return day
}

// DayOfYear returns the day of year specified by d.
func (d Date) DayOfYear() int {
	_, dayOfYear := ordinalToOrdinalDate(d.ordinal)
	return dayOfYear
}

// Weekday returns the day of week specified by d.
func (d Date) Weekday() time.Weekday {
	weekday := (d.ordinal + 1) % 7 // Day 0 is monday.
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

	thursday := d.ordinal + delta
	year, dayOfYear := ordinalToOrdinalDate(thursday)
	return year, (dayOfYear-1)/7 + 1
}

// norm1 normalize the hi and lo into [1, base].
func norm1(hi, lo, base int) (int, int) {
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
	year, month, day := ordinalToCalendar(d.ordinal)

	year += years
	month += months
	day += days

	year, month = norm1(year, month, 12)
	n := ordinalBeforeYear(year)
	n += daysBeforeMonth(year, month)
	n += day

	return Date{ordinal: n}
}

// AddDays returns the date corresponding to adding the given number of days to d.
func (d Date) AddDays(days int) Date {
	return Date{ordinal: d.ordinal + days}
}

// Sub returns the days d-dd.
// If the result exceeds the integer scope, the maximum (or minimum) integer will be returned.
func (d Date) Sub(dd Date) int {
	days := d.ordinal - dd.ordinal
	switch {
	case d.ordinal >= 0 && dd.ordinal <= 0 && days < 0:
		return math.MaxInt
	case d.ordinal <= 0 && dd.ordinal >= 0 && days > 0:
		return math.MinInt
	default:
		return days
	}
}

// IsZero reports whether the date d is the zero value, January 1 of year 1.
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
	daysEvery400Years = ordinalBeforeYear(401) + 1
	daysEvery100Years = ordinalBeforeYear(101) + 1
	daysEvery4Years   = ordinalBeforeYear(5) + 1
)

// ordinalBeforeYear returns the ordinal of last year's last day.
// Ordinal day 0 is January 1 of year 1.
func ordinalBeforeYear(year int) int {
	var delta int
	if year > 0 {
		y := year - 1 // If year is 5, delta reach 1.
		delta = y/4 - y/100 + y/400
	} else {
		y := year                       // If year is -4, delta reach -1.
		delta = y/4 - y/100 + y/400 - 1 // Handle year 0, it is a leap year.
	}
	return (year-1)*365 + delta - 1 // Shift for January 1 of year 1.
}

func ordinalDateToOrdinal(year, dayOfYear int) int {
	return ordinalBeforeYear(year) + dayOfYear
}

func calendarToOrdinal(year, month, day int) int {
	dayOfYear := daysBeforeMonth(year, month) + day
	return ordinalDateToOrdinal(year, dayOfYear)
}

func ordinalToOrdinalDate(n int) (year, dayOfYear int) {
	n400, n := norm1(0, n, daysEvery400Years)
	n400, n = norm1(n400, n+1, daysEvery400Years) // Shift for January 1 of year 1, prevent integer overflow.
	year = n400*400 + 1                           // Start from year 1.

	n100 := (n - 1) / daysEvery100Years
	n100 -= n100 >> 2 // Handle the leap day every 400 years. If n100 is 4, set it to 3.
	year += n100 * 100
	n -= daysEvery100Years * n100

	n4 := (n - 1) / daysEvery4Years
	year += n4 * 4
	n -= daysEvery4Years * n4

	n1 := (n - 1) / 365
	n1 -= n1 >> 2 // Handle the leap day every 4 years. If n1 is 4, set it to 3.
	year += n1
	n -= 365 * n1

	return year, n
}

func ordinalDateToCalendar(year, dayOfYear int) (int, int, int) {
	day := dayOfYear
	if isLeap(year) {
		switch {
		case day == 31+29:
			return year, 2, 29 // Leap day.
		case day > 31+29:
			day-- // Remove leap day.
		}
	}

	month := (dayOfYear-1)/31 + 1

	days := int(daysInYear[month])
	if day > days {
		month++
		day -= days
	} else {
		day -= int(daysInYear[month-1])
	}

	return year, month, day
}

func ordinalToCalendar(n int) (year, month, day int) {
	year, dayOfYear := ordinalToOrdinalDate(n)
	return ordinalDateToCalendar(year, dayOfYear)
}
