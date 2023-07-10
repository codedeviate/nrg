package helper

import (
	"os"
	"path/filepath"
	"strings"
)

func Contains(needle interface{}, haystack []interface{}) bool {
	for _, value := range haystack {
		if value == needle {
			return true
		}
	}
	return false
}

func arrayIncludesString(array []string, value string) bool {
	for _, v := range array {
		if v == value {
			return true
		}
	}
	return false
}

func arrayStringIndex(array []string, value string) int {
	for i, v := range array {
		if v == value {
			return i
		}
	}
	return -1
}

func arrayIncludesStringInArray(array []string, values []string) (bool, string) {
	for _, v := range values {
		if arrayIncludesString(array, v) {
			return true, v
		}
	}
	return false, ""
}

func removeFirstStringFromArray(array []string, value string) []string {
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

func CloneCommand(command *Command) *Command {
	clonedCommand := &Command{}
	clonedCommand.Commands = append([]string{}, command.Commands...)
	clonedCommand.Args = append([]string{}, command.Args...)
	clonedCommand.Env = make(map[string]string, len(command.Env))
	for key, value := range command.Env {
		clonedCommand.Env[key] = value
	}
	clonedCommand.Flags = command.Flags
	clonedCommand.Passthru = command.Passthru && true
	return clonedCommand
}

func CanRun(scriptName string) bool {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return false
	}
	stack := GetStack()
	scriptName = strings.TrimSpace(scriptName)
	scriptName = strings.TrimPrefix(scriptName, "@")
	scriptName = strings.TrimSuffix(scriptName, ".js")
	var scriptPath string
	if stack.ActiveProject != nil {
		projectDir := GetStack().ActiveProject.Path
		for _, p := range GetStack().ActiveProject.ScriptPaths {
			scriptPath = filepath.Join(projectDir, p, scriptName+".js")
			if _, err := os.Stat(scriptPath); !os.IsNotExist(err) {
				break
			}
		}
	}
	// Check if the script is in the project scripts folder
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) || len(scriptPath) == 0 {
		if stack.ActiveProject != nil {
			projectDir := GetStack().ActiveProject.Path
			scriptPath = filepath.Join(projectDir, "/.nrg/scripts/", scriptName+".js")
		}
	}
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		scriptPath = filepath.Join(homeDir, "/.nrg/scripts/", scriptName+".js")
		if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func RunesToString(runes []rune) string {
	return string(runes)
}
func StringToRunes(str string) []rune {
	return []rune(str)
}
