package timex

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseDate(t *testing.T) {
	tests := []struct {
		layout string
		value  string
	}{
		{"YYYY-MM-DD", "2010-02-04"},
		{"MMMM DD YYYY", "February 04 2010"},
		{"MMMM D, YYYY", "February 4, 2010"},
		{"MMM DD YYYY", "Feb 04 2010"},
		{"MMM D YYYY", "Feb 4 2010"},
		{"DD MMM YYYY", "04 Feb 2010"},
		{"DD-MMM-YY", "04-Feb-10"},
		// Case-insensitive
		{"MMM D YYYY", "FEB 4 2010"},
		// Chinese
		{"YYYY年M月D日", "2010年2月4日"},
	}

	for _, tt := range tests {
		date, err := ParseDate(tt.layout, tt.value)
		assert.NoError(t, err)

		assert.Equal(t, 2010, date.Year())
		assert.Equal(t, 2, date.Month())
		assert.Equal(t, 4, date.Day())
	}
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
