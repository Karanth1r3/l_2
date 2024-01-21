package dev02

import (
	"fmt"
	"strconv"
	"unicode"
	"unicode/utf8"
)

func main() {
	(UnpackString(`ew4r2k`))
	fmt.Println(UnpackString(`ew\\5r\2k`))

	fmt.Println(UnpackString(`ew\\0r2k`))
}

func loopChar(r rune, count int) []rune {
	res := make([]rune, count)
	for i := 0; i < count; i++ {
		res[i] = r
	}
	return res
}

func UnpackString(s string) (string, error) {
	result := make([]rune, 0)

	if s == "" {
		return "", fmt.Errorf("incorrect string")
	}
	var letterCount = 0
	for _, elem := range s {
		if !unicode.IsDigit(elem) {
			letterCount++
		}
	}
	if letterCount == 0 {
		return "", fmt.Errorf("incorrect string")
	}

	ln := utf8.RuneCountInString(s)
	rs := []rune(s)
	if unicode.IsDigit(rs[0]) {
		return "", fmt.Errorf("incorrect string")
	}

	for i := 0; i < ln; i++ {
		if i < ln-1 {
			switch {
			// Should have thought about escapes from the beginning i guess. This looks fucked up now
			case rs[i] == '\\' && !unicode.IsDigit(rs[i+2]):
				result = append(result, rs[i+1])
				i++

			case rs[i] == '\\' && unicode.IsDigit(rs[i+2]):

				temp := make([]rune, 0)
				for j := i + 2; j < ln; j++ {
					if unicode.IsDigit(rs[j]) {
						temp = append(temp, rs[j])

					} else {
						break
					}
				}
				fmt.Println(temp)
				num, err := strconv.Atoi(string(temp))
				if err != nil {
					fmt.Println("Reading symbol count went wrong")
				}
				result = append(result, loopChar(rs[i+1], num)...)
				i++

			case !unicode.IsDigit(rs[i]) && !unicode.IsDigit(rs[i+1]) && rs[i] != '\\':
				result = append(result, rs[i])

			case !unicode.IsDigit(rs[i]) && unicode.IsDigit(rs[i+1]) && rs[i] != '\\':
				temp := make([]rune, 0)
				for j := i + 1; j < ln; j++ {
					if unicode.IsDigit(rs[j]) {
						temp = append(temp, rs[j])
						//	i++
					} else {
						break
					}
				}
				num, err := strconv.Atoi(string(temp))
				if err != nil {
					fmt.Println("Reading symbol count went wrong")
				}
				result = append(result, loopChar(rs[i], num)...)
			}
		} else if i == ln-1 {
			if !unicode.IsDigit(rs[i]) {
				result = append(result, rs[i])
			}
		}
	}
	return string(result), nil
}
