package jsio

import (
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var screenWidth = 0
var screenHeight = 0
var lastScreenUpdate int64 = 0

func UpdateScreenSize() {
	if time.Now().Unix()-lastScreenUpdate > 1 {
		cmd := exec.Command("stty", "size")
		cmd.Stdin = os.Stdin
		out, err := cmd.Output()
		if err != nil {
			log.Fatal(err)
		}
		parts := strings.Split(string(out), " ")
		screenWidth, err = strconv.Atoi(strings.TrimSpace(parts[1]))
		if err != nil {
			log.Fatal(err)
		}
		screenHeight, err = strconv.Atoi(strings.TrimSpace(parts[0]))
		lastScreenUpdate = time.Now().Unix()
	}
}

func GetScreenWidth() int {
	UpdateScreenSize()
	return screenWidth
}

func GetScreenHeight() int {
	UpdateScreenSize()
	return screenHeight
}
