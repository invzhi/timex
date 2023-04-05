package timex

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAppendInt(t *testing.T) {
	tests := []struct {
		n     int
		width int
		s     string
	}{
		{-12, 0, "-12"},
		{-12, 1, "-2"},
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
		{12, 1, "2"},
		{12, 2, "12"},
		{12, 3, "012"},
		{12, 4, "0012"},
		{123, 0, "123"},
		{123, 1, "3"},
		{123, 2, "23"},
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
			// MustNewDate(0, 2, 4),    // Zero year
			// MustNewDate(-123, 2, 4), // Negative year
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
}

func TestParseDateErrors(t *testing.T) {
	tests := []struct {
		layout    string
		value     string
		errString string
	}{
		{RFC3339, "22-10-25", `parsing date "22-10-25" as "YYYY-MM-DD": cannot parse "22-10-25" as "YYYY"`},
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
