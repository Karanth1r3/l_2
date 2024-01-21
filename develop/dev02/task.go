package dev02

import (
	"fmt"
	"strconv"
	"unicode"
	"unicode/utf8"
)

func main() {
	(UnpackString("ew4r2k"))
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

			case !unicode.IsDigit(rs[i]) && !unicode.IsDigit(rs[i+1]):
				result = append(result, rs[i])
			case !unicode.IsDigit(rs[i]) && unicode.IsDigit(rs[i+1]):
				temp := make([]rune, 0)
				for j := i + 1; j < ln; j++ {
					if unicode.IsDigit(rs[j]) {
						//	fmt.Println(string(rs[i]))
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
				fmt.Println(num)
				result = append(result, loopChar(rs[i], num)...)
				//	fmt.Println(string(rs))
			}
		} else if i == ln-1 {
			if !unicode.IsDigit(rs[i]) {
				result = append(result, rs[i])
			}
		}
	}
	return string(result), nil
}
