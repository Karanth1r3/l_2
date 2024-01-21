package main

import (
	"fmt"
	"slices"
	"sort"
	"strings"
	"sync"
	"unicode/utf8"
)

func main() {
	dict := []string{"Fuck", "Uckf", "Care", "cufk", "бука", "убак", "Тяпка", "куба", "пятка", "пятак", "fork", "rofk"}
	dictToLower(dict)
	fmt.Println(&dict)
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
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	for _, elem := range *dct { // ...
		wg.Add(1)
		go findAnagramGroup(&wg, &mu, elem, dct, result)
	}
	wg.Wait()
	return result
}

func findAnagramGroup(wg *sync.WaitGroup, mu *sync.Mutex, word string, dct *[]string, anagrams map[string][]string) {
	mu.Lock()
	result := make(map[string][]string)
	found := []string{}
	//fmt.Println(word)
	for _, elem := range *dct {
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
			found = append(found, elem)
		}
	}
	if len(result[word]) > 1 { // if  group has at least 2 elements, it can be added to anagrams map
		rs := []string{}
		sort.StringSlice.Sort(result[word])
		rs = append(rs, result[word]...)
		anagrams[word] = rs
	}

	for _, e := range found {
		*dct = slices.DeleteFunc(*dct, func(s string) bool {
			return s == e
		})
	}
	*dct = slices.DeleteFunc(*dct, func(s string) bool { // Can delete only [0] actually
		return s == word
	})

	fmt.Println(*dct)
	mu.Unlock()
	wg.Done()
}
