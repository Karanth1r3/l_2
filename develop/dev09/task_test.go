package dev09_test

import (
	"os"
	"path"
	"testing"

	"github.com/Karanth1r3/l_2/develop/dev09"
)

func TestWGet(t *testing.T) {
	t.Skip("Dev purpose")
	tests := []struct {
		url  string
		err  error
		isOk bool
	}{
		{
			url:  "https://ru.wikipedia.org/wiki/Telnet",
			err:  nil,
			isOk: true,
		},
		{
			url:  "",
			isOk: false,
		},
	}
	for _, test := range tests {
		t.Run(test.url, func(t *testing.T) {
			errr := dev09.Wget(test.url)
			if errr != nil && test.isOk {
				t.Fatal("unexpected error")
			}
			os.Remove(path.Base(test.url))
		})
	}
}
