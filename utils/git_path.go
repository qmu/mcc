package utils

import (
	"io/ioutil"
	"path/filepath"
)

// GetDotGitPath returns the .git path from the arg path
func GetDotGitPath(path string) (result string, err error) {
	dir, err := filepath.Abs(path)
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		dir = filepath.Dir(path)
		return GetDotGitPath(dir)
	}

	for _, file := range files {
		if file.IsDir() && file.Name() == ".git" {
			return dir, err
		}
	}

	if dir != "/" {
		dir = filepath.Dir(dir)
		return GetDotGitPath(dir)
	}
	return result, err
}
