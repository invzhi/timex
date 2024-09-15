package timex_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/invzhi/timex"
)

func TestDateScan(t *testing.T) {
	tests := []struct {
		value interface{}
		s     string
	}{
		{[]byte("2006-01-02"), "2006-01-02"},
		{[]byte("1999-12-24"), "1999-12-24"},
		{"2006-01-02", "2006-01-02"},
		{"1999-12-24", "1999-12-24"},
		{"2006-01-02 00:00:00", "2006-01-02"},
		{"1999-12-24 00:00:00", "1999-12-24"},
		{time.Date(2006, 1, 2, 0, 0, 0, 0, time.UTC), "2006-01-02"},
	}

	for _, tt := range tests {
		var date timex.Date
		err := date.Scan(tt.value)
		assert.NoError(t, err)

		assert.Equal(t, tt.s, date.Format(timex.RFC3339))
	}

	t.Run("NullDate", func(t *testing.T) {
		var date timex.NullDate
		err := date.Scan(nil)
		assert.NoError(t, err)
		assert.False(t, date.Valid)

		for _, tt := range tests {
			err = date.Scan(tt.value)
			assert.NoError(t, err)

			assert.Equal(t, tt.s, date.Date.Format(timex.RFC3339))
			assert.True(t, date.Valid)
		}
	})
}

func TestDateScanErrors(t *testing.T) {
	assert.EqualError(t, new(timex.Date).Scan(nil), "unsupported type <nil>")
	assert.EqualError(t, new(timex.Date).Scan(uint64(1)), "unsupported type uint64")

	t.Run("NullDate", func(t *testing.T) {
		assert.EqualError(t, new(timex.NullDate).Scan(uint64(1)), "unsupported type uint64")
	})
}

func TestDateValue(t *testing.T) {
	tests := []struct {
		date  timex.Date
		value interface{}
	}{
		{timex.Date{}, "0001-01-01"},
		{timex.MustNewDate(2006, 1, 2), "2006-01-02"},
		{timex.MustNewDate(1996, 12, 24), "1996-12-24"},
	}

	for _, tt := range tests {
		value, err := tt.date.Value()
		assert.NoError(t, err)
		assert.Equal(t, tt.value, value)
	}

	t.Run("NullDate", func(t *testing.T) {
		value, err := timex.NullDate{}.Value()
		assert.NoError(t, err)
		assert.Nil(t, value)

		for _, tt := range tests {
			date := timex.NullDate{Date: tt.date, Valid: true}
			value, err := date.Value()
			assert.NoError(t, err)
			assert.Equal(t, tt.value, value)
		}
	})
}
