package basics

import (
	"os"
	"path/filepath"
)

// FilesUnderPath returns a list of files contained below the given path
func FilesUnderPath(path string) ([]string, error) {
	var result []string

	f, err := os.Open(path)
	if err != nil {
		return result, err
	}
	defer f.Close()

	infos, err := f.Readdir(0)
	if err != nil {
		return result, err
	}

	for _, info := range infos {
		fullPath := filepath.Join(path, info.Name())

		if info.IsDir() {
			subfiles, err := FilesUnderPath(fullPath)
			if err != nil {
				return result, err
			}
			result = append(result, subfiles...)
		} else {
			result = append(result, fullPath)
		}
	}

	return result[:], nil
}
