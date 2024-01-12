package dev01_test

import (
	"dev01/task"
	"testing"
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
