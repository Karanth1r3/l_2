package dev01

import (
	"fmt"

	"github.com/beevik/ntp"
)

func main() {
	PrintTime()
}

func PrintTime() (string, error) {
	time, err := ntp.Time("0.beevik-ntp.pool.ntp.org")
	if err != nil {
		return "", err
	}
	res := fmt.Sprint(time)
	fmt.Println(res)
	return res, nil
}
