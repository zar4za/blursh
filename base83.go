package blursh

import (
	"math"
	"strings"
)

var (
	digitCharacters = [...]byte{
		'0',
		'1',
		'2',
		'3',
		'4',
		'5',
		'6',
		'7',
		'8',
		'9',
		'A',
		'B',
		'C',
		'D',
		'E',
		'F',
		'G',
		'H',
		'I',
		'J',
		'K',
		'L',
		'M',
		'N',
		'O',
		'P',
		'Q',
		'R',
		'S',
		'T',
		'U',
		'V',
		'W',
		'X',
		'Y',
		'Z',
		'a',
		'b',
		'c',
		'd',
		'e',
		'f',
		'g',
		'h',
		'i',
		'j',
		'k',
		'l',
		'm',
		'n',
		'o',
		'p',
		'q',
		'r',
		's',
		't',
		'u',
		'v',
		'w',
		'x',
		'y',
		'z',
		'#',
		'$',
		'%',
		'*',
		'+',
		',',
		'-',
		'.',
		':',
		';',
		'=',
		'?',
		'@',
		'[',
		']',
		'^',
		'_',
		'{',
		'|',
		'}',
		'~',
	}
)

func Encode83(builder *strings.Builder, n int, length int) {
	// TODO: probably this also needs refactor

	for i := 1; i <= length; i++ {
		digit := int(math.Floor(math.Floor(float64(n)/math.Pow(83, float64(length-i))))) % 83
		builder.WriteByte(digitCharacters[digit])
	}
}

func Decode83(hash string) int {
	// TODO: pls rewrite this misunderstanding

	value := 0

	for _, char := range hash {
		for digit, ch := range digitCharacters {
			if ch == byte(char) {
				return value*83 + digit
			}
		}
	}

	return value
}
