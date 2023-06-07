package timex

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func BenchmarkDateParse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = ParseDate("YYYY-MM-DD", "2006-01-02")
	}
}

func BenchmarkDateFormat(b *testing.B) {
	date := MustNewDate(2006, 1, 2)
	for i := 0; i < b.N; i++ {
		_ = date.Format(RFC3339)
	}
}

func BenchmarkDateString(b *testing.B) {
	date := MustNewDate(2006, 1, 2)
	for i := 0; i < b.N; i++ {
		_ = date.String()
	}
}

func BenchmarkDateGoString(b *testing.B) {
	date := MustNewDate(2006, 1, 2)
	for i := 0; i < b.N; i++ {
		_ = date.GoString()
	}
}

func BenchmarkDateMarshalJSON(b *testing.B) {
	date := MustNewDate(2006, 1, 2)
	for i := 0; i < b.N; i++ {
		_, _ = date.MarshalJSON()
	}
}

func TestAppendInt(t *testing.T) {
	tests := []struct {
		n     int
		width int
		s     string
	}{
		{-12, 0, "-12"},
		{-12, 1, "-12"},
		{-12, 2, "-12"},
		{-12, 3, "-012"},
		{-12, 4, "-0012"},
		{-1, 0, "-1"},
		{-1, 1, "-1"},
		{-1, 2, "-01"},
		{-1, 3, "-001"},
		{0, 0, "0"},
		{0, 1, "0"},
		{0, 2, "00"},
		{0, 3, "000"},
		{1, 0, "1"},
		{1, 1, "1"},
		{1, 2, "01"},
		{1, 3, "001"},
		{12, 0, "12"},
		{12, 1, "12"},
		{12, 2, "12"},
		{12, 3, "012"},
		{12, 4, "0012"},
		{123, 0, "123"},
		{123, 1, "123"},
		{123, 2, "123"},
		{123, 3, "123"},
		{123, 4, "0123"},
	}

	for _, tt := range tests {
		bytes := appendInt(nil, tt.n, tt.width)
		assert.Equal(t, tt.s, string(bytes))
	}
}

func TestParseDate(t *testing.T) {
	tests := []struct {
		layout string
		value  string
		unsafe bool
	}{
		{RFC3339, "2010-02-04", false},
		{"MMMM DD YYYY", "February 04 2010", false},
		{"MMMM D, YYYY", "February 4, 2010", false},
		{"MMM DD YYYY", "Feb 04 2010", false},
		{"MMM D YYYY", "Feb 4 2010", false},
		{"DD MMM YYYY", "04 Feb 2010", false},
		{"DD-MMM-YY", "04-Feb-10", true},
		// Case-insensitive
		{"MMM D YYYY", "FEB 4 2010", false},
		// Chinese
		{"YYYY年M月D日", "2010年2月4日", false},
	}

	for _, tt := range tests {
		date, err := ParseDate(tt.layout, tt.value)
		assert.NoError(t, err)

		assert.Equal(t, 2010, date.Year())
		assert.Equal(t, 2, date.Month())
		assert.Equal(t, 4, date.Day())
	}

	t.Run("Format", func(t *testing.T) {
		dates := []Date{
			MustNewDate(2010, 2, 4),
			MustNewDate(1990, 2, 4), // Nineties year
		}

		for _, tt := range tests {
			for _, d1 := range dates {
				d2, err := ParseDate(tt.layout, d1.Format(tt.layout))
				assert.NoError(t, err)

				assert.Equal(t, d1, d2)
			}
		}
	})

	t.Run("FormatUnsafe", func(t *testing.T) {
		dates := []Date{
			MustNewDate(0, 2, 4),    // Zero year
			MustNewDate(-123, 2, 4), // Negative year
		}

		for _, tt := range tests {
			for _, d1 := range dates {
				if tt.unsafe {
					continue
				}

				d2, err := ParseDate(tt.layout, d1.Format(tt.layout))
				assert.NoError(t, err)

				assert.Equal(t, d1, d2)
			}
		}
	})

	t.Run("FormatBigYear", func(t *testing.T) {
		d1 := MustNewDate(-12345, 2, 4) // Big Negative year
		d2 := MustNewDate(12345, 2, 4)  // Big Positive year

		assert.Equal(t, "-12345-02-04", d1.Format(RFC3339))
		assert.Equal(t, "12345-02-04", d2.Format(RFC3339))
	})
}

func TestParseDateErrors(t *testing.T) {
	tests := []struct {
		layout    string
		value     string
		errString string
	}{
		{RFC3339, "22-10-25", `parsing date "22-10-25" as "YYYY-MM-DD": cannot parse "22-10-25" as "YYYY"`},
		{RFC3339, "12022-10-25", `parsing date "12022-10-25" as "YYYY-MM-DD": cannot parse "2-10-25" as "-"`},
		{" YYYY-MM-DD", "2010-02-04", `parsing date "2010-02-04" as " YYYY-MM-DD": cannot parse "2010-02-04" as " "`},
		{" YYYY-MM-DD", "", `parsing date "" as " YYYY-MM-DD": cannot parse "" as "YYYY"`},
		{"YY-MM-DD", "a2-10-25", `parsing date "a2-10-25" as "YY-MM-DD": cannot parse "a2-10-25" as "YY"`},
		{"YY-M-DD", "22-a0-25", `parsing date "22-a0-25" as "YY-M-DD": cannot parse "a0-25" as "M"`},
		{"D MMM YY", "4 --- 00", `parsing date "4 --- 00" as "D MMM YY": cannot parse "--- 00" as "MMM"`},
		{"D MMMM YY", "4 --- 00", `parsing date "4 --- 00" as "D MMMM YY": cannot parse "--- 00" as "MMMM"`},
	}

	for _, tt := range tests {
		_, err := ParseDate(tt.layout, tt.value)
		assert.EqualError(t, err, tt.errString)
	}
}

func FuzzParseDate(f *testing.F) {
	f.Add("YYYY-MM-DD", "2006-01-02")
	f.Add(" YYYY-MM-DD", "")
	f.Fuzz(func(t *testing.T, layout, value string) {
		assert.NotPanics(t, func() {
			_, _ = ParseDate(layout, value)
		})
	})
}

func FuzzDateFormat(f *testing.F) {
	f.Add(0, "YYYY-MM-DD")
	f.Add(0, " YYYY-MM-DD")
	f.Fuzz(func(t *testing.T, n int, layout string) {
		assert.NotPanics(t, func() {
			date := Date{ordinal: n}
			_ = date.Format(layout)
		})
	})
}

func TestDateString(t *testing.T) {
	tests := []struct {
		year, month, day int
		str, goStr       string
	}{
		{2006, 1, 2, "2006-01-02", "timex.MustNewDate(2006, 1, 2)"},
		{1, 12, 11, "0001-12-11", "timex.MustNewDate(1, 12, 11)"},
		{0, 1, 2, "0000-01-02", "timex.MustNewDate(0, 1, 2)"},
		{-2000, 1, 2, "-2000-01-02", "timex.MustNewDate(-2000, 1, 2)"},
		{10001, 1, 2, "10001-01-02", "timex.MustNewDate(10001, 1, 2)"},
		{-10001, 1, 2, "-10001-01-02", "timex.MustNewDate(-10001, 1, 2)"},
	}

	for _, tt := range tests {
		date := MustNewDate(tt.year, tt.month, tt.day)
		assert.Equal(t, tt.str, date.String())
		assert.Equal(t, tt.goStr, date.GoString())
	}
}

func TestDateMarshalJSON(t *testing.T) {
	d1 := MustNewDate(2006, 1, 2)
	bytes, err := d1.MarshalJSON()
	assert.NoError(t, err)

	var d2 Date
	err = d2.UnmarshalJSON(bytes)
	assert.NoError(t, err)
	assert.Equal(t, d1, d2)

	var d3 Date
	err = d3.UnmarshalJSON([]byte("null"))
	assert.NoError(t, err)
	assert.True(t, d3.IsZero())
}

func TestDateMarshalJSONError(t *testing.T) {
	tests := []struct {
		year, month, day int
		errString        string
	}{
		{-1, 1, 2, "year is out of range [0,9999]"},
		{-102, 1, 2, "year is out of range [0,9999]"},
		{-1000, 1, 2, "year is out of range [0,9999]"},
		{10000, 1, 2, "year is out of range [0,9999]"},
		{10200, 1, 2, "year is out of range [0,9999]"},
	}
	for _, tt := range tests {
		date := MustNewDate(tt.year, tt.month, tt.day)
		_, err := date.MarshalJSON()
		assert.EqualError(t, err, tt.errString)
	}
}
