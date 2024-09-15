package timex_test

import (
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/invzhi/timex"
)

func BenchmarkDateOrdinalDate(b *testing.B) {
	date := timex.MustNewDate(2006, 12, 20)

	b.Run("Timex", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			date.OrdinalDate()
		}
	})
	b.Run("Time", func(b *testing.B) {
		d := date.Time(time.UTC)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			d.YearDay()
		}
	})
}

func BenchmarkDateDate(b *testing.B) {
	date := timex.MustNewDate(2006, 12, 20)

	b.Run("Timex", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			date.Date()
		}
	})
	b.Run("Time", func(b *testing.B) {
		d := date.Time(time.UTC)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			d.Date()
		}
	})
}

func BenchmarkDateYear(b *testing.B) {
	date := timex.MustNewDate(2006, 12, 20)

	b.Run("Timex", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			date.Year()
		}
	})
	b.Run("Time", func(b *testing.B) {
		d := date.Time(time.UTC)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			d.Year()
		}
	})
}

func BenchmarkDateMonth(b *testing.B) {
	date := timex.MustNewDate(2006, 12, 20)

	b.Run("Timex", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			date.Month()
		}
	})
	b.Run("Time", func(b *testing.B) {
		d := date.Time(time.UTC)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			d.Month()
		}
	})
}

func BenchmarkDateDay(b *testing.B) {
	date := timex.MustNewDate(2006, 12, 20)

	b.Run("Timex", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			date.Day()
		}
	})
	b.Run("Time", func(b *testing.B) {
		d := date.Time(time.UTC)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			d.Day()
		}
	})
}

func BenchmarkDateDayOfYear(b *testing.B) {
	date := timex.MustNewDate(2006, 12, 20)

	b.Run("Timex", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			date.DayOfYear()
		}
	})
	b.Run("Time", func(b *testing.B) {
		d := date.Time(time.UTC)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			d.YearDay()
		}
	})
}

func BenchmarkDateWeekday(b *testing.B) {
	date := timex.MustNewDate(2006, 12, 20)

	b.Run("Timex", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			date.Weekday()
		}
	})
	b.Run("Time", func(b *testing.B) {
		d := date.Time(time.UTC)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			d.Weekday()
		}
	})
}

func BenchmarkDateISOWeek(b *testing.B) {
	date := timex.MustNewDate(2006, 12, 20)

	b.Run("Timex", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			date.ISOWeek()
		}
	})
	b.Run("Time", func(b *testing.B) {
		d := date.Time(time.UTC)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			d.ISOWeek()
		}
	})
}

func BenchmarkDateAdd(b *testing.B) {
	date := timex.MustNewDate(2006, 12, 20)

	b.Run("Timex", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			date.Add(10, 100, 1000)
		}
	})
	b.Run("Time", func(b *testing.B) {
		d := date.Time(time.UTC)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			d.AddDate(10, 100, 1000)
		}
	})
}

func BenchmarkDateAddDays(b *testing.B) {
	date := timex.MustNewDate(2006, 12, 20)

	b.Run("Timex", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			date.AddDays(1000)
		}
	})
	b.Run("Time", func(b *testing.B) {
		d := date.Time(time.UTC)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			d.Add(1000 * 24 * time.Hour)
		}
	})
}

func BenchmarkDateSub(b *testing.B) {
	d1 := timex.MustNewDate(2006, 1, 2)
	d2 := timex.MustNewDate(1970, 1, 1)

	b.Run("Timex", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			d1.Sub(d2)
		}
	})
	b.Run("Time", func(b *testing.B) {
		t1 := d1.Time(time.UTC)
		t2 := d2.Time(time.UTC)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			t1.Sub(t2)
		}
	})
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
		datetime := time.Date(tt.year, time.Month(tt.month), tt.day, 0, 0, 0, 0, time.UTC)
		date := timex.DateFromTime(datetime)

		year, month, day := date.Date()
		assert.Equal(t, tt.year, year)
		assert.Equal(t, tt.month, month)
		assert.Equal(t, tt.day, day)

		assert.Equal(t, tt.year, date.Year())
		assert.Equal(t, tt.month, date.Month())
		assert.Equal(t, tt.day, date.Day())

		assert.Equal(t, datetime, date.Time(time.UTC))
	}

	t.Run("Today", func(t *testing.T) {
		now := time.Now().In(time.UTC)
		date := timex.Today(time.UTC)

		year, month, day := date.Date()
		assert.Equal(t, now.Year(), year)
		assert.Equal(t, now.Month(), time.Month(month))
		assert.Equal(t, now.Day(), day)

		assert.Equal(t, now.Year(), date.Year())
		assert.Equal(t, now.Month(), time.Month(date.Month()))
		assert.Equal(t, now.Day(), date.Day())
	})
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
		_, err := timex.NewDate(tt.year, tt.month, tt.day)
		assert.EqualError(t, err, tt.errString)

		assert.Panicsf(t, func() {
			_ = timex.MustNewDate(tt.year, tt.month, tt.day)
		}, "timex: NewDate: "+tt.errString)
	}
}

func TestDateFromOrdinalDate(t *testing.T) {
	tests := []struct {
		year, dayOfYear int
		month, day      int
	}{
		{2000, 1, 1, 1},
		{2000, 366, 12, 31},
		{2001, 1, 1, 1},
		{2001, 365, 12, 31},
	}

	for _, tt := range tests {
		date := timex.MustDateFromOrdinalDate(tt.year, tt.dayOfYear)

		year, month, day := date.Date()
		assert.Equal(t, tt.year, year)
		assert.Equal(t, tt.month, month)
		assert.Equal(t, tt.day, day)

		assert.Equal(t, tt.year, date.Year())
		assert.Equal(t, tt.month, date.Month())
		assert.Equal(t, tt.day, date.Day())

		year, dayOfYear := date.OrdinalDate()
		assert.Equal(t, tt.year, year)
		assert.Equal(t, tt.dayOfYear, dayOfYear)
		assert.Equal(t, tt.dayOfYear, date.DayOfYear())
	}

	t.Run("Errors", func(t *testing.T) {
		errTests := []struct {
			year, dayOfYear int
			errString       string
		}{
			{2000, -1, "day of year is out of range [1,366]"},
			{2000, 0, "day of year is out of range [1,366]"},
			{2000, 367, "day of year is out of range [1,366]"},
			{2001, 366, "day of year is out of range [1,365]"},
		}

		for _, tt := range errTests {
			_, err := timex.DateFromOrdinalDate(tt.year, tt.dayOfYear)
			assert.EqualError(t, err, tt.errString)

			assert.Panicsf(t, func() {
				_ = timex.MustDateFromOrdinalDate(tt.year, tt.dayOfYear)
			}, "timex: DateFromOrdinalDate: "+tt.errString)
		}
	})
}

func TestDateWeekday(t *testing.T) {
	const numOfWeeks = 1024

	// January 1 of year 1 is monday.
	weekday := time.Monday
	for n := 0; n <= 7*numOfWeeks; n++ {
		date := timex.Date{}.AddDays(n)
		assert.Equal(t, weekday, date.Weekday())

		if weekday == time.Saturday {
			weekday = time.Sunday
		} else {
			weekday++
		}
	}

	// January 1 of year 1 is monday.
	weekday = time.Monday
	for n := 0; n <= 7*numOfWeeks; n++ {
		date := timex.Date{}.AddDays(-n)
		assert.Equal(t, weekday, date.Weekday())

		if weekday == time.Sunday {
			weekday = time.Saturday
		} else {
			weekday--
		}
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
		date := timex.MustNewDate(tt.year, tt.month, tt.day)

		y, w := date.ISOWeek()
		assert.Equal(t, tt.y, y)
		assert.Equal(t, tt.w, w, tt)
	}

	// The only real invariant: Jan 04 is in week 1.
	for year := 1950; year < 2100; year++ {
		date := timex.MustNewDate(year, 1, 4)

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
		date := timex.MustNewDate(2000, tt.month, 1)

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

	d1 := timex.MustNewDate(2011, 11, 18)
	d2 := timex.MustNewDate(2016, 3, 19)

	for _, tt := range tests {
		d3 := d1.Add(tt.years, tt.months, tt.days)
		assert.Equal(t, d2, d3)
	}

	t.Run("Overflow", func(t *testing.T) {
		assert.Equal(t,
			timex.MustNewDate(2001, 3, 1),
			timex.MustNewDate(2000, 2, 29).Add(1, 0, 0),
		)
		assert.Equal(t,
			timex.MustNewDate(1999, 3, 1),
			timex.MustNewDate(2000, 2, 29).Add(-1, 0, 0),
		)
		assert.Equal(t,
			timex.MustNewDate(2000, 12, 1),
			timex.MustNewDate(2000, 10, 31).Add(0, 1, 0),
		)
		assert.Equal(t,
			timex.MustNewDate(2000, 10, 1),
			timex.MustNewDate(2000, 10, 31).Add(0, -1, 0),
		)
	})
}

func TestDateAddDays(t *testing.T) {
	tests := []struct {
		d1, d2 timex.Date
		days   int
	}{
		{timex.MustNewDate(2009, 11, 23), timex.MustNewDate(2009, 11, 24), 1},
		{timex.MustNewDate(2009, 11, 24), timex.MustNewDate(2009, 11, 23), -1},
		{timex.MustNewDate(-2009, 11, 24), timex.MustNewDate(-2009, 11, 23), -1},
		{timex.MustNewDate(2290, 1, 1), timex.MustNewDate(2000, 1, 1), -290*365 - 71},
		{timex.MustNewDate(2000, 1, 1), timex.MustNewDate(2290, 1, 1), 290*365 + 71},
	}

	for _, tt := range tests {
		d3 := tt.d1.AddDays(tt.days)
		assert.Equal(t, tt.d2, d3)
	}
}

func TestDateSub(t *testing.T) {
	maxDate := timex.Date{}.AddDays(math.MaxInt)
	minDate := timex.Date{}.AddDays(math.MinInt)

	tests := []struct {
		d1, d2 timex.Date
		days   int
	}{
		{timex.Date{}, timex.Date{}, 0},
		{minDate, maxDate, math.MinInt},
		{maxDate, minDate, math.MaxInt},
		{timex.Date{}, maxDate, -math.MaxInt},
		{maxDate, timex.Date{}, math.MaxInt},
		{timex.Date{}, minDate, math.MaxInt},
		{minDate, timex.Date{}, math.MinInt},
		{timex.MustNewDate(2009, 11, 23), timex.MustNewDate(2009, 11, 24), -1},
		{timex.MustNewDate(2009, 11, 24), timex.MustNewDate(2009, 11, 23), 1},
		{timex.MustNewDate(-2009, 11, 24), timex.MustNewDate(-2009, 11, 23), 1},
		{timex.MustNewDate(2290, 1, 1), timex.MustNewDate(2000, 1, 1), 290*365 + 71},
		{timex.MustNewDate(2000, 1, 1), timex.MustNewDate(2290, 1, 1), -290*365 - 71},
	}

	for _, tt := range tests {
		days := tt.d1.Sub(tt.d2)
		assert.Equal(t, tt.days, days)
	}
}

func TestDateIsZero(t *testing.T) {
	assert.True(t, timex.Date{}.IsZero())
	assert.True(t, timex.MustNewDate(1, 1, 1).IsZero())
	assert.True(t, timex.DateFromTime(time.Time{}).IsZero())
}

func TestDateBeforeAfter(t *testing.T) {
	tests := []struct {
		d1, d2        timex.Date
		before, after bool
	}{
		{timex.MustNewDate(2000, 1, 1), timex.MustNewDate(2000, 1, 1), false, false},

		{timex.MustNewDate(2000, 1, 1), timex.MustNewDate(2000, 1, 2), true, false},
		{timex.MustNewDate(2000, 1, 2), timex.MustNewDate(2000, 1, 1), false, true},

		{timex.MustNewDate(2000, 1, 1), timex.MustNewDate(2000, 2, 1), true, false},
		{timex.MustNewDate(2000, 2, 1), timex.MustNewDate(2000, 1, 1), false, true},

		{timex.MustNewDate(1999, 12, 31), timex.MustNewDate(2000, 1, 1), true, false},
		{timex.MustNewDate(2000, 1, 1), timex.MustNewDate(1999, 12, 31), false, true},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.before, tt.d1.Before(tt.d2))
		assert.Equal(t, tt.after, tt.d1.After(tt.d2))
		assert.Equal(t, !tt.before && !tt.after, tt.d1.Equal(tt.d2))
	}
}
