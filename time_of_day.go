package timex

import (
	"errors"
	"time"
)

// TimeOfDay represents a specific time in a day.
type TimeOfDay struct {
	n int64
}

// NewTimeOfDay returns the time of day corresponding to hour, minute, second, and nanosecond.
func NewTimeOfDay(hour, min, sec, nsec int) (TimeOfDay, error) {
	if hour < 0 || hour > 23 {
		return TimeOfDay{}, errors.New("hour is out of range [0,23]")
	}
	if min < 0 || min > 59 {
		return TimeOfDay{}, errors.New("minute is out of range [0,59]")
	}
	if sec < 0 || sec > 59 {
		return TimeOfDay{}, errors.New("second is out of range [0,59]")
	}
	if nsec < 0 || nsec >= 1e9 {
		return TimeOfDay{}, errors.New("nanosecond is out of range [0,1e9)")
	}

	return TimeOfDay{n: timeToNanoseconds(hour, min, sec, nsec)}, nil
}

// MustNewTimeOfDay is like NewTimeOfDay but panics if the time of day cannot be created.
func MustNewTimeOfDay(hour, min, sec, nsec int) TimeOfDay {
	timeOfDay, err := NewTimeOfDay(hour, min, sec, nsec)
	if err != nil {
		panic(`timex: NewTimeOfDay: ` + err.Error())
	}
	return timeOfDay
}

// TimeOfDayFromTime returns the time of day specified by t.
func TimeOfDayFromTime(t time.Time) TimeOfDay {
	hour, min, sec := t.Clock()
	nsec := t.Nanosecond()
	return TimeOfDay{n: timeToNanoseconds(hour, min, sec, nsec)}
}

// TimeOfDayNow returns the current time of day in the given location.
func TimeOfDayNow(location *time.Location) TimeOfDay {
	t := time.Now().In(location)
	return TimeOfDayFromTime(t)
}

// Clock returns the hour, minute, second and nanosecond specified by t.
func (t TimeOfDay) Clock() (hour, min, sec, nsec int) {
	return nanosecondsToTime(t.n)
}

// Hour returns the hour specified by t.
func (t TimeOfDay) Hour() int {
	return int(t.n / nsecsEveryHour)
}

// Minute returns the minute specified by t.
func (t TimeOfDay) Minute() int {
	return int(t.n%nsecsEveryHour) / nsecsEveryMinute
}

// Second returns the second specified by t.
func (t TimeOfDay) Second() int {
	return int(t.n%nsecsEveryMinute) / nsecsEverySecond
}

// Nanosecond returns the nanosecond specified by t.
func (t TimeOfDay) Nanosecond() int {
	return int(t.n % nsecsEverySecond)
}

// norm0 normalize the hi and lo into [0, base].
func norm0(hi, lo, base int64) (int64, int64) {
	if lo < 0 {
		n := (-lo-1)/base + 1
		lo += n * base
		hi -= n
	}
	if lo >= base {
		n := lo / base
		lo -= n * base
		hi += n
	}
	return hi, lo
}

// Add returns the exceeded days and the time of day corresponding to adding the given number of hours, minutes, seconds and nanoseconds to t.
func (t TimeOfDay) Add(hours, mins, secs, nsecs int) (int, TimeOfDay) {
	n := t.n +
		int64(hours)*nsecsEveryHour +
		int64(mins)*nsecsEveryMinute +
		int64(secs)*nsecsEverySecond +
		int64(nsecs)

	day, n := norm0(0, n, nsecsEveryDay)
	return int(day), TimeOfDay{n: n}
}

// AddDuration returns the exceeded days and the time of day corresponding to adding the given time.Duration to t.
func (t TimeOfDay) AddDuration(d time.Duration) (int, TimeOfDay) {
	n := t.n + d.Nanoseconds()
	day, n := norm0(0, n, nsecsEveryDay)
	return int(day), TimeOfDay{n: n}
}

// Sub returns the duration t-tt.
func (t TimeOfDay) Sub(tt TimeOfDay) time.Duration {
	return time.Duration(t.n - tt.n)
}

// IsZero reports whether the time of day t is the zero value, 00:00:00.
func (t TimeOfDay) IsZero() bool {
	return t.n == 0
}

// Before reports whether the time of day t is before tt.
func (t TimeOfDay) Before(tt TimeOfDay) bool {
	return t.n < tt.n
}

// After reports whether the time of day t is after tt.
func (t TimeOfDay) After(tt TimeOfDay) bool {
	return t.n > tt.n
}

// Equal reports whether the time of day t and tt is the same time.
func (t TimeOfDay) Equal(tt TimeOfDay) bool {
	return t.n == tt.n
}

const (
	nsecsEverySecond = 1e9
	nsecsEveryMinute = 60 * nsecsEverySecond
	nsecsEveryHour   = 60 * nsecsEveryMinute
	nsecsEveryDay    = 24 * nsecsEveryHour
)

func timeToNanoseconds(hour, min, sec, nsec int) int64 {
	return int64(hour)*nsecsEveryHour +
		int64(min)*nsecsEveryMinute +
		int64(sec)*nsecsEverySecond +
		int64(nsec)
}

func nanosecondsToTime(nsecs int64) (int, int, int, int) {
	hour := nsecs / nsecsEveryHour
	nsecs -= hour * nsecsEveryHour
	min := nsecs / nsecsEveryMinute
	nsecs -= min * nsecsEveryMinute
	sec := nsecs / nsecsEverySecond
	nsecs -= sec * nsecsEverySecond
	return int(hour), int(min), int(sec), int(nsecs)
}
