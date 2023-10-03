package jsio

import (
	"fmt"
	"strings"
)

func Readln() string {
	var input string
	fmt.Scanln(&input)
	return input
}

func ReadYN(defaultInput string) bool {
	var input string
	fmt.Scanln(&input)

	if strings.ToLower(input) == "y" || strings.ToLower(defaultInput) == "y" {
		return true
	}
	return false
}
