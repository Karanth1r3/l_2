package dev04

import (
	"fmt"
	"sort"
	"strings"
	"unicode/utf8"
)

func main() {
	dict := []string{"Fuck", "Uckf", "Care"}
	dictToLower(dict)
	fmt.Println(dict)
	t := getAnagramGroups(&dict)
	fmt.Println(t)

}

func dictToLower(dict []string) {
	for _, elem := range dict {
		elem = strings.ToLower(elem)
	}
}

func hasEqualRuneCount(s1, s2 string) bool {
	return utf8.RuneCountInString(s1) == utf8.RuneCountInString(s2)
}

func toLowercase(s string) string {
	return strings.ToLower(s)
}

func getAnagramGroups(dct *[]string) (result map[string][]string) {
	result = make(map[string][]string)
	for _, elem := range *dct {
		findAnagramGroup(elem, dct, result)
	}
	return result
}

func findAnagramGroup(word string, dct *[]string, anagrams map[string][]string) {
	result := make(map[string][]string)
	for index, elem := range *dct {
		if elem == word { // if words are equal - they aren't anagrams
			fmt.Println("==")
			continue
		}
		if !hasEqualRuneCount(elem, word) { // if words have different symbol count - they aren't anagrams
			fmt.Println("==Count")
			continue
		}
		chars := []rune(elem)
		orChars := []rune(word)
		sort.Slice(orChars, func(i, j int) bool { // sort word and element of the dictionary alphabetically
			return orChars[i] < orChars[j]
		})
		sort.Slice(chars, func(i, j int) bool {
			return chars[i] < chars[j]
		})
		fmt.Println(string(chars), string(orChars))
		if string(orChars) == string(chars) { // if sorted strings are equal - they are anagrams. Check for original word & original element equality was before
			fmt.Println("Found")
			result[word] = append(result[word], elem)
			*dct = append((*dct)[:index], (*dct)[index+1:]...)
		}
	}
	if len(result) > 1 { // if  group has at least 2 elements, it can be added to anagrams map
		anagrams[word] = result[word]
	}
}
