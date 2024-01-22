package dev05_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/Karanth1r3/l_2/develop/dev05"
)

func TestGrep(t *testing.T) {
	tests := []struct {
		testID                                    string
		inputData                                 []string
		typedString                               string
		after, before, context                    int
		ignoreCase, invert, fixed, lineNum, count bool // Linenum is harder to test (in this way at least), so it will make no difference here
		expectedResult                            map[int]string
		isOk                                      bool
	}{
		{
			testID:         "0",
			inputData:      []string{"Fuck", "Uckf", "Care", "cufk", "бука", "убак", "Тяпка", "куба", "пятка", "пятак", "fork", "rofk"},
			typedString:    "Fuck",
			after:          0,
			before:         0,
			context:        0,
			ignoreCase:     false,
			invert:         false,
			fixed:          false,
			lineNum:        false,
			count:          false,
			isOk:           true,
			expectedResult: map[int]string{0: "Fuck"},
		},
		{
			testID:         "1",
			inputData:      []string{"Fuck", "Uckf", "Care", "cufk", "бука", "убак", "Тяпка", "куба", "пятка", "пятак", "fork", "rofk"},
			typedString:    "Fuck",
			after:          0,
			before:         0,
			context:        0,
			ignoreCase:     false,
			invert:         false,
			fixed:          false,
			lineNum:        false,
			count:          true,
			isOk:           true,
			expectedResult: map[int]string{0: "1"},
		},
		{
			testID:         "2",
			inputData:      []string{"Fuck", "Uckf", "Care", "cufk", "бука", "убак", "Тяпка", "куба", "пятка", "пятак", "fork", "rofk"},
			typedString:    "fuck",
			after:          0,
			before:         0,
			context:        2,
			ignoreCase:     true,
			invert:         false,
			fixed:          false,
			lineNum:        true,
			count:          false,
			isOk:           true,
			expectedResult: map[int]string{0: "Fuck"},
		},
		{
			testID:         "3",
			inputData:      []string{"Fuck", "Uckf", "Care", "cufk", "бука", "убак", "Тяпка", "куба", "пятка", "пятак", "fork", "rofk"},
			typedString:    "fuck",
			after:          0,
			before:         0,
			context:        2,
			ignoreCase:     false,
			invert:         true,
			fixed:          true,
			lineNum:        true, // Will make no difference
			count:          false,
			isOk:           true,
			expectedResult: map[int]string{0: "Fuck", 1: "Uckf", 2: "Care", 3: "cufk", 4: "бука", 5: "убак", 6: "Тяпка", 7: "куба", 8: "пятка", 9: "пятак", 10: "fork", 11: "rofk"},
		},
	}
	for _, test := range tests {
		t.Run(test.testID, func(t *testing.T) {
			//	fmt.Println(dev04.GetAnagramGroups(&test.inputData))
			m := dev05.PrintGrep(test.inputData, test.typedString, test.after,
				test.before, test.context, test.ignoreCase, test.invert,
				test.fixed, test.lineNum, test.count)
			if !reflect.DeepEqual(m, test.expectedResult) && test.isOk {
				s := fmt.Sprintf("got: %v, expected: %v", m, test.expectedResult)
				t.Fatal("unexpected result", s)
			}
			if !reflect.DeepEqual(m, test.expectedResult) && !test.isOk {
				t.Fatal("unexpected result")
			}
		})
	}
}
