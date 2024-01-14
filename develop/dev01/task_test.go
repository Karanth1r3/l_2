package dev01_test

import (
	"testing"

	"github.com/Karanth1r3/l_2/develop/dev01"
)

func TestTime(t *testing.T) {
	r, err := dev01.PrintTime()

	if err != nil {
		t.Fatal(err)
	}
	unexpected := ""
	if r == unexpected {
		t.Fatal("unexpected behaviour")
	}
}
