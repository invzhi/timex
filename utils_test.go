package timex

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAppendInt(t *testing.T) {
	tests := []struct {
		n     int
		width int
		s     string
	}{
		{-12, 0, "-12"},
		{-12, 1, "-12"},
		{-12, 2, "-12"},
		{-12, 3, "-012"},
		{-12, 4, "-0012"},
		{-1, 0, "-1"},
		{-1, 1, "-1"},
		{-1, 2, "-01"},
		{-1, 3, "-001"},
		{0, 0, "0"},
		{0, 1, "0"},
		{0, 2, "00"},
		{0, 3, "000"},
		{1, 0, "1"},
		{1, 1, "1"},
		{1, 2, "01"},
		{1, 3, "001"},
		{12, 0, "12"},
		{12, 1, "12"},
		{12, 2, "12"},
		{12, 3, "012"},
		{12, 4, "0012"},
		{123, 0, "123"},
		{123, 1, "123"},
		{123, 2, "123"},
		{123, 3, "123"},
		{123, 4, "0123"},
	}

	for _, tt := range tests {
		bytes := appendInt(nil, tt.n, tt.width)
		assert.Equal(t, tt.s, string(bytes))
	}
}
