package utils

import (
	"os"
)

func ReadDir(dir string) []os.FileInfo {
	d, err := os.Open(dir)

	if err != nil {
		panic(err)
	}
	defer d.Close()

	files, err := d.Readdir(-1)
	if err != nil {
		panic(err)
	}

	filteredDirs := make([]os.FileInfo, 0)
	for _, file := range files {
		if file.IsDir() {
			filteredDirs = append(filteredDirs, file)
		}
	}

	return filteredDirs
}
