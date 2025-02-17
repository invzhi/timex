package timex

import "errors"

const (
	tokenMidday = iota + 1
	tokenMiddayUppercase
	token24Hour
	token24HourTwoDigit
	token12Hour
	token12HourTwoDigit
	tokenMinute
	tokenMinuteTwoDigit
	tokenSecond
	tokenSecondTwoDigit
)

const (
	RFC3339Time = "HH:mm:ss"
)

func nextTimeToken(layout string) (prefix string, token int, suffix string) {
	for i := 0; i < len(layout); i++ {
		switch layout[i] {
		case 'a':
			return layout[:i], tokenMidday, layout[i+1:]
		case 'A':
			return layout[:i], tokenMiddayUppercase, layout[i+1:]
		case 'H': // H, HH
			if len(layout) >= i+2 && layout[i:i+2] == "HH" {
				return layout[:i], token24HourTwoDigit, layout[i+2:]
			}
			if len(layout) >= i+1 && layout[i:i+1] == "H" {
				return layout[:i], token24Hour, layout[i+1:]
			}
		case 'h': // h, hh
			if len(layout) >= i+2 && layout[i:i+2] == "hh" {
				return layout[:i], token12HourTwoDigit, layout[i+2:]
			}
			if len(layout) >= i+1 && layout[i:i+1] == "h" {
				return layout[:i], token12Hour, layout[i+1:]
			}
		case 'm': // m, mm
			if len(layout) >= i+2 && layout[i:i+2] == "mm" {
				return layout[:i], tokenMinuteTwoDigit, layout[i+2:]
			}
			if len(layout) >= i+1 && layout[i:i+1] == "m" {
				return layout[:i], tokenMinute, layout[i+1:]
			}
		case 's': // s, ss
			if len(layout) >= i+2 && layout[i:i+2] == "ss" {
				return layout[:i], tokenSecondTwoDigit, layout[i+2:]
			}
			if len(layout) >= i+1 && layout[i:i+1] == "s" {
				return layout[:i], tokenSecond, layout[i+1:]
			}
		}
	}
	return layout, 0, ""
}

func parseStrictRFC3339Time(b []byte) (TimeOfDay, error) {
	ok := true
	parseUint := func(s []byte) (n int) {
		for _, c := range s {
			if !isDigit(c) {
				ok = false
				return 0
			}
			n = n*10 + int(c-'0')
		}
		return n
	}

	hour := parseUint(b[0:2])
	min := parseUint(b[3:5])
	if !ok || b[2] != ':' || b[5] != ':' {
		return TimeOfDay{}, &ParseError{Layout: RFC3339Time, Value: string(b)}
	}
	sec, nsec, _, ok := atof(string(b[6:]), 2, 2, 9)
	if !ok {
		return TimeOfDay{}, &ParseError{Layout: RFC3339Time, Value: string(b)}
	}

	return NewTimeOfDay(hour, min, sec, nsec)
}

// ParseTimeOfDay parses a formatted string and returns the time of day it represents.
//
//	a   am/pm  ante meridiem or post meridiem
//	A   AM/PM  ante meridiem or post meridiem
//	H    0-23  Two-digit hour, 24-hour clock
//	HH  00-23  Hour, 24-hour clock
//	h    1-12  Two-digit hour, 12-hour clock
//	hh  01-12  Hour, 12-hour clock
//	m    0-59  Minute
//	mm  00-59  Minute, 2-digits
//	s    0-59  Second, including fraction
//	ss  00-59  Second, 2-digits, including fraction
func ParseTimeOfDay(layout, value string) (TimeOfDay, error) {
	var hour, min, sec, nsec int
	var amSet, pmSet bool

	originLayout, originValue := layout, value
	var layoutElem, valueElem string
	for {
		prefix, token, suffix := nextTimeToken(layout)
		if token == 0 {
			break
		}

		layoutElem = layout[len(prefix) : len(layout)-len(suffix)]

		layout = suffix
		if len(value) < len(prefix) {
			return TimeOfDay{}, &ParseError{Layout: originLayout, Value: originValue, LayoutElem: layoutElem, ValueElem: valueElem}
		}
		if value[:len(prefix)] != prefix {
			return TimeOfDay{}, &ParseError{Layout: originLayout, Value: originValue, LayoutElem: prefix, ValueElem: value}
		}
		value = value[len(prefix):]

		valueElem = value

		var ok bool

		switch token {
		case tokenMidday:
			var index int
			index, value, ok = searchName([]string{"am", "pm"}, value)
			switch index {
			case 0:
				amSet = true
			case 1:
				pmSet = true
			}
		case tokenMiddayUppercase:
			var index int
			index, value, ok = searchName([]string{"AM", "PM"}, value)
			switch index {
			case 0:
				amSet = true
			case 1:
				pmSet = true
			}
		case token24Hour, token12Hour:
			hour, value, ok = atoi(value, 1, 2)
		case token24HourTwoDigit, token12HourTwoDigit:
			hour, value, ok = atoi(value, 2, 2)
		case tokenMinute:
			min, value, ok = atoi(value, 1, 2)
		case tokenMinuteTwoDigit:
			min, value, ok = atoi(value, 2, 2)
		case tokenSecond:
			sec, nsec, value, ok = atof(value, 1, 2, 9)
		case tokenSecondTwoDigit:
			sec, nsec, value, ok = atof(value, 2, 2, 9)
		}

		if !ok {
			return TimeOfDay{}, &ParseError{Layout: originLayout, Value: originValue, LayoutElem: layoutElem, ValueElem: valueElem}
		}
	}

	if amSet && hour == 12 {
		hour = 0
	} else if pmSet && hour < 12 {
		hour += 12
	}

	return NewTimeOfDay(hour, min, sec, nsec)
}

func (t TimeOfDay) appendRFC3339(b []byte) []byte {
	hour, min, sec, nsec := nanosecondsToTime(t.n)

	b = appendInt(b, hour, 2)
	b = append(b, ':')
	b = appendInt(b, min, 2)
	b = append(b, ':')
	b = appendInt(b, sec, 2)
	b = appendFraction(b, nsec, 9)
	return b
}

// hour12 returns the hour of 12-hour clock from 24-hour clock.
func hour12(hour int) int {
	hour12 := hour % 12
	if hour12 == 0 {
		hour12 = 12
	}
	return hour12
}

// midday returns the string of am/pm from 24-hour clock.
func midday(hour int, am, pm string) string {
	if hour < 12 {
		return am
	}
	return pm
}

func (t TimeOfDay) format(layout string) string {
	hour, min, sec, nsec := nanosecondsToTime(t.n)
	bytes := make([]byte, 0, len(layout)+10)

	for {
		prefix, token, suffix := nextTimeToken(layout)
		bytes = append(bytes, prefix...)
		if token == 0 {
			break
		}

		layout = suffix

		switch token {
		case tokenMidday:
			bytes = append(bytes, midday(hour, "am", "pm")...)
		case tokenMiddayUppercase:
			bytes = append(bytes, midday(hour, "AM", "PM")...)
		case token24Hour:
			bytes = appendInt(bytes, hour, 0)
		case token24HourTwoDigit:
			bytes = appendInt(bytes, hour, 2)
		case token12Hour:
			bytes = appendInt(bytes, hour12(hour), 0)
		case token12HourTwoDigit:
			bytes = appendInt(bytes, hour12(hour), 2)
		case tokenMinute:
			bytes = appendInt(bytes, min, 0)
		case tokenMinuteTwoDigit:
			bytes = appendInt(bytes, min, 2)
		case tokenSecond:
			bytes = appendInt(bytes, sec, 0)
			bytes = appendFraction(bytes, nsec, 9)
		case tokenSecondTwoDigit:
			bytes = appendInt(bytes, sec, 2)
			bytes = appendFraction(bytes, nsec, 9)
		}
	}

	return string(bytes)
}

// Format returns a textual representation of the time of day.
//
//	a   am/pm  ante meridiem or post meridiem
//	A   AM/PM  ante meridiem or post meridiem
//	H    0-23  Two-digit hour, 24-hour clock
//	HH  00-23  Hour, 24-hour clock
//	h    1-12  Two-digit hour, 12-hour clock
//	hh  01-12  Hour, 12-hour clock
//	m    0-59  Minute
//	mm  00-59  Minute, 2-digits
//	s    0-59  Second
//	ss  00-59  Second, 2-digits
func (t TimeOfDay) Format(layout string) string {
	switch layout {
	case RFC3339Time:
		b := make([]byte, 0, len(RFC3339Time))
		b = t.appendRFC3339(b)
		return string(b)
	default:
		return t.format(layout)
	}
}

// String returns the textual representation of the time of day.
func (t TimeOfDay) String() string {
	return t.Format(RFC3339Time)
}

// GoString returns the Go syntax of the time of day.
func (t TimeOfDay) GoString() string {
	hour, min, sec, nsec := nanosecondsToTime(t.n)

	bytes := make([]byte, 0, 32)

	bytes = append(bytes, "timex.MustNewTimeOfDay("...)
	bytes = appendInt(bytes, hour, 0)

	bytes = append(bytes, ", "...)
	bytes = appendInt(bytes, min, 0)

	bytes = append(bytes, ", "...)
	bytes = appendInt(bytes, sec, 0)

	bytes = append(bytes, ", "...)
	bytes = appendInt(bytes, nsec, 0)

	bytes = append(bytes, ')')

	return string(bytes)
}

// MarshalJSON implements the json.Marshaler interface.
// The time of day is a quoted string in RFC 3339 format.
func (t TimeOfDay) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(RFC3339Time)+2)
	b = append(b, '"')
	b = t.appendRFC3339(b)
	b = append(b, '"')
	return b, nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// The time of day is expected to be a quoted string in RFC 3339 format.
func (t *TimeOfDay) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	if len(data) < 2 || data[0] != '"' || data[len(data)-1] != '"' {
		return errors.New("TimeOfDay.UnmarshalJSON: input is not a JSON string")
	}

	var err error
	*t, err = parseStrictRFC3339Time(data[1 : len(data)-1])
	return err
}
