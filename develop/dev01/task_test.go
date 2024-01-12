package dev01_test

import (
	"testing"

	"github.com/Karanth1r3/l_2/develop/dev01/task"
)

func TestPrintTime(t *testing.T) {
	s, err := task.PrintTime()
	if err != nil {
		t.Fatal(err)
	}
	unexpected := ""
	if s == unexpected {
		t.Fatal("Unexpected")
	}
}
