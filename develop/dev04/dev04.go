package dev04

import (
	"fmt"
	"sort"
	"strings"
	"unicode/utf8"
)

func main() {
	dict := []string{"Fuck", "Uckf", "Care", "cufk", "Тяпка", "пятка", "пятак"}
	dictToLower(dict)
	t := getAnagramGroups(&dict)
	fmt.Println(t)

}

func dictToLower(dict []string) {
	for i := 0; i < len(dict); i++ {
		dict[i] = strings.ToLower(dict[i])
	}
}

func hasEqualRuneCount(s1, s2 string) bool {
	return utf8.RuneCountInString(s1) == utf8.RuneCountInString(s2)
}

func getAnagramGroups(dct *[]string) (result map[string][]string) {
	result = make(map[string][]string)

	go func() {

	}()
out:
	for _, elem := range *dct { // ...
		for _, slc := range result {
			for _, val := range slc {
				if val == elem {
					continue out
				}
			}
		}

		findAnagramGroup(elem, dct, result)
	}
	return result
}

func findAnagramGroup(word string, dct *[]string, anagrams map[string][]string) {
	result := make(map[string][]string)
	foundIndexes := []int{}

	for index, elem := range *dct {
		if elem == word { // if words are equal - they aren't anagrams
			continue
		}
		if !hasEqualRuneCount(elem, word) { // if words have different symbol count - they aren't anagrams
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
		if string(orChars) == string(chars) { // if sorted strings are equal - they are anagrams. Check for original word & original element equality was before
			result[word] = append(result[word], elem)
			foundIndexes = append(foundIndexes, index)
		}
	}
	if len(result[word]) > 1 { // if  group has at least 2 elements, it can be added to anagrams map
		rs := []string{}
		rs = append(rs, result[word]...)
		anagrams[word] = rs
	}

	/*// block for deleting found elements from dictionary
	  for _, e := range foundIndexes {
	    *dct = slices.Delete(*dct, e)
	  }
	*/
	/*
	   for _, e := range result[word] {
	     for i, e2 := range *dct {
	       if e == e2 {
	         //fmt.Println(*dct, i)
	         *dct = append((*dct)[:i], (*dct)[i+1:]...)
	         //fmt.Println("deleted")
	       }
	     }
	   }
	*/
	fmt.Println(*dct)

}
