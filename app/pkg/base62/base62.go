package base62

import (
	"errors"
	"math"
	"strings"
)

const (
	encode    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	encodeLen = uint64(len(encode))
)

func Encode(n uint64) string {
	if n == 0 {
		return string(encode[0])
	}

	s := int(math.Ceil(math.Log10(float64(n+1)) / math.Log10(float64(encodeLen))))
	r := make([]rune, s)

	for ; n > 0; n /= encodeLen {
		r[s-1] = rune(encode[n%encodeLen])
		s--
	}

	return string(r)
}

func Decode(encoded string) (uint64, error) {
	var number uint64

	encodedLen := len(encoded)

	for i := encodedLen - 1; i >= 0; i-- {
		symbol := rune(encoded[i])
		alphabeticPosition := strings.IndexRune(encode, symbol)

		if alphabeticPosition == -1 {
			return uint64(alphabeticPosition), errors.New("invalid character: " + string(symbol))
		}

		number += uint64(alphabeticPosition) * uint64(math.Pow(float64(encodeLen), float64(encodedLen-i-1)))
	}

	return number, nil
}
