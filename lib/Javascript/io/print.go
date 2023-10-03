package jsio

import (
	"fmt"
	"strconv"
)

func PrintPadded(s string, length int, rest interface{}) {
	if length <= 0 {
		length = len(s)
	}
	if rest != nil {
		fmt.Print(rest)
	}
}
func PrintPaddedLn(s string, length int, rest interface{}) {
	if length <= 0 {
		length = len(s)
	}
	fmt.Printf("%-"+strconv.Itoa(length)+"s", s)
	if rest != nil {
		fmt.Print(rest)
	}
	fmt.Println()
}
func PrintWidth(s string, length int, rest interface{}) {
	if length <= 0 {
		length = len(s)
	} else if len(s) > length {
		s = s[0:length]
	}
	fmt.Printf("%-"+strconv.Itoa(length)+"s", s)
	if rest != nil {
		fmt.Print(rest)
	}
}
func PrintWidthLn(s string, length int, rest interface{}) {
	if length <= 0 {
		length = len(s)
	} else if len(s) > length {
		s = s[0:length]
	}
	fmt.Printf("%-"+strconv.Itoa(length)+"s", s)
	if rest != nil {
		fmt.Print(rest)
	}
	fmt.Println()
}

func PrintMidPad(s1 string, length int, s2 string) {
	if length <= 0 {
		length = len(s1) + len(s2)
	}
	fmt.Print(s1)
	fmt.Printf("%"+strconv.Itoa(length-(len(s1)+len(s2)))+"s", "")
	fmt.Print(s2)
}
func PrintMidPadLn(s1 string, length int, s2 string) {
	if length <= 0 {
		length = len(s1) + len(s2)
	}
	fmt.Print(s1)
	fmt.Printf("%"+strconv.Itoa(length-(len(s1)+len(s2)))+"s", "")
	fmt.Println(s2)
}
