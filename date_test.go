package timex

import (
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func BenchmarkDateOrdinalDate(b *testing.B) {
	date := MustNewDate(2006, 1, 2)
	for i := 0; i < b.N; i++ {
		date.OrdinalDate()
	}
}

func BenchmarkDateDate(b *testing.B) {
	date := MustNewDate(2006, 1, 2)
	for i := 0; i < b.N; i++ {
		date.Date()
	}
}

func BenchmarkDateWeekday(b *testing.B) {
	date := MustNewDate(2006, 1, 2)
	for i := 0; i < b.N; i++ {
		date.Weekday()
	}
}

func BenchmarkDateISOWeek(b *testing.B) {
	date := MustNewDate(2006, 1, 2)
	for i := 0; i < b.N; i++ {
		date.ISOWeek()
	}
}

func BenchmarkDateQuarter(b *testing.B) {
	date := MustNewDate(2006, 1, 2)
	for i := 0; i < b.N; i++ {
		date.Quarter()
	}
}

func BenchmarkDateAdd(b *testing.B) {
	var date Date
	for i := 0; i < b.N; i++ {
		date.Add(10, 100, 3650000)
	}
}

func TestLeapYear(t *testing.T) {
	for year := -10000; year <= 10000; year++ {
		days := ordinalBeforeYear(year+1) - ordinalBeforeYear(year)

		if isLeap(year) {
			assert.Equal(t, 366, days)
		} else {
			assert.Equal(t, 365, days)
		}
	}
}

func TestDaysInMonth(t *testing.T) {
	tests := []struct {
		year, month int
		days        int
	}{
		{-1, 1, 31},
		{-1, 2, 28}, // Non-leap year
		{-1, 3, 31},
		{-1, 4, 30},
		{-1, 5, 31},
		{-1, 6, 30},
		{-1, 7, 31},
		{-1, 8, 31},
		{-1, 9, 30},
		{-1, 10, 31},
		{-1, 11, 30},
		{-1, 12, 31},
		{0, 1, 31},
		{0, 2, 29}, // Leap year
		{0, 3, 31},
		{0, 4, 30},
		{0, 5, 31},
		{0, 6, 30},
		{0, 7, 31},
		{0, 8, 31},
		{0, 9, 30},
		{0, 10, 31},
		{0, 11, 30},
		{0, 12, 31},
		{1, 1, 31},
		{1, 2, 28}, // Non-leap year
		{1, 3, 31},
		{1, 4, 30},
		{1, 5, 31},
		{1, 6, 30},
		{1, 7, 31},
		{1, 8, 31},
		{1, 9, 30},
		{1, 10, 31},
		{1, 11, 30},
		{1, 12, 31},
		{4, 1, 31},
		{4, 2, 29}, // Leap year
		{4, 3, 31},
		{4, 4, 30},
		{4, 5, 31},
		{4, 6, 30},
		{4, 7, 31},
		{4, 8, 31},
		{4, 9, 30},
		{4, 10, 31},
		{4, 11, 30},
		{4, 12, 31},
	}

	for _, tt := range tests {
		days := daysInMonth(tt.year, tt.month)
		assert.Equal(t, tt.days, days)
	}
}

func TestOrdinalBeforeYear(t *testing.T) {
	assert.Equal(t, 365*400+97, daysEvery400Years)
	assert.Equal(t, 365*100+24, daysEvery100Years)
	assert.Equal(t, 365*4+1, daysEvery4Years)

	assert.Equal(t, -367, ordinalBeforeYear(0))
	assert.Equal(t, -1, ordinalBeforeYear(1))
	assert.Equal(t, 364, ordinalBeforeYear(2))

	// Days in positive years and negative years should be equal.
	for year := 1; year < 10000; year++ {
		daysInPositiveYears := ordinalBeforeYear(year) + 1
		daysInNegativeYears := ordinalBeforeYear(0) - ordinalBeforeYear(-year+1)

		assert.Equal(t, daysInPositiveYears, daysInNegativeYears)
	}

	// Returned ordinal day should be the last day of last year.
	for year := 1; year < 10000; year++ {
		n := ordinalBeforeYear(year)

		y1, m1, d1 := fromOrdinal(n)
		assert.Equal(t, year-1, y1)
		assert.Equal(t, 12, m1)
		assert.Equal(t, 31, d1)

		y2, m2, d2 := fromOrdinal(n + 1)
		assert.Equal(t, year, y2)
		assert.Equal(t, 1, m2)
		assert.Equal(t, 1, d2)
	}
}

func TestOrdinalFromTo(t *testing.T) {
	assert.Equal(t, -1, toOrdinal(0, 12, 31))
	assert.Equal(t, 0, toOrdinal(1, 1, 1))
	assert.Equal(t, 1, toOrdinal(1, 1, 2))

	for n := -3650000; n <= 3650000; n++ {
		year, month, day := fromOrdinal(n)

		assert.GreaterOrEqual(t, month, 1)
		assert.LessOrEqual(t, month, 12)

		assert.GreaterOrEqual(t, day, 1)
		assert.LessOrEqual(t, day, daysInMonth(year, month))

		assert.Equal(t, n, toOrdinal(year, month, day))
	}
}

func TestDateFromTime(t *testing.T) {
	tests := []struct {
		year, month, day int
	}{
		{-1, 1, 1},
		{0, 1, 1},
		{0, 12, 31},
		{1, 1, 1},
		{2006, 1, 2},
	}

	for _, tt := range tests {
		dateTime := time.Date(tt.year, time.Month(tt.month), tt.day, 0, 0, 0, 0, time.UTC)
		date := DateFromTime(dateTime)

		year, month, day := date.Date()
		assert.Equal(t, tt.year, year)
		assert.Equal(t, tt.month, month)
		assert.Equal(t, tt.day, day)

		assert.Equal(t, dateTime, date.Time(time.UTC))
	}

	t.Run("Today", func(t *testing.T) {
		now := time.Now()
		date := Today(time.UTC)

		assert.Equal(t, now.Year(), date.Year())
		assert.Equal(t, now.YearDay(), date.DayOfYear())
	})
}

func TestDateOrdinalDate(t *testing.T) {
	for year := -10000; year <= 10000; year++ {
		days := 365
		if isLeap(year) {
			days = 366
		}
		for dayOfYear := 1; dayOfYear <= days; dayOfYear++ {
			{
				y, m, d := fromOrdinalDate(year, dayOfYear)
				y, yearDay := MustNewDate(y, m, d).OrdinalDate()
				assert.Equal(t, year, y)
				assert.Equal(t, dayOfYear, yearDay)
			}
			{
				n := ordinalBeforeYear(year) + dayOfYear
				assert.Equal(t, dayOfYear, Date{ordinal: n}.DayOfYear())
			}
		}
	}
}

func TestNewDateErrors(t *testing.T) {
	tests := []struct {
		year, month, day int
		errString        string
	}{
		{2000, -1, 1, "month is out of range [1,12]"},
		{2000, 0, 1, "month is out of range [1,12]"},
		{2000, 13, 1, "month is out of range [1,12]"},
		{2000, 2, -1, "day is out of range [1,29]"},
		{2000, 2, 30, "day is out of range [1,29]"},
		{2001, 2, -1, "day is out of range [1,28]"},
		{2001, 2, 29, "day is out of range [1,28]"},
	}

	for _, tt := range tests {
		_, err := NewDate(tt.year, tt.month, tt.day)
		assert.EqualError(t, err, tt.errString)

		assert.Panicsf(t, func() {
			_ = MustNewDate(tt.year, tt.month, tt.day)
		}, "timex: NewDate: "+tt.errString)
	}
}

func TestDateWeekday(t *testing.T) {
	// January 1 of year 1 is monday.
	n := 0
	for weekday := time.Monday; n <= 3650000; weekday++ {
		if weekday > time.Saturday {
			weekday = time.Sunday
		}
		assert.Equal(t, weekday, Date{ordinal: n}.Weekday())

		n++
	}
	n = 0
	for weekday := time.Monday; n >= -3650000; weekday-- {
		if weekday < time.Sunday {
			weekday = time.Saturday
		}
		assert.Equal(t, weekday, Date{ordinal: n}.Weekday())

		n--
	}
}

func TestDateISOWeek(t *testing.T) {
	tests := []struct {
		year  int
		month int
		day   int
		y     int
		w     int
	}{
		{1981, 1, 1, 1981, 1},
		{1982, 1, 1, 1981, 53},
		{1983, 1, 1, 1982, 52},
		{1984, 1, 1, 1983, 52},
		{1985, 1, 1, 1985, 1},
		{1986, 1, 1, 1986, 1},
		{1987, 1, 1, 1987, 1},
		{1988, 1, 1, 1987, 53},
		{1989, 1, 1, 1988, 52},
		{1990, 1, 1, 1990, 1},
		{1991, 1, 1, 1991, 1},
		{1992, 1, 1, 1992, 1},
		{1993, 1, 1, 1992, 53},
		{1994, 1, 1, 1993, 52},
		{1995, 1, 2, 1995, 1},
		{1996, 1, 1, 1996, 1},
		{1996, 1, 7, 1996, 1},
		{1996, 1, 8, 1996, 2},
		{1997, 1, 1, 1997, 1},
		{1998, 1, 1, 1998, 1},
		{1999, 1, 1, 1998, 53},
		{2000, 1, 1, 1999, 52},
		{2001, 1, 1, 2001, 1},
		{2002, 1, 1, 2002, 1},
		{2003, 1, 1, 2003, 1},
		{2004, 1, 1, 2004, 1},
		{2005, 1, 1, 2004, 53},
		{2006, 1, 1, 2005, 52},
		{2007, 1, 1, 2007, 1},
		{2008, 1, 1, 2008, 1},
		{2009, 1, 1, 2009, 1},
		{2010, 1, 1, 2009, 53},
		{2010, 1, 1, 2009, 53},
		{2011, 1, 1, 2010, 52},
		{2011, 1, 2, 2010, 52},
		{2011, 1, 3, 2011, 1},
		{2011, 1, 4, 2011, 1},
		{2011, 1, 5, 2011, 1},
		{2011, 1, 6, 2011, 1},
		{2011, 1, 7, 2011, 1},
		{2011, 1, 8, 2011, 1},
		{2011, 1, 9, 2011, 1},
		{2011, 1, 10, 2011, 2},
		{2011, 1, 11, 2011, 2},
		{2011, 6, 12, 2011, 23},
		{2011, 6, 13, 2011, 24},
		{2011, 12, 25, 2011, 51},
		{2011, 12, 26, 2011, 52},
		{2011, 12, 27, 2011, 52},
		{2011, 12, 28, 2011, 52},
		{2011, 12, 29, 2011, 52},
		{2011, 12, 30, 2011, 52},
		{2011, 12, 31, 2011, 52},
		{1995, 1, 1, 1994, 52},
		{2012, 1, 1, 2011, 52},
		{2012, 1, 2, 2012, 1},
		{2012, 1, 8, 2012, 1},
		{2012, 1, 9, 2012, 2},
		{2012, 12, 23, 2012, 51},
		{2012, 12, 24, 2012, 52},
		{2012, 12, 30, 2012, 52},
		{2012, 12, 31, 2013, 1},
		{2013, 1, 1, 2013, 1},
		{2013, 1, 6, 2013, 1},
		{2013, 1, 7, 2013, 2},
		{2013, 12, 22, 2013, 51},
		{2013, 12, 23, 2013, 52},
		{2013, 12, 29, 2013, 52},
		{2013, 12, 30, 2014, 1},
		{2014, 1, 1, 2014, 1},
		{2014, 1, 5, 2014, 1},
		{2014, 1, 6, 2014, 2},
		{2015, 1, 1, 2015, 1},
		{2016, 1, 1, 2015, 53},
		{2017, 1, 1, 2016, 52},
		{2018, 1, 1, 2018, 1},
		{2019, 1, 1, 2019, 1},
		{2020, 1, 1, 2020, 1},
		{2021, 1, 1, 2020, 53},
		{2022, 1, 1, 2021, 52},
		{2023, 1, 1, 2022, 52},
		{2024, 1, 1, 2024, 1},
		{2025, 1, 1, 2025, 1},
		{2026, 1, 1, 2026, 1},
		{2027, 1, 1, 2026, 53},
		{2028, 1, 1, 2027, 52},
		{2029, 1, 1, 2029, 1},
		{2030, 1, 1, 2030, 1},
		{2031, 1, 1, 2031, 1},
		{2032, 1, 1, 2032, 1},
		{2033, 1, 1, 2032, 53},
		{2034, 1, 1, 2033, 52},
		{2035, 1, 1, 2035, 1},
		{2036, 1, 1, 2036, 1},
		{2037, 1, 1, 2037, 1},
		{2038, 1, 1, 2037, 53},
		{2039, 1, 1, 2038, 52},
		{2040, 1, 1, 2039, 52},
	}

	for _, tt := range tests {
		date, err := NewDate(tt.year, tt.month, tt.day)
		assert.NoError(t, err)

		y, w := date.ISOWeek()
		assert.Equal(t, tt.y, y)
		assert.Equal(t, tt.w, w, tt)
	}

	// The only real invariant: Jan 04 is in week 1
	for year := 1950; year < 2100; year++ {
		date, err := NewDate(year, 1, 4)
		assert.NoError(t, err)

		y, w := date.ISOWeek()
		assert.Equal(t, year, y)
		assert.Equal(t, 1, w)
	}
}

func TestDateQuarter(t *testing.T) {
	tests := []struct {
		month   int
		quarter int
	}{
		{1, 1},
		{2, 1},
		{3, 1},
		{4, 2},
		{5, 2},
		{6, 2},
		{7, 3},
		{8, 3},
		{9, 3},
		{10, 4},
		{11, 4},
		{12, 4},
	}
	for _, tt := range tests {
		date, err := NewDate(2000, tt.month, 1)
		assert.NoError(t, err)

		assert.Equal(t, tt.quarter, date.Quarter())
	}
}

func TestDateAdd(t *testing.T) {
	tests := []struct {
		years, months, days int
	}{
		{4, 4, 1},
		{3, 16, 1},
		{3, 15, 30},
		{5, -6, -18 - 30 - 12},
		{5, -5, -18 - 31 - 30 - 12},
		{6, -17, -18 - 31 - 30 - 12},
	}

	d1, err := NewDate(2011, 11, 18)
	assert.NoError(t, err)
	d2, err := NewDate(2016, 3, 19)
	assert.NoError(t, err)

	for _, tt := range tests {
		d3 := d1.Add(tt.years, tt.months, tt.days)
		assert.Equal(t, d2, d3)
	}

	t.Run("Overflow", func(t *testing.T) {
		assert.Equal(t,
			MustNewDate(2001, 3, 1),
			MustNewDate(2000, 2, 29).Add(1, 0, 0),
		)
		assert.Equal(t,
			MustNewDate(1999, 3, 1),
			MustNewDate(2000, 2, 29).Add(-1, 0, 0),
		)
		assert.Equal(t,
			MustNewDate(2000, 12, 1),
			MustNewDate(2000, 10, 31).Add(0, 1, 0),
		)
		assert.Equal(t,
			MustNewDate(2000, 10, 1),
			MustNewDate(2000, 10, 31).Add(0, -1, 0),
		)
	})
}

func TestDateSub(t *testing.T) {
	maxDate := Date{ordinal: math.MaxInt}
	minDate := Date{ordinal: math.MinInt}

	tests := []struct {
		d1, d2 Date
		days   int
	}{
		{Date{}, Date{}, 0},
		{minDate, maxDate, math.MinInt},
		{maxDate, minDate, math.MaxInt},
		{Date{}, maxDate, -math.MaxInt},
		{maxDate, Date{}, math.MaxInt},
		{Date{}, minDate, math.MaxInt},
		{minDate, Date{}, math.MinInt},
		{MustNewDate(2009, 11, 23), MustNewDate(2009, 11, 24), -1},
		{MustNewDate(2009, 11, 24), MustNewDate(2009, 11, 23), 1},
		{MustNewDate(-2009, 11, 24), MustNewDate(-2009, 11, 23), 1},
		{MustNewDate(2290, 1, 1), MustNewDate(2000, 1, 1), 290*365 + 71},
		{MustNewDate(2000, 1, 1), MustNewDate(2290, 1, 1), -290*365 - 71},
	}

	for _, tt := range tests {
		days := tt.d1.Sub(tt.d2)
		assert.Equal(t, tt.days, days)
	}
}

func TestDateIsZero(t *testing.T) {
	assert.True(t, Date{}.IsZero())
	assert.True(t, MustNewDate(1, 1, 1).IsZero())
	assert.True(t, DateFromTime(time.Time{}).IsZero())
}

func TestDateBeforeAfter(t *testing.T) {
	tests := []struct {
		d1, d2        Date
		before, after bool
	}{
		{MustNewDate(2000, 1, 1), MustNewDate(2000, 1, 1), false, false},

		{MustNewDate(2000, 1, 1), MustNewDate(2000, 1, 2), true, false},
		{MustNewDate(2000, 1, 2), MustNewDate(2000, 1, 1), false, true},

		{MustNewDate(2000, 1, 1), MustNewDate(2000, 2, 1), true, false},
		{MustNewDate(2000, 2, 1), MustNewDate(2000, 1, 1), false, true},

		{MustNewDate(1999, 12, 31), MustNewDate(2000, 1, 1), true, false},
		{MustNewDate(2000, 1, 1), MustNewDate(1999, 12, 31), false, true},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.before, tt.d1.Before(tt.d2))
		assert.Equal(t, tt.after, tt.d1.After(tt.d2))
		assert.Equal(t, !tt.before && !tt.after, tt.d1.Equal(tt.d2))
	}
}
