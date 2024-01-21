package dev04_test

import (
	"fmt"
	"testing"

	"github.com/Karanth1r3/l_2/develop/dev04"
)

func TestAnagram(t *testing.T) {
	tests := []struct {
		testID      string
		inputData   []string
		expectedLen int
		isOk        bool
	}{
		{
			testID:      "0",
			inputData:   []string{"Fuck", "Uckf", "Care", "cufk", "бука", "убак", "Тяпка", "куба", "пятка", "пятак", "fork", "rofk"},
			expectedLen: 3, // Concurrency makes map composition unpredictable => i guess i can only compare lengths
			isOk:        true,
		},
	}
	for _, test := range tests {
		t.Run(test.testID, func(t *testing.T) {
			//	fmt.Println(dev04.GetAnagramGroups(&test.inputData))
			res := dev04.GetAnagramGroups(&test.inputData)
			if len(res) != test.expectedLen && test.isOk {
				s := fmt.Sprintf("expected: %d, got: %d", test.expectedLen, len(res))
				//s := fmt.Sprintf("%v", res)
				t.Fatal("unexpected behaviour: ", s)
			}
		})
	}
}
