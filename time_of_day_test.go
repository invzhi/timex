package timex_test

import (
	"testing"
	"time"

	"github.com/invzhi/timex"
	"github.com/stretchr/testify/assert"
)

func BenchmarkTimeOfDayClock(b *testing.B) {
	timeOfDay := timex.MustNewTimeOfDay(15, 4, 5, 6)

	b.Run("Timex", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			timeOfDay.Clock()
		}
	})
	b.Run("Time", func(b *testing.B) {
		hour, min, sec, nsec := timeOfDay.Clock()
		t := time.Date(2006, 1, 2, hour, min, sec, nsec, time.UTC)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			t.Clock()
		}
	})
}

func BenchmarkTimeOfDayHour(b *testing.B) {
	timeOfDay := timex.MustNewTimeOfDay(15, 4, 5, 6)

	b.Run("Timex", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			timeOfDay.Hour()
		}
	})
	b.Run("Time", func(b *testing.B) {
		hour, min, sec, nsec := timeOfDay.Clock()
		t := time.Date(2006, 1, 2, hour, min, sec, nsec, time.UTC)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			t.Hour()
		}
	})
}

func BenchmarkTimeOfDayMinute(b *testing.B) {
	timeOfDay := timex.MustNewTimeOfDay(15, 4, 5, 6)

	b.Run("Timex", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			timeOfDay.Minute()
		}
	})
	b.Run("Time", func(b *testing.B) {
		hour, min, sec, nsec := timeOfDay.Clock()
		t := time.Date(2006, 1, 2, hour, min, sec, nsec, time.UTC)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			t.Minute()
		}
	})
}

func BenchmarkTimeOfDaySecond(b *testing.B) {
	timeOfDay := timex.MustNewTimeOfDay(15, 4, 5, 6)

	b.Run("Timex", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			timeOfDay.Second()
		}
	})
	b.Run("Time", func(b *testing.B) {
		hour, min, sec, nsec := timeOfDay.Clock()
		t := time.Date(2006, 1, 2, hour, min, sec, nsec, time.UTC)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			t.Second()
		}
	})
}

func BenchmarkTimeOfDayNanosecond(b *testing.B) {
	timeOfDay := timex.MustNewTimeOfDay(15, 4, 5, 6)

	b.Run("Timex", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			timeOfDay.Nanosecond()
		}
	})
	b.Run("Time", func(b *testing.B) {
		hour, min, sec, nsec := timeOfDay.Clock()
		t := time.Date(2006, 1, 2, hour, min, sec, nsec, time.UTC)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			t.Nanosecond()
		}
	})
}

func BenchmarkTimeOfDayAdd(b *testing.B) {
	timeOfDay := timex.MustNewTimeOfDay(15, 4, 5, 6)

	b.Run("Timex", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			timeOfDay.AddDuration(24 * time.Hour)
		}
	})
	b.Run("Time", func(b *testing.B) {
		hour, min, sec, nsec := timeOfDay.Clock()
		t := time.Date(2006, 1, 2, hour, min, sec, nsec, time.UTC)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			t.Add(24 * time.Hour)
		}
	})
}

func BenchmarkTimeOfDaySub(b *testing.B) {
	t1 := timex.MustNewTimeOfDay(15, 4, 5, 6)
	t2 := timex.MustNewTimeOfDay(14, 4, 5, 6)

	b.Run("Timex", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			t1.Sub(t2)
		}
	})
	b.Run("Time", func(b *testing.B) {
		hour, min, sec, nsec := t1.Clock()
		tt1 := time.Date(2006, 1, 2, hour, min, sec, nsec, time.UTC)
		tt2 := time.Date(2006, 1, 2, hour, min, sec, nsec, time.UTC)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			tt1.Sub(tt2)
		}
	})
}

func TestNewTimeOfDayFromTime(t *testing.T) {
	tests := []struct {
		hour, min, sec, nsec int
	}{
		{0, 0, 0, 0},
		{10, 0, 0, 0},
		{12, 0, 0, 0},
		{15, 04, 05, 0},
		{23, 59, 59, 1e9 - 1},
	}

	for _, tt := range tests {
		datetime := time.Date(2003, 1, 2, tt.hour, tt.min, tt.sec, tt.nsec, time.UTC)
		timeOfDay := timex.TimeOfDayFromTime(datetime)

		hour, min, sec, nsec := timeOfDay.Clock()
		assert.Equal(t, tt.hour, hour)
		assert.Equal(t, tt.min, min)
		assert.Equal(t, tt.sec, sec)
		assert.Equal(t, tt.nsec, nsec)

		assert.Equal(t, tt.hour, timeOfDay.Hour())
		assert.Equal(t, tt.min, timeOfDay.Minute())
		assert.Equal(t, tt.sec, timeOfDay.Second())
		assert.Equal(t, tt.nsec, timeOfDay.Nanosecond())
	}

	t.Run("Now", func(t *testing.T) {
		now := time.Now().In(time.UTC)
		timeOfDay := timex.TimeOfDayNow(time.UTC)

		hour, min, sec, nsec := timeOfDay.Clock()
		assert.LessOrEqual(t, now.Hour(), hour)
		assert.LessOrEqual(t, now.Minute(), min)
		assert.LessOrEqual(t, now.Second(), sec)
		assert.LessOrEqual(t, now.Nanosecond(), nsec)

		assert.LessOrEqual(t, now.Hour(), timeOfDay.Hour())
		assert.LessOrEqual(t, now.Minute(), timeOfDay.Minute())
		assert.LessOrEqual(t, now.Second(), timeOfDay.Second())
		assert.LessOrEqual(t, now.Nanosecond(), timeOfDay.Nanosecond())
	})
}

func TestNewTimeOfDayErrors(t *testing.T) {
	tests := []struct {
		hour, min, sec, nsec int
		errString            string
	}{
		{-3, 0, 0, 0, "hour is out of range [0,23]"},
		{-1, 0, 0, 0, "hour is out of range [0,23]"},
		{24, 0, 0, 0, "hour is out of range [0,23]"},
		{30, 0, 0, 0, "hour is out of range [0,23]"},
		{12, -3, 0, 0, "minute is out of range [0,59]"},
		{12, -1, 0, 0, "minute is out of range [0,59]"},
		{12, 60, 0, 0, "minute is out of range [0,59]"},
		{12, 90, 0, 0, "minute is out of range [0,59]"},
		{12, 0, -3, 0, "second is out of range [0,59]"},
		{12, 0, -1, 0, "second is out of range [0,59]"},
		{12, 0, 60, 0, "second is out of range [0,59]"},
		{12, 0, 90, 0, "second is out of range [0,59]"},
		{12, 0, 0, -3, "nanosecond is out of range [0,1e9)"},
		{12, 0, 0, -1, "nanosecond is out of range [0,1e9)"},
		{12, 0, 0, 1e9, "nanosecond is out of range [0,1e9)"},
		{12, 0, 0, 1e9 + 1, "nanosecond is out of range [0,1e9)"},
	}

	for _, tt := range tests {
		_, err := timex.NewTimeOfDay(tt.hour, tt.min, tt.sec, tt.nsec)
		assert.EqualError(t, err, tt.errString)

		assert.Panicsf(t, func() {
			_ = timex.MustNewTimeOfDay(tt.hour, tt.min, tt.sec, tt.nsec)
		}, "timex: NewTimeOfDay: "+tt.errString)
	}
}

func TestTimeOfDayAdd(t *testing.T) {
	tests := []struct {
		hours, mins, secs, nsecs   int
		days, hour, min, sec, nsec int
	}{
		{
			0, 0, 0, 0,
			0, 0, 0, 0, 0,
		},
		{
			0, 0, 0, -1,
			-1, 23, 59, 59, 1e9 - 1,
		},
		{
			-23, 0, 0, 0,
			-1, 1, 0, 0, 0,
		},
		{
			-24, 0, 0, 0,
			-1, 0, 0, 0, 0,
		},
		{
			-1, 60, 0, 1,
			0, 0, 0, 0, 1,
		},
		{
			0, 0, 0, 1,
			0, 0, 0, 0, 1,
		},
		{
			0, 0, 0, 24 * 60 * 60 * 1e9,
			1, 0, 0, 0, 0,
		},
		{
			24, 0, 0, 0,
			1, 0, 0, 0, 0,
		},
		{
			7 * 24, 0, 0, 0,
			7, 0, 0, 0, 0,
		},
		{
			7 * 24, 61, 61, 0,
			7, 1, 2, 1, 0,
		},
	}

	timeOfDay := timex.MustNewTimeOfDay(0, 0, 0, 0)
	t.Run("Add", func(t *testing.T) {
		for _, tt := range tests {
			days, newTimeOfDay := timeOfDay.Add(tt.hours, tt.mins, tt.secs, tt.nsecs)

			hour, min, sec, nsec := newTimeOfDay.Clock()
			assert.Equal(t, tt.days, days)
			assert.Equal(t, tt.hour, hour)
			assert.Equal(t, tt.min, min)
			assert.Equal(t, tt.sec, sec)
			assert.Equal(t, tt.nsec, nsec)
		}
	})
	t.Run("AddDuration", func(t *testing.T) {
		for _, tt := range tests {
			d := time.Duration(tt.hours)*time.Hour +
				time.Duration(tt.mins)*time.Minute +
				time.Duration(tt.secs)*time.Second +
				time.Duration(tt.nsecs)*time.Nanosecond
			days, newTimeOfDay := timeOfDay.AddDuration(d)

			hour, min, sec, nsec := newTimeOfDay.Clock()
			assert.Equal(t, tt.days, days)
			assert.Equal(t, tt.hour, hour)
			assert.Equal(t, tt.min, min)
			assert.Equal(t, tt.sec, sec)
			assert.Equal(t, tt.nsec, nsec)
		}
	})
}

func TestTimeOfDaySub(t *testing.T) {
	tests := []struct {
		t1, t2 timex.TimeOfDay
		d      time.Duration
	}{
		{timex.TimeOfDay{}, timex.TimeOfDay{}, 0},
		{timex.MustNewTimeOfDay(0, 0, 0, 0), timex.MustNewTimeOfDay(23, 0, 0, 0), -23 * time.Hour},
		{timex.MustNewTimeOfDay(12, 0, 0, 0), timex.MustNewTimeOfDay(0, 0, 0, 0), 12 * time.Hour},
		{timex.MustNewTimeOfDay(12, 0, 0, 0), timex.MustNewTimeOfDay(11, 59, 0, 0), 1 * time.Minute},
		{timex.MustNewTimeOfDay(3, 0, 0, 0), timex.MustNewTimeOfDay(11, 29, 58, 589235), -8*time.Hour - 29*time.Minute - 58*time.Second - 589235},
	}

	for _, tt := range tests {
		d := tt.t1.Sub(tt.t2)
		assert.Equal(t, tt.d, d)
	}
}

func TestTimeOfDayIsZero(t *testing.T) {
	assert.True(t, timex.TimeOfDay{}.IsZero())
	assert.True(t, timex.MustNewTimeOfDay(0, 0, 0, 0).IsZero())
	assert.True(t, timex.TimeOfDayFromTime(time.Time{}).IsZero())
}

func TestTimeOfDayBeforeAfter(t *testing.T) {
	tests := []struct {
		t1, t2        timex.TimeOfDay
		before, after bool
	}{
		{timex.MustNewTimeOfDay(0, 0, 0, 0), timex.MustNewTimeOfDay(0, 0, 0, 0), false, false},

		{timex.MustNewTimeOfDay(0, 0, 0, 0), timex.MustNewTimeOfDay(0, 0, 0, 1), true, false},
		{timex.MustNewTimeOfDay(0, 0, 0, 1), timex.MustNewTimeOfDay(0, 0, 0, 2), true, false},
		{timex.MustNewTimeOfDay(0, 0, 1, 0), timex.MustNewTimeOfDay(0, 0, 1, 1), true, false},
		{timex.MustNewTimeOfDay(0, 0, 1, 1), timex.MustNewTimeOfDay(0, 0, 1, 2), true, false},
		{timex.MustNewTimeOfDay(0, 1, 0, 0), timex.MustNewTimeOfDay(0, 1, 0, 1), true, false},
		{timex.MustNewTimeOfDay(0, 1, 0, 1), timex.MustNewTimeOfDay(0, 1, 0, 2), true, false},
		{timex.MustNewTimeOfDay(1, 0, 0, 0), timex.MustNewTimeOfDay(1, 0, 0, 1), true, false},
		{timex.MustNewTimeOfDay(1, 0, 0, 1), timex.MustNewTimeOfDay(1, 0, 0, 2), true, false},

		{timex.MustNewTimeOfDay(12, 0, 0, 0), timex.MustNewTimeOfDay(11, 59, 59, 0), false, true},
		{timex.MustNewTimeOfDay(0, 0, 0, 0), timex.MustNewTimeOfDay(23, 59, 59, 0), true, false},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.before, tt.t1.Before(tt.t2))
		assert.Equal(t, tt.after, tt.t1.After(tt.t2))
		assert.Equal(t, !tt.before && !tt.after, tt.t1.Equal(tt.t2))
		assert.Equal(t, !tt.before && !tt.after, tt.t2.Equal(tt.t1))
	}
}
