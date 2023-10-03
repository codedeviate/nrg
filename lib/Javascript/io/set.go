package jsio

import "fmt"

func SetBold() {
	fmt.Print("\033[1m")
}

func SetNormal() {
	fmt.Print("\033[0m")
}

func SetRed(text string) {
	fmt.Print("\033[31m")
}

func SetGreen(text string) {
	fmt.Print("\033[32m")
}

func SetYellow(text string) {
	fmt.Print("\033[33m")
}

func SetBlue(text string) {
	fmt.Print("\033[34m")
}

func SetMagenta(text string) {
	fmt.Print("\033[35m")
}

func SetCyan(text string) {
	fmt.Print("\033[36m")
}

func SetWhite(text string) {
	fmt.Print("\033[37m")
}

func SetTabTitle(title string) {
	fmt.Printf("\033]0;%s\007", title)
}
