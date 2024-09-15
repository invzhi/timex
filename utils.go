package timex

func isDigit(c byte) bool { return c >= '0' && c <= '9' }

func fromDigit(c byte) int { return int(c - '0') }

func toDigit(n int) byte { return byte(n) + '0' }

// match reports whether s1 and s2 match ignoring case.
// It is assumed s1 and s2 are the same length.
func match(s1, s2 string) bool {
	for i := 0; i < len(s1); i++ {
		c1 := s1[i]
		c2 := s2[i]
		if c1 != c2 {
			// Switch to lower-case.
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
	// Optimization for the most common scenario.
	if min == 2 && max == 2 && len(s) >= 2 {
		if !isDigit(s[0]) || !isDigit(s[1]) {
			goto SLOW
		}

		n := fromDigit(s[0])*1e1 + fromDigit(s[1])
		return n, s[2:], true
	}
	if min == 4 && max == 4 && len(s) >= 4 {
		if !isDigit(s[0]) || !isDigit(s[1]) || !isDigit(s[2]) || !isDigit(s[3]) {
			goto SLOW
		}

		n := fromDigit(s[0])*1e3 + fromDigit(s[1])*1e2 + fromDigit(s[2])*1e1 + fromDigit(s[3])
		return n, s[4:], true
	}

SLOW:
	var negative bool
	if len(s) > 0 && (s[0] == '-' || s[0] == '+') {
		negative = s[0] == '-'
		s = s[1:]
	}

	var n, index int
	for ; index < len(s) && index < max; index++ {
		c := s[index]
		if !isDigit(c) {
			break
		}

		n = n*10 + fromDigit(c)
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

	// Optimization for the most common scenario.
	switch {
	case min == 2 && n < 1e2:
		return append(b, toDigit(n/1e1), toDigit(n%1e1))
	case min == 4 && n < 1e4:
		return append(b, toDigit(n/1e3), toDigit(n/1e2%1e1), toDigit(n/1e1%1e1), toDigit(n%1e1))
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
		b[index] = toDigit(n - next*10)
		n = next
	}

	return b
}
