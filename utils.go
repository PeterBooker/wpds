package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"unicode"
	"unicode/utf8"
)

// isConfirmationRequired checks if files already exist in the destination and the user
// needs to confirm beginning a fresh download.
func isConfirmationRequired(dir string) bool {

	path := filepath.Join(wd, dir)

	// Check if directory exists
	if _, err := os.Stat(path); !os.IsNotExist(err) {

		// If exists, check if empty
		empty, _ := isDirEmpty(path)

		if empty == false {
			return true
		}

	}

	return false

}

// getUserConfirmation prompts the user to confirm they are happy starting a fresh download
// to avoid accidentally polluting any current files.
func getUserConfirmation() bool {

	// Ask for confirmation
	fmt.Println("Are you sure you want to begin a fresh download? (Y/N)")

	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		fmt.Println(err)
	}

	runes, _ := utf8.DecodeRuneInString(response)
	response = string(unicode.ToLower(runes))

	if string(response[0]) == "y" {
		return true
	}

	return false

}

// isDirEmpty checks if the given directory is empty or not.
func isDirEmpty(dir string) (bool, error) {

	f, err := os.Open(dir)
	if err != nil {
		return false, err
	}
	defer f.Close()

	// Read first file
	_, err = f.Readdir(1)

	// If EOF then the Dir is empty
	if err == io.EOF {
		return true, nil
	}

	return false, err

}

func removeDuplicates(items *[]string) {

	found := make(map[string]bool)
	i := 0
	for k, v := range *items {
		if !found[v] {
			found[v] = true
			(*items)[i] = (*items)[k]
			i++
		}
	}
	*items = (*items)[:i]

}
