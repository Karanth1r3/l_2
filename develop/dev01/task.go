package dev01

import (
	"fmt"

	"github.com/beevik/ntp"
)

func main() {
	PrintTime("0.beevik-ntp.pool.ntp.org")
}

func PrintTime(address string) (string, error) {
	time, err := ntp.Time(address)
	if err != nil {
		return "", err
	}
	res := fmt.Sprint(time)
	fmt.Println(res)
	return res, nil
}
