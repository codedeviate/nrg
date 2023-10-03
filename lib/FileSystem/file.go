package FileSystem

import "os"

func ReadFile(filename string) (string, error) {
	fileContent, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	// Convert []byte to string
	return string(fileContent), nil
}
