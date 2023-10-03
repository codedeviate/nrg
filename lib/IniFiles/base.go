package IniFiles

import (
	"nrg/lib/FileSystem"
	"nrg/lib/Strings"
	"strings"
)

func ParseIniFile(filename string) (map[string]map[string]string, error) {
	data, err := FileSystem.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return ParseIni(data)
}

func ParseIni(data string) (map[string]map[string]string, error) {
	lines := Strings.SplitLines(data)
	// Parse INI file
	var section string
	var err error
	var key string
	var value string
	var sectionMap map[string]string
	var iniMap map[string]map[string]string
	iniMap = make(map[string]map[string]string)
	sectionMap = make(map[string]string)
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		if line[0] == '[' {
			section = strings.TrimSpace(line[1 : len(line)-1])
			sectionMap = make(map[string]string)
			iniMap[section] = sectionMap
			continue
		}
		key, value, err = ParseIniLine(line)
		if err != nil {
			return nil, err
		}
		sectionMap[key] = value
	}
	return iniMap, nil
}

func ParseIniLine(line string) (string, string, error) {
	// Parse INI line
	var key string
	var value string
	var i int
	i = 0
	for i < len(line) {
		if line[i] == '=' {
			break
		}
		i++
	}
	if i >= len(line) {
		return "", "", nil
	}
	key = strings.TrimSpace(line[0:i])
	value = strings.TrimSpace(line[i+1:])
	return key, value, nil
}
