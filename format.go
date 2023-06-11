package timex

import (
	"errors"
	"fmt"
)

const (
	tokenYearTwoDigit = iota + 1
	tokenYearFourDigit
	tokenMonth
	tokenMonthTwoDigit
	tokenMonthShortName
	tokenMonthLongName
	tokenDayOfMonth
	tokenDayOfMonthTwoDigit
)

const (
	RFC3339 = "YYYY-MM-DD"
)

var monthShortNames = []string{
	"Jan",
	"Feb",
	"Mar",
	"Apr",
	"May",
	"Jun",
	"Jul",
	"Aug",
	"Sep",
	"Oct",
	"Nov",
	"Dec",
}

var monthLongNames = []string{
	"January",
	"February",
	"March",
	"April",
	"May",
	"June",
	"July",
	"August",
	"September",
	"October",
	"November",
	"December",
}

// ParseError describes a problem parsing a date string.
type ParseError struct {
	Layout     string
	Value      string
	LayoutElem string
	ValueElem  string
}

// Error returns the string representation of a ParseError.
func (e *ParseError) Error() string {
	if len(e.LayoutElem) == 0 && len(e.ValueElem) == 0 {
		return fmt.Sprintf("parsing date %q as %q", e.Value, e.Layout)
	}
	return fmt.Sprintf("parsing date %q as %q: cannot parse %q as %q", e.Value, e.Layout, e.ValueElem, e.LayoutElem)
}

// match reports whether s1 and s2 match ignoring case.
// It is assumed s1 and s2 are the same length.
func match(s1, s2 string) bool {
	for i := 0; i < len(s1); i++ {
		c1 := s1[i]
		c2 := s2[i]
		if c1 != c2 {
			// switch to lower-case
			c1 |= 'a' - 'A'
			c2 |= 'a' - 'A'
			if c1 != c2 || c1 < 'a' || c1 > 'z' {
				return false
			}
		}
	}
	return true
}

// searchName reports whether the prefix of value exist in names.
// It returns the index and left string.
func searchName(names []string, value string) (int, string, bool) {
	for i, name := range names {
		if len(value) >= len(name) && match(value[:len(name)], name) {
			return i, value[len(name):], true
		}
	}
	return -1, value, false
}

// atoi converts a string to integer with minimum and maximum digit length.
func atoi(s string, min, max int) (int, string, bool) {
	var negative bool
	if len(s) > 0 && (s[0] == '-' || s[0] == '+') {
		negative = s[0] == '-'
		s = s[1:]
	}

	var n, index int
	for ; index < len(s) && index < max; index++ {
		c := s[index]
		if c < '0' || c > '9' {
			break
		}

		n = n*10 + int(c-'0')
	}
	if index < min {
		return 0, "", false
	}

	if negative {
		n = -n
	}

	return n, s[index:], true
}

// appendInt appends the decimal form of integer with specified minimum digit length.
// If specified digit width is zero, the original form of integer will be followed.
func appendInt(b []byte, n int, min int) []byte {
	if n < 0 {
		b = append(b, '-')
		n = -n
	}

	var width int
	if n == 0 {
		width = 1
	}
	for i := n; i > 0; i /= 10 {
		width++
	}

	if min > width {
		width = min
	}

	if len(b)+width <= cap(b) {
		b = b[:len(b)+width]
	} else {
		b = append(b, make([]byte, width)...)
	}

	for i := 0; i < width; i++ {
		index := len(b) - 1 - i

		next := n / 10
		b[index] = byte(n-next*10) + '0'
		n = next
	}

	return b
}

func nextToken(layout string) (prefix string, token int, suffix string) {
	for i := 0; i < len(layout); i++ {
		switch layout[i] {
		case 'Y': // YY, YYYY
			if len(layout) >= i+4 && layout[i:i+4] == "YYYY" {
				return layout[:i], tokenYearFourDigit, layout[i+4:]
			}
			if len(layout) >= i+2 && layout[i:i+2] == "YY" {
				return layout[:i], tokenYearTwoDigit, layout[i+2:]
			}
		case 'M': // M, MM, MMM, MMMM
			if len(layout) >= i+4 && layout[i:i+4] == "MMMM" {
				return layout[:i], tokenMonthLongName, layout[i+4:]
			}
			if len(layout) >= i+3 && layout[i:i+3] == "MMM" {
				return layout[:i], tokenMonthShortName, layout[i+3:]
			}
			if len(layout) >= i+2 && layout[i:i+2] == "MM" {
				return layout[:i], tokenMonthTwoDigit, layout[i+2:]
			}
			if len(layout) >= i+1 && layout[i:i+1] == "M" {
				return layout[:i], tokenMonth, layout[i+1:]
			}
		case 'D': // D, DD
			if len(layout) >= i+2 && layout[i:i+2] == "DD" {
				return layout[:i], tokenDayOfMonthTwoDigit, layout[i+2:]
			}
			if len(layout) >= i+1 && layout[i:i+1] == "D" {
				return layout[:i], tokenDayOfMonth, layout[i+1:]
			}
		}
	}
	return layout, 0, ""
}

func parseStrictRFC3339(b []byte) (Date, error) {
	if len(b) < len(RFC3339) {
		return Date{}, &ParseError{Layout: RFC3339, Value: string(b)}
	}

	ok := true
	parseUint := func(s []byte) (n int) {
		for _, c := range s {
			if c < '0' || c > '9' {
				ok = false
				return 0
			}
			n = n*10 + int(c-'0')
		}
		return n
	}

	year := parseUint(b[0:4])
	month := parseUint(b[5:7])
	day := parseUint(b[8:10])
	if !ok || b[4] != '-' || b[7] != '-' {
		return Date{}, &ParseError{Layout: RFC3339, Value: string(b)}
	}

	return NewDate(year, month, day)
}

// ParseDate parses a formatted string and returns the date it represents.
//
//	YY       01             Two-digit year
//	YYYY   2001             Four-digit year
//	M      1-12             Month, beginning at 1
//	MM    01-12             Month, 2-digits
//	MMM   Jan-Dec           The abbreviated month name
//	MMMM  January-December  The full month name
//	D      1-31             Day of month
//	DD    01-31             Day of month, 2-digits
func ParseDate(layout, value string) (Date, error) {
	var year, month, day int

	err := &ParseError{Layout: layout, Value: value}
	for {
		prefix, token, suffix := nextToken(layout)
		if token == 0 {
			break
		}

		err.LayoutElem = layout[len(prefix) : len(layout)-len(suffix)]

		layout = suffix
		if len(value) < len(prefix) {
			return Date{}, err
		}
		if value[:len(prefix)] != prefix {
			err.LayoutElem = prefix
			err.ValueElem = value
			return Date{}, err
		}
		value = value[len(prefix):]

		err.ValueElem = value

		var ok bool

		switch token {
		case tokenYearTwoDigit:
			year, value, ok = atoi(value, 2, 2)
			if year >= 69 {
				year += 1900
			} else {
				year += 2000
			}
		case tokenYearFourDigit:
			year, value, ok = atoi(value, 4, 4)
		case tokenMonth:
			month, value, ok = atoi(value, 1, 2)
		case tokenMonthTwoDigit:
			month, value, ok = atoi(value, 2, 2)
		case tokenMonthShortName:
			month, value, ok = searchName(monthShortNames, value)
			month++
		case tokenMonthLongName:
			month, value, ok = searchName(monthLongNames, value)
			month++
		case tokenDayOfMonth:
			day, value, ok = atoi(value, 1, 2)
		case tokenDayOfMonthTwoDigit:
			day, value, ok = atoi(value, 2, 2)
		}

		if !ok {
			return Date{}, err
		}
	}

	return NewDate(year, month, day)
}

func (d Date) appendRFC3339(b []byte) []byte {
	year, month, day := ordinalToCalendar(d.ordinal)

	b = appendInt(b, year, 4)
	b = append(b, '-')
	b = appendInt(b, month, 2)
	b = append(b, '-')
	b = appendInt(b, day, 2)
	return b
}

func (d Date) appendStrictRFC3339(b []byte) ([]byte, error) {
	year, month, day := ordinalToCalendar(d.ordinal)
	if year < 0 || year > 9999 {
		return nil, errors.New("year is out of range [0,9999]")
	}

	b = appendInt(b, year, 4)
	b = append(b, '-')
	b = appendInt(b, month, 2)
	b = append(b, '-')
	b = appendInt(b, day, 2)
	return b, nil
}

func (d Date) format(layout string) string {
	year, month, day := ordinalToCalendar(d.ordinal)
	bytes := make([]byte, 0, len(layout)+10)

	for {
		prefix, token, suffix := nextToken(layout)
		bytes = append(bytes, prefix...)
		if token == 0 {
			break
		}

		layout = suffix

		switch token {
		case tokenYearTwoDigit:
			bytes = appendInt(bytes, year%100, 2)
		case tokenYearFourDigit:
			bytes = appendInt(bytes, year, 4)
		case tokenMonth:
			bytes = appendInt(bytes, month, 0)
		case tokenMonthTwoDigit:
			bytes = appendInt(bytes, month, 2)
		case tokenMonthShortName:
			bytes = append(bytes, monthShortNames[month-1]...)
		case tokenMonthLongName:
			bytes = append(bytes, monthLongNames[month-1]...)
		case tokenDayOfMonth:
			bytes = appendInt(bytes, day, 0)
		case tokenDayOfMonthTwoDigit:
			bytes = appendInt(bytes, day, 2)
		}
	}

	return string(bytes)
}

// Format returns a textual representation of the date.
//
//	YY       01             Two-digit year
//	YYYY   2001             Four-digit year
//	M      1-12             Month, beginning at 1
//	MM    01-12             Month, 2-digits
//	MMM   Jan-Dec           The abbreviated month name
//	MMMM  January-December  The full month name
//	D      1-31             Day of month
//	DD    01-31             Day of month, 2-digits
func (d Date) Format(layout string) string {
	switch layout {
	case RFC3339:
		b := make([]byte, 0, len(RFC3339))
		b = d.appendRFC3339(b)
		return string(b)
	default:
		return d.format(layout)
	}
}

// String returns the textual representation of the date.
func (d Date) String() string {
	return d.Format(RFC3339)
}

// GoString returns the Go syntax of the date.
func (d Date) GoString() string {
	year, month, day := ordinalToCalendar(d.ordinal)

	bytes := make([]byte, 0, 32)

	bytes = append(bytes, "timex.MustNewDate("...)
	bytes = appendInt(bytes, year, 0)

	bytes = append(bytes, ", "...)
	bytes = appendInt(bytes, month, 0)

	bytes = append(bytes, ", "...)
	bytes = appendInt(bytes, day, 0)
	bytes = append(bytes, ')')

	return string(bytes)
}

// MarshalJSON implements the json.Marshaler interface.
// The date is a quoted string in RFC 3339 format.
func (d Date) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(RFC3339)+2)
	b = append(b, '"')
	b, err := d.appendStrictRFC3339(b)
	if err != nil {
		return nil, err
	}
	b = append(b, '"')
	return b, nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// The date is expected to be a quoted string in RFC 3339 format.
func (d *Date) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	if len(data) < 2 || data[0] != '"' || data[len(data)-1] != '"' {
		return errors.New("Date.UnmarshalJSON: input is not a JSON string")
	}

	var err error
	*d, err = parseStrictRFC3339(data[1 : len(data)-1])
	return err
}
