package Strings

import "strings"

func ArrayIncludesString(array []string, value string) bool {
	for _, v := range array {
		if v == value {
			return true
		}
	}
	return false
}

func Contains(needle interface{}, haystack []interface{}) bool {
	for _, value := range haystack {
		if value == needle {
			return true
		}
	}
	return false
}

func ArrayStringIndex(array []string, value string) int {
	for i, v := range array {
		if v == value {
			return i
		}
	}
	return -1
}

func ArrayIncludesStringInArray(array []string, values []string) (bool, string) {
	for _, v := range values {
		if ArrayIncludesString(array, v) {
			return true, v
		}
	}
	return false, ""
}

func RemoveFirstStringFromArray(array []string, value string) []string {
	newArray := []string{}
	firstFound := false
	for _, val := range array {
		if val == value && firstFound == false {
			firstFound = true
			continue
		}
		newArray = append(newArray, val)
	}
	return newArray
}

func SplitLines(data string) []string {
	// Split string into lines
	var lines []string
	lines = strings.Split(data, "\n")
	return lines
}
