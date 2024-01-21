package dev04_test

import (
	"fmt"
	"reflect"
	"sort"
	"testing"

	"github.com/Karanth1r3/l_2/develop/dev04"
)

func TestAnagram(t *testing.T) {
	tests := []struct {
		testID         string
		inputData      []string
		expectedLen    int
		expectedGroups [][]string
		isOk           bool
	}{
		{
			testID:         "0",
			inputData:      []string{"Fuck", "Uckf", "Care", "cufk", "бука", "убак", "Тяпка", "куба", "пятка", "пятак", "fork", "rofk"},
			expectedLen:    3, // Concurrency makes map composition unpredictable => i guess i can only compare lengths
			expectedGroups: [][]string{{"cufk", "fuck", "uckf"}, {"бука", "куба", "убак"}, {"пятак", "пятка", "тяпка"}},
			isOk:           true,
		},
	}
	for _, test := range tests {
		t.Run(test.testID, func(t *testing.T) {
			//	fmt.Println(dev04.GetAnagramGroups(&test.inputData))
			x := make([][]string, 0)
			res := dev04.GetAnagramGroups(&test.inputData)
			for k, _ := range res {
				rs := []string{k}
				for _, elem := range res[k] {
					rs = append(rs, elem)
				}
				sort.StringSlice.Sort(rs)
				x = append(x, rs)
			}
			s1 := fmt.Sprintf("%v", x)
			s2 := fmt.Sprintf("%v", test.expectedGroups)
			s3 := fmt.Sprintf("expected: %s, got: %s", s2, s1)
			if len(res) != test.expectedLen && test.isOk {
				s := fmt.Sprintf("expected: %d, got: %d", test.expectedLen, len(res))
				//s := fmt.Sprintf("%v", res)
				t.Fatal("unexpected behaviour: ", s)
			}
			if !compare(x, test.expectedGroups) && test.isOk {
				t.Fatal(s3)
			}
			if compare(x, test.expectedGroups) && !test.isOk {
				t.Fatal("unexpected behaviour")
			}
		})
	}
}

func compare(s1 [][]string, s2 [][]string) bool {
	equalityCount := len(s1)
	count := 0
	for _, elem := range s1 {
		for _, elem2 := range s2 {
			if reflect.DeepEqual(elem, elem2) {
				count++
			}
		}
	}
	if equalityCount == count {
		return true
	} else {
		return false
	}
}
