package utils

import (
	"io/fs"
	"os"
	"path/filepath"
	"sort"
)

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

func CreateFile(path string) error {
	_, err := os.Create(path)
	if err != nil {
		return err
	}
	return nil
}
