package timex

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewDateErrors(t *testing.T) {
	tests := []struct {
		year, month, day int
		errString        string
	}{
		{-1, 1, 1, "year is out of range [0,9999]"},
		{10000, 1, 1, "year is out of range [0,9999]"},
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

func TestDateFromTime(t *testing.T) {
	tt := time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)
	for i := 0; tt.Year() < 10000; i++ {
		date := DateFromTime(tt)
		assert.Equal(t, tt.Year(), date.Year())
		assert.Equal(t, int(tt.Month()), date.Month())
		assert.Equal(t, tt.Day(), date.Day())
		assert.Equal(t, tt.YearDay(), date.DayOfYear())
		assert.Equal(t, tt.Weekday(), date.Weekday())

		tt = tt.AddDate(0, 0, 1)
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

func TestDateAdd(t *testing.T) {
	tests := []struct {
		years, months, days int
	}{
		{4, 4, 1},
		{3, 16, 1},
		{3, 15, 30},
		{5, -6, -18 - 30 - 12},
	}

	d1, err := NewDate(2011, 11, 18)
	assert.NoError(t, err)
	d2, err := NewDate(2016, 3, 19)
	assert.NoError(t, err)

	for _, tt := range tests {
		d3 := d1.Add(tt.years, tt.months, tt.days)
		assert.Equal(t, d3, d2)
	}
}
