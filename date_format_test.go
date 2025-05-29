package timex_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/invzhi/timex"
)

func BenchmarkDateParse(b *testing.B) {
	value := "02/01/2006"

	b.Run("Timex", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := timex.ParseDate("DD/MM/YYYY", value)
			assert.NoError(b, err)
		}
	})
	b.Run("Time", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := time.ParseInLocation("02/01/2006", value, time.UTC)
			assert.NoError(b, err)
		}
	})
}

func BenchmarkDateFormat(b *testing.B) {
	date := timex.MustNewDate(2006, 1, 2)

	b.Run("Timex", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			date.Format("YYYY-MM-DD")
		}
	})
	b.Run("Time", func(b *testing.B) {
		d := date.Time(time.UTC)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			d.Format("2006-01-02")
		}
	})
}

func BenchmarkDateString(b *testing.B) {
	date := timex.MustNewDate(2006, 1, 2)
	for i := 0; i < b.N; i++ {
		_ = date.String()
	}
}

func BenchmarkDateGoString(b *testing.B) {
	date := timex.MustNewDate(2006, 1, 2)
	for i := 0; i < b.N; i++ {
		_ = date.GoString()
	}
}

func BenchmarkDateMarshalJSON(b *testing.B) {
	date := timex.MustNewDate(2006, 1, 2)

	b.Run("Timex", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := date.MarshalJSON()
			assert.NoError(b, err)
		}
	})
	b.Run("Time", func(b *testing.B) {
		d := date.Time(time.UTC)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := d.MarshalJSON()
			assert.NoError(b, err)
		}
	})
}

func BenchmarkDateUnmarshalJSON(b *testing.B) {
	b.Run("Timex", func(b *testing.B) {
		var date timex.Date

		for i := 0; i < b.N; i++ {
			err := date.UnmarshalJSON([]byte(`"2006-01-02"`))
			assert.NoError(b, err)
		}
	})
	b.Run("Time", func(b *testing.B) {
		bytes, _ := time.Date(2006, 1, 2, 0, 0, 0, 0, time.UTC).MarshalJSON()

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var d time.Time
			err := d.UnmarshalJSON(bytes)
			assert.NoError(b, err)
		}
	})
}

func TestParseDate(t *testing.T) {
	tests := []struct {
		layout string
		value  string
		unsafe bool
	}{
		{timex.RFC3339Date, "2010-02-04", false},
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
		date, err := timex.ParseDate(tt.layout, tt.value)
		assert.NoError(t, err)

		assert.Equal(t, 2010, date.Year())
		assert.Equal(t, 2, date.Month())
		assert.Equal(t, 4, date.Day())
	}

	t.Run("Format", func(t *testing.T) {
		dates := []timex.Date{
			timex.MustNewDate(2010, 2, 4),
			timex.MustNewDate(1990, 2, 4), // Nineties year
		}

		for _, tt := range tests {
			for _, d1 := range dates {
				d2, err := timex.ParseDate(tt.layout, d1.Format(tt.layout))
				assert.NoError(t, err)

				assert.Equal(t, d1, d2)
			}
		}
	})

	t.Run("FormatUnsafe", func(t *testing.T) {
		dates := []timex.Date{
			timex.MustNewDate(0, 2, 4),    // Zero year
			timex.MustNewDate(-123, 2, 4), // Negative year
		}

		for _, tt := range tests {
			for _, d1 := range dates {
				if tt.unsafe {
					continue
				}

				d2, err := timex.ParseDate(tt.layout, d1.Format(tt.layout))
				assert.NoError(t, err)

				assert.Equal(t, d1, d2)
			}
		}
	})

	t.Run("FormatBigYear", func(t *testing.T) {
		d1 := timex.MustNewDate(-12345, 2, 4) // Big Negative year
		d2 := timex.MustNewDate(12345, 2, 4)  // Big Positive year

		assert.Equal(t, "-12345-02-04", d1.Format(timex.RFC3339Date))
		assert.Equal(t, "12345-02-04", d2.Format(timex.RFC3339Date))
	})
}

func TestParseDateErrors(t *testing.T) {
	tests := []struct {
		layout    string
		value     string
		errString string
	}{
		{timex.RFC3339Date, "22-10-25", `parsing "22-10-25" as "YYYY-MM-DD": cannot parse "22-10-25" as "YYYY"`},
		{timex.RFC3339Date, "12022-10-25", `parsing "12022-10-25" as "YYYY-MM-DD": cannot parse "2-10-25" as "-"`},
		{" YYYY-MM-DD", "2010-02-04", `parsing "2010-02-04" as " YYYY-MM-DD": cannot parse "2010-02-04" as " "`},
		{" YYYY-MM-DD", "", `parsing "" as " YYYY-MM-DD": cannot parse "" as "YYYY"`},
		{"YY-MM-DD", "a2-10-25", `parsing "a2-10-25" as "YY-MM-DD": cannot parse "a2-10-25" as "YY"`},
		{"YY-M-DD", "22-a0-25", `parsing "22-a0-25" as "YY-M-DD": cannot parse "a0-25" as "M"`},
		{"D MMM YY", "4 --- 00", `parsing "4 --- 00" as "D MMM YY": cannot parse "--- 00" as "MMM"`},
		{"D MMMM YY", "4 --- 00", `parsing "4 --- 00" as "D MMMM YY": cannot parse "--- 00" as "MMMM"`},
	}

	for _, tt := range tests {
		_, err := timex.ParseDate(tt.layout, tt.value)
		assert.EqualError(t, err, tt.errString)
	}
}

func FuzzParseDate(f *testing.F) {
	f.Add("YYYY-MM-DD", "2006-01-02")
	f.Add(" YYYY-MM-DD", "")
	f.Fuzz(func(t *testing.T, layout, value string) {
		assert.NotPanics(t, func() {
			_, _ = timex.ParseDate(layout, value)
		})
	})
}

func FuzzDateFormat(f *testing.F) {
	f.Add(0, "YYYY-MM-DD")
	f.Add(0, " YYYY-MM-DD")
	f.Fuzz(func(t *testing.T, n int, layout string) {
		assert.NotPanics(t, func() {
			date := timex.Date{}.AddDays(n)
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
		date := timex.MustNewDate(tt.year, tt.month, tt.day)
		assert.Equal(t, tt.str, date.String())
		assert.Equal(t, tt.goStr, date.GoString())
	}
}

func TestDateMarshalJSON(t *testing.T) {
	tests := []struct {
		year, month, day int
	}{
		{2006, 1, 2},
		{1, 12, 11},
		{0, 1, 2},
	}

	for _, tt := range tests {
		d1 := timex.MustNewDate(tt.year, tt.month, tt.day)
		bytes, err := d1.MarshalJSON()
		assert.NoError(t, err)

		var d2 timex.Date
		err = d2.UnmarshalJSON(bytes)
		assert.NoError(t, err)
		assert.Equal(t, d1, d2)
	}

	t.Run("Null", func(t *testing.T) {
		var date timex.Date
		err := date.UnmarshalJSON([]byte("null"))
		assert.NoError(t, err)
		assert.True(t, date.IsZero())
	})
}

func TestStringDateMarshalJSON(t *testing.T) {
	tests := []struct {
		year, month, day int
	}{
		{2006, 1, 2},
		{1, 12, 11},
		{0, 1, 2},
	}

	for _, tt := range tests {
		d1 := timex.StringDate{Date: timex.MustNewDate(tt.year, tt.month, tt.day)}
		bytes, err := d1.MarshalJSON()
		assert.NoError(t, err)

		var d2 timex.StringDate
		err = d2.UnmarshalJSON(bytes)
		assert.NoError(t, err)
		assert.Equal(t, d1, d2)
	}

	t.Run("Null", func(t *testing.T) {
		var date timex.StringDate
		err := date.UnmarshalJSON([]byte("null"))
		assert.NoError(t, err)
		assert.True(t, date.Date.IsZero())
	})

	t.Run("EmptyString", func(t *testing.T) {
		var d1 timex.StringDate
		bytes, err := d1.MarshalJSON()
		assert.NoError(t, err)
		assert.Equal(t, `""`, string(bytes))

		var d2 timex.StringDate
		err = d2.UnmarshalJSON(bytes)
		assert.NoError(t, err)
		assert.True(t, d2.Date.IsZero())
	})
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
		date := timex.MustNewDate(tt.year, tt.month, tt.day)
		_, err := date.MarshalJSON()
		assert.EqualError(t, err, tt.errString)
	}
}

func TestDateUnmarshalJSONError(t *testing.T) {
	tests := []struct {
		s         string
		errString string
	}{
		{`2006-01-02`, `Date.UnmarshalJSON: input is not a JSON string`},
		{`""`, `parsing "" as "YYYY-MM-DD"`},
		{`"-1-01-02"`, `parsing "-1-01-02" as "YYYY-MM-DD"`},
		{`"10000-01-02"`, `parsing "10000-01-02" as "YYYY-MM-DD"`},
		{`"12345-01-02"`, `parsing "12345-01-02" as "YYYY-MM-DD"`},
		{`"2006+01+02"`, `parsing "2006+01+02" as "YYYY-MM-DD"`},
		{`"YYYY-01-02"`, `parsing "YYYY-01-02" as "YYYY-MM-DD"`},
	}

	for _, tt := range tests {
		var date timex.Date
		err := date.UnmarshalJSON([]byte(tt.s))
		assert.EqualError(t, err, tt.errString)
	}
}

func FuzzDateUnmarshalJSON(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		assert.NotPanics(t, func() {
			var date timex.Date
			_ = date.UnmarshalJSON(data)
		})
	})
}
