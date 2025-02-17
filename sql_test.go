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

		assert.Equal(t, tt.s, date.Format(timex.RFC3339Date))
	}

	t.Run("NullDate", func(t *testing.T) {
		var date timex.NullDate
		err := date.Scan(nil)
		assert.NoError(t, err)
		assert.False(t, date.Valid)

		for _, tt := range tests {
			err = date.Scan(tt.value)
			assert.NoError(t, err)

			assert.Equal(t, tt.s, date.Date.Format(timex.RFC3339Date))
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

func TestTimeOfDayScan(t *testing.T) {
	tests := []struct {
		value interface{}
		s     string
	}{
		{[]byte("15:04:05"), "15:04:05"},
		{"15:04:05", "15:04:05"},
		{"15:04:05.123456789", "15:04:05.123456789"},
		{time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC), "15:04:05"},
	}

	for _, tt := range tests {
		var timeOfDay timex.TimeOfDay
		err := timeOfDay.Scan(tt.value)
		assert.NoError(t, err)

		assert.Equal(t, tt.s, timeOfDay.Format(timex.RFC3339Time))
	}

	t.Run("NullTimeOfDay", func(t *testing.T) {
		var timeOfDay timex.NullTimeOfDay
		err := timeOfDay.Scan(nil)
		assert.NoError(t, err)
		assert.False(t, timeOfDay.Valid)

		for _, tt := range tests {
			err = timeOfDay.Scan(tt.value)
			assert.NoError(t, err)

			assert.Equal(t, tt.s, timeOfDay.TimeOfDay.Format(timex.RFC3339Time))
			assert.True(t, timeOfDay.Valid)
		}
	})
}

func TestTimeOfDayScanErrors(t *testing.T) {
	assert.EqualError(t, new(timex.TimeOfDay).Scan(nil), "unsupported type <nil>")
	assert.EqualError(t, new(timex.TimeOfDay).Scan(uint64(1)), "unsupported type uint64")

	t.Run("NullTimeOfDay", func(t *testing.T) {
		assert.EqualError(t, new(timex.NullTimeOfDay).Scan(uint64(1)), "unsupported type uint64")
	})
}

func TestTimeOfDayValue(t *testing.T) {
	tests := []struct {
		timeOfDay timex.TimeOfDay
		value     interface{}
	}{
		{timex.TimeOfDay{}, "00:00:00"},
		{timex.MustNewTimeOfDay(3, 0, 0, 0), "03:00:00"},
		{timex.MustNewTimeOfDay(12, 0, 0, 0), "12:00:00"},
		{timex.MustNewTimeOfDay(15, 4, 5, 0), "15:04:05"},
		{timex.MustNewTimeOfDay(15, 4, 5, 1e8), "15:04:05.1"},
		{timex.MustNewTimeOfDay(15, 4, 5, 2e7), "15:04:05.02"},
		{timex.MustNewTimeOfDay(15, 4, 5, 6780), "15:04:05.00000678"},
		{timex.MustNewTimeOfDay(23, 59, 59, 1e9-1), "23:59:59.999999999"},
	}

	for _, tt := range tests {
		value, err := tt.timeOfDay.Value()
		assert.NoError(t, err)
		assert.Equal(t, tt.value, value)
	}

	t.Run("NullTimeOfDay", func(t *testing.T) {
		value, err := timex.NullTimeOfDay{}.Value()
		assert.NoError(t, err)
		assert.Nil(t, value)

		for _, tt := range tests {
			timeOfDay := timex.NullTimeOfDay{TimeOfDay: tt.timeOfDay, Valid: true}
			value, err := timeOfDay.Value()
			assert.NoError(t, err)
			assert.Equal(t, tt.value, value)
		}
	})
}
