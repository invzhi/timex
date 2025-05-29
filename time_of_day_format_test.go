package timex_test

import (
	"testing"
	"time"

	"github.com/invzhi/timex"
	"github.com/stretchr/testify/assert"
)

func TestParseTimeOfDay(t *testing.T) {
	tests := []struct {
		layout string
		value  string
	}{
		{timex.RFC3339Time, "15:04:05.000000006"},
		{"HHmmss", "150405.000000006"},
		{"H:m:s", "15:4:5.000000006"},
		{"hh:mm:ssA", "03:04:05.000000006PM"},
		{"hh:mm:ssa", "03:04:05.000000006pm"},
		{"h:m:sa", "3:4:5.000000006pm"},
		// Chinese
		{"HH时m分s秒", "15时4分5.000000006秒"},
	}

	for _, tt := range tests {
		timeOfDay, err := timex.ParseTimeOfDay(tt.layout, tt.value)
		assert.NoError(t, err)

		assert.Equal(t, 15, timeOfDay.Hour())
		assert.Equal(t, 4, timeOfDay.Minute())
		assert.Equal(t, 5, timeOfDay.Second())
		assert.Equal(t, 6, timeOfDay.Nanosecond())
	}

	t.Run("Format", func(t *testing.T) {
		ts := []timex.TimeOfDay{
			timex.MustNewTimeOfDay(0, 0, 0, 0),
			timex.MustNewTimeOfDay(3, 0, 0, 0),
			timex.MustNewTimeOfDay(12, 0, 0, 0),
			timex.MustNewTimeOfDay(15, 4, 5, 0),
			timex.MustNewTimeOfDay(15, 4, 5, 1e8),
			timex.MustNewTimeOfDay(15, 4, 5, 2e7),
			timex.MustNewTimeOfDay(15, 4, 5, 6780),
			timex.MustNewTimeOfDay(23, 59, 59, 1e9-1),
		}

		for _, tt := range tests {
			for _, t1 := range ts {
				t2, err := timex.ParseTimeOfDay(tt.layout, t1.Format(tt.layout))
				assert.NoError(t, err)

				assert.Equal(t, t1, t2)
			}
		}
	})
}

func TestParseTimeOfDayErrors(t *testing.T) {
	tests := []struct {
		layout    string
		value     string
		errString string
	}{
		{timex.RFC3339Time, "3:04:05", `parsing "3:04:05" as "HH:mm:ss": cannot parse "3:04:05" as "HH"`},
		{timex.RFC3339Time, "015:04:05", `parsing "015:04:05" as "HH:mm:ss": cannot parse "5:04:05" as ":"`},
		{" HH:mm:ss", "15:03:04", `parsing "15:03:04" as " HH:mm:ss": cannot parse "15:03:04" as " "`},
		{" HH:mm:ss", "", `parsing "" as " HH:mm:ss": cannot parse "" as "HH"`},
		{"hh:mm:ss", "a3:04:05", `parsing "a3:04:05" as "hh:mm:ss": cannot parse "a3:04:05" as "hh"`},
		{"hh:mm:ss", "03:a4:05", `parsing "03:a4:05" as "hh:mm:ss": cannot parse "a4:05" as "mm"`},
		{"hh:mm:ss", "03:04:a5", `parsing "03:04:a5" as "hh:mm:ss": cannot parse "a5" as "ss"`},
	}

	for _, tt := range tests {
		_, err := timex.ParseTimeOfDay(tt.layout, tt.value)
		assert.EqualError(t, err, tt.errString)
	}
}

func FuzzParseTimeOfDay(f *testing.F) {
	f.Add("HH:mm:ss", "15:04:05")
	f.Add(" HH:mm:ss", "")
	f.Fuzz(func(t *testing.T, layout, value string) {
		assert.NotPanics(t, func() {
			_, _ = timex.ParseTimeOfDay(layout, value)
		})
	})
}

func FuzzTimeOfDayFormat(f *testing.F) {
	f.Add(0, "HH:mm:ss")
	f.Add(0, " HH:mm:ss")
	f.Fuzz(func(t *testing.T, n int, layout string) {
		assert.NotPanics(t, func() {
			_, timeOfDay := timex.TimeOfDay{}.AddDuration(time.Duration(n))
			_ = timeOfDay.Format(layout)
		})
	})
}

func TestTimeOfDayString(t *testing.T) {
	tests := []struct {
		hour, min, sec, nsec int
		str, goStr           string
	}{
		{0, 0, 0, 0, "00:00:00", "timex.MustNewTimeOfDay(0, 0, 0, 0)"},
		{3, 0, 0, 0, "03:00:00", "timex.MustNewTimeOfDay(3, 0, 0, 0)"},
		{12, 0, 0, 0, "12:00:00", "timex.MustNewTimeOfDay(12, 0, 0, 0)"},
		{15, 04, 05, 0, "15:04:05", "timex.MustNewTimeOfDay(15, 4, 5, 0)"},
		{15, 04, 05, 1e8, "15:04:05.1", "timex.MustNewTimeOfDay(15, 4, 5, 100000000)"},
		{15, 04, 05, 2e7, "15:04:05.02", "timex.MustNewTimeOfDay(15, 4, 5, 20000000)"},
		{15, 04, 05, 6780, "15:04:05.00000678", "timex.MustNewTimeOfDay(15, 4, 5, 6780)"},
		{23, 59, 59, 1e9 - 1, "23:59:59.999999999", "timex.MustNewTimeOfDay(23, 59, 59, 999999999)"},
	}

	for _, tt := range tests {
		timeOfDay := timex.MustNewTimeOfDay(tt.hour, tt.min, tt.sec, tt.nsec)
		assert.Equal(t, tt.str, timeOfDay.String())
		assert.Equal(t, tt.goStr, timeOfDay.GoString())
	}
}

func TestTimeOfDayMarshalJSON(t *testing.T) {
	tests := []struct {
		hour, min, sec, nsec int
	}{
		{0, 0, 0, 0},
		{3, 0, 0, 0},
		{12, 0, 0, 0},
		{15, 04, 05, 0},
		{15, 04, 05, 1e8},
		{15, 04, 05, 2e7},
		{15, 04, 05, 6780},
		{23, 59, 59, 1e9 - 1},
	}

	for _, tt := range tests {
		t1 := timex.MustNewTimeOfDay(tt.hour, tt.min, tt.sec, tt.nsec)
		bytes, err := t1.MarshalJSON()
		assert.NoError(t, err)

		var t2 timex.TimeOfDay
		err = t2.UnmarshalJSON(bytes)
		assert.NoError(t, err)
		assert.Equal(t, t1, t2)
	}

	t.Run("Null", func(t *testing.T) {
		var timeOfDay timex.TimeOfDay
		err := timeOfDay.UnmarshalJSON([]byte("null"))
		assert.NoError(t, err)
		assert.True(t, timeOfDay.IsZero())
	})
}

func TestTimeOfDayUnmarshalJSONError(t *testing.T) {
	tests := []struct {
		s         string
		errString string
	}{
		{`15:04:05`, `TimeOfDay.UnmarshalJSON: input is not a JSON string`},
		{`""`, `parsing "" as "HH:mm:ss"`},
		{`"3:04:05"`, `parsing "3:04:05" as "HH:mm:ss"`},
		{`"-3:04:05"`, `parsing "-3:04:05" as "HH:mm:ss"`},
		{`"15:04:5"`, `parsing "15:04:5" as "HH:mm:ss"`},
		{`"100:04:05"`, `parsing "100:04:05" as "HH:mm:ss"`},
		{`"12345:04:05"`, `parsing "12345:04:05" as "HH:mm:ss"`},
		{`"15.04.05"`, `parsing "15.04.05" as "HH:mm:ss"`},
		{`"15-04-05"`, `parsing "15-04-05" as "HH:mm:ss"`},
		{`"HH.04.05"`, `parsing "HH.04.05" as "HH:mm:ss"`},
	}

	for _, tt := range tests {
		var timeOfDay timex.TimeOfDay
		err := timeOfDay.UnmarshalJSON([]byte(tt.s))
		assert.EqualError(t, err, tt.errString)
	}
}

func FuzzTimeOfDayUnmarshalJSON(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		assert.NotPanics(t, func() {
			var timeOfDay timex.TimeOfDay
			_ = timeOfDay.UnmarshalJSON(data)
		})
	})
}
