package dev02_test

import (
	"testing"
	"unicode"

	"github.com/Karanth1r3/l_2/develop/dev02"
)

func TestDev02(t *testing.T) {
	tests := []struct {
		inputData string
		expected  string
		isOk      bool
	}{
		{
			inputData: "ew4r2k",
			expected:  "ewwwwrrk",
			isOk:      true,
		},
		{
			inputData: "45",
			expected:  "",
			isOk:      false,
		},
		{
			inputData: "",
			expected:  "",
			isOk:      false,
		},
		{
			inputData: "e5",
			expected:  "eeeee",
			isOk:      true,
		},
		{
			inputData: "e5",
			expected:  "eeeee",
			isOk:      true,
		},
		{
			inputData: `ew\\5r\1k`,
			expected:  `ew\\\\\r1k`,
			isOk:      true,
		},
		{
			inputData: `ew\\5r2k`,
			expected:  `ew\\\\\rrk`,
			isOk:      true,
		},
		{
			inputData: "0vsvvv",
			expected:  "",
			isOk:      false,
		},
		{
			inputData: "v0svvv",
			expected:  "svvv",
			isOk:      true,
		},
	}
	for _, test := range tests {
		t.Run(test.inputData, func(t *testing.T) {
			res, err := dev02.UnpackString(test.inputData)
			if err != nil {
				if test.isOk {
					t.Fatal("unexpected behaviour")
				}
			}
			if res == "" && test.isOk {
				t.Fatal("unexpected behaviour")
			}
			if res != "" && !test.isOk {
				t.Fatal("Probably smth wrong")
			}
			if len(test.inputData) > 0 {
				if unicode.IsDigit(rune(test.inputData[0])) && test.isOk {
					t.Fatal("unexpected behaviour")
				}
			}
		})
	}
}
