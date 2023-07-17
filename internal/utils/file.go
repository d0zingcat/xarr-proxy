package utils

import (
	"io/fs"
	"os"
	"path/filepath"
	"sort"
)

// read file and return content as string
func ReadFile(path string) (string, error) {
	// check if the file exists
	if !FileExists(path) {
		return "", os.ErrNotExist
	}
	// read file
	res, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(res), nil
}

// check if the file FileExists
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// walk over the dir and find all the files, order by filename asc
func WalkDir(dir string, fn func(filename string) bool) ([]string, error) {
	var files []string
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// filter func
		if fn != nil {
			if fn(path) {
				files = append(files, path)
			}
			return nil
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Slice(files, func(i, j int) bool {
		return files[i] < files[j]
	})
	return files, nil
}

func CreateFile(path string, content *string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	if content != nil {
		_, err := f.WriteString(*content)
		if err != nil {
			return err
		}
	}
	return nil
}
