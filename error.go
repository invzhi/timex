package timex

import "fmt"

// ParseError describes a problem parsing a string.
type ParseError struct {
	Layout     string
	Value      string
	LayoutElem string
	ValueElem  string
}

// Error returns the string representation of a ParseError.
func (e *ParseError) Error() string {
	if len(e.LayoutElem) == 0 && len(e.ValueElem) == 0 {
		return fmt.Sprintf("parsing %q as %q", e.Value, e.Layout)
	}
	return fmt.Sprintf("parsing %q as %q: cannot parse %q as %q", e.Value, e.Layout, e.ValueElem, e.LayoutElem)
}
