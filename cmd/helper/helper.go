package helper

import (
	"os"
	"path/filepath"
	"strings"
)

func Contains(files []string, file string) bool {
	for _, f := range files {
		if strings.Contains(file, f) {
			return true
		}
	}
	return false
}

func ListDirectory(path string) ([]string, error) {
	var directories []string

	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			directories = append(directories, filepath.Join(path, entry.Name()))
		}
	}

	return directories, nil
}

func RecursiveDirectorySearch(rootPath string) ([]string, error) {
	var directories []string
	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			directories = append(directories, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return directories, nil
}
