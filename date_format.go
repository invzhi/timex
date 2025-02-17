package timex

import "errors"

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
	// Deprecated: Use [RFC3339Date] instead.
	RFC3339 = RFC3339Date

	RFC3339Date = "YYYY-MM-DD"
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

func nextDateToken(layout string) (prefix string, token int, suffix string) {
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

func parseStrictRFC3339Date(b []byte) (Date, error) {
	if len(b) != len(RFC3339Date) {
		return Date{}, &ParseError{Layout: RFC3339Date, Value: string(b)}
	}

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

	year := parseUint(b[0:4])
	month := parseUint(b[5:7])
	day := parseUint(b[8:10])
	if !ok || b[4] != '-' || b[7] != '-' {
		return Date{}, &ParseError{Layout: RFC3339Date, Value: string(b)}
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

	originLayout, originValue := layout, value
	var layoutElem, valueElem string
	for {
		prefix, token, suffix := nextDateToken(layout)
		if token == 0 {
			break
		}

		layoutElem = layout[len(prefix) : len(layout)-len(suffix)]

		layout = suffix
		if len(value) < len(prefix) {
			return Date{}, &ParseError{Layout: originLayout, Value: originValue, LayoutElem: layoutElem, ValueElem: valueElem}
		}
		if value[:len(prefix)] != prefix {
			return Date{}, &ParseError{Layout: originLayout, Value: originValue, LayoutElem: prefix, ValueElem: value}
		}
		value = value[len(prefix):]

		valueElem = value

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
			var index int
			index, value, ok = searchName(monthShortNames, value)
			month = index + 1
		case tokenMonthLongName:
			var index int
			index, value, ok = searchName(monthLongNames, value)
			month = index + 1
		case tokenDayOfMonth:
			day, value, ok = atoi(value, 1, 2)
		case tokenDayOfMonthTwoDigit:
			day, value, ok = atoi(value, 2, 2)
		}

		if !ok {
			return Date{}, &ParseError{Layout: originLayout, Value: originValue, LayoutElem: layoutElem, ValueElem: valueElem}
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
		prefix, token, suffix := nextDateToken(layout)
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
	case RFC3339Date:
		b := make([]byte, 0, len(RFC3339Date))
		b = d.appendRFC3339(b)
		return string(b)
	default:
		return d.format(layout)
	}
}

// String returns the textual representation of the date.
func (d Date) String() string {
	return d.Format(RFC3339Date)
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
	b := make([]byte, 0, len(RFC3339Date)+2)
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
	*d, err = parseStrictRFC3339Date(data[1 : len(data)-1])
	return err
}
