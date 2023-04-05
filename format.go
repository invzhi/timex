package timex

import (
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

// atoi converts a string to integer with min and max digit length.
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

// appendInt appends the decimal form of integer with specified digit width.
// If specified digit width is zero, the original form of integer will be followed.
func appendInt(b []byte, n int, width int) []byte {
	if n < 0 {
		b = append(b, '-')
		n = -n
	}

	if width == 0 {
		if n == 0 {
			width = 1
		}
		for i := n; i > 0; i /= 10 {
			width++
		}
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
			bytes = appendInt(bytes, d.year, 2)
		case tokenYearFourDigit:
			bytes = appendInt(bytes, d.year, 4)
		case tokenMonth:
			bytes = appendInt(bytes, d.month, 0)
		case tokenMonthTwoDigit:
			bytes = appendInt(bytes, d.month, 2)
		case tokenMonthShortName:
			bytes = append(bytes, monthShortNames[d.month-1]...)
		case tokenMonthLongName:
			bytes = append(bytes, monthLongNames[d.month-1]...)
		case tokenDayOfMonth:
			bytes = appendInt(bytes, d.day, 0)
		case tokenDayOfMonthTwoDigit:
			bytes = appendInt(bytes, d.day, 2)
		}
	}

	return string(bytes)
}
