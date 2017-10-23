package utils

import (
	"io"
	"os"
)

// CreateDir creates the directory at path.
func CreateDir(path string) error {

	err := os.MkdirAll(path, 0755)

	return err

}

// DirExists checks if the path exists.
func DirExists(path string) bool {

	if _, err := os.Stat(path); err == nil {
		return true
	}

	return false

}

// RemoveDir deletes the path and everything contained in it.
func RemoveDir(path string) error {

	err := os.RemoveAll(path)

	return err

}

// IsDirEmpty checks if the given directory is empty or not.
func IsDirEmpty(path string) bool {

	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()

	// Read first file
	_, err = f.Readdir(1)

	// If EOF then the Dir is empty
	if err == io.EOF {
		return true
	}

	return false

}
