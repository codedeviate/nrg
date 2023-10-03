package jsfs

import (
	"fmt"
	"os"
	"path/filepath"
)

func FindFilename(filename string) []string {
	paths := []string{}
	err := filepath.Walk(".",
		func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				if info.Name() == filename {
					paths = append(paths, path)
					return filepath.SkipDir
				}
			}
			return nil
		})
	if err != nil {
		fmt.Println(err)
		return []string{}
	}
	return paths
}

func FindDirname(dirname string) []string {
	paths := []string{}
	err := filepath.WalkDir(".",
		func(path string, info os.DirEntry, err error) error {
			if info.Name() == dirname {
				paths = append(paths, path)
				return filepath.SkipDir
			}
			return nil
		})
	if err != nil {
		fmt.Println(err)
		return []string{}
	}
	return paths
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
