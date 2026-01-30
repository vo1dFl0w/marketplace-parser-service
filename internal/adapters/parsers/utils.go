package parsers

import (
	"strconv"
	"strings"
	"unicode"
)

func ParseStringToFloat64(s string) (float64, error) {
	rs := []rune(s)

	var builder strings.Builder
	for _, r := range rs {
		if unicode.IsDigit(r) || r == '.' || r == ',' {
			if r == ',' {
				r = '.'
			}
			builder.WriteRune(r)
		}
	}

	str := builder.String()
	res, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0.0, err
	}

	return res, nil
}

func ParseStringToInteger(s string) (int, error) {
	rs := []rune(s)

	var builder strings.Builder
	for _, r := range rs {
		if unicode.IsDigit(r) {
			builder.WriteRune(r)
		}
	}

	str := builder.String()
	res, err := strconv.Atoi(str)
	if err != nil {
		return 0, err
	}

	return res, nil
}