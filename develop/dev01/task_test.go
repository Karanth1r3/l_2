package dev01_test

import (
	"testing"

	"github.com/Karanth1r3/l_2/develop/dev01"
)

func TestTime(t *testing.T) {

	tests := []struct {
		address string
		ok      bool
	}{
		{
			address: "0.beevik-ntp.pool.ntp.org",
			ok:      true,
		},
		{
			address: "",
			ok:      false,
		},
	}

	for _, test := range tests {
		t.Run(test.address, func(t *testing.T) {
			res, err := dev01.PrintTime(test.address)
			if err != nil {
				if test.ok {
					t.Fatal(err)
				}
			}
			if res == "" && test.ok {
				t.Fatal("unexpected behaviour")
			}
		})
	}
}
