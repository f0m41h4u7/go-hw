package hw02_unpack_string //nolint:golint,stylecheck

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(input string) (string, error) {
	if input == "" {
		return "", nil
	}

	str := []rune(input)
	var lastChar = str[0]

	// Valid string cannot begin with number
	if unicode.IsDigit(lastChar) {
		return "", ErrInvalidString
	}

	var tmp rune                   // Temporary saves symbol when reading it's amount
	var resBuilder strings.Builder // Result string
	var qtyBuilder strings.Builder // Counts quantity of symbols

	for _, currChar := range str[1:] {
		switch {
		case unicode.IsDigit(currChar):
			if !unicode.IsDigit(lastChar) {
				tmp = lastChar
			}
			qtyBuilder.WriteRune(currChar)
			lastChar = currChar
		case currChar == lastChar:
			return "", ErrInvalidString
		case unicode.IsDigit(lastChar):
			qty, _ := strconv.Atoi(qtyBuilder.String())
			resBuilder.WriteString(strings.Repeat(string(tmp), qty))
			lastChar = currChar
			qtyBuilder.Reset()
		default:
			resBuilder.WriteRune(lastChar)
			lastChar = currChar
			qtyBuilder.Reset()
		}
	}

	// Process last symbols
	if unicode.IsDigit(lastChar) {
		qty, _ := strconv.Atoi(qtyBuilder.String())
		resBuilder.WriteString(strings.Repeat(string(tmp), qty))
	} else {
		resBuilder.WriteRune(lastChar)
	}

	return resBuilder.String(), nil
}
