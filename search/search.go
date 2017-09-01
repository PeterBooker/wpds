package search

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

// Match holds data about search regex matches.
type Match struct {
	filename string
	index    int
	snippet  string
}

func NewStringSearch(regex string, dir string) error {

	/*
		results := make(chan Match)
		var searchFinished bool

		for {

			if searchFinished {
				return nil
			}

			select {
			case match, open := <-results:

				if !open {

					searchFinished = true

				} else {

					fmt.Println(match)

				}

			}

		}*/

	re, err := regexp.Compile(regex)
	if err != nil {
		fmt.Println("Invalid StringSearch Regex")
		os.Exit(1)
	}

	wd, _ := os.Getwd()

	path := filepath.Join(wd, dir)

	err = filepath.Walk(path, searchFile(re))
	if err != nil {

	}

	return nil

}

func searchFile(regex *regexp.Regexp) filepath.WalkFunc {

	var results [][]int

	return func(path string, fi os.FileInfo, err error) error {

		// We cannot do anything, move on
		// Perhaps we can warn the user in the future?
		if err != nil {
			return nil
		}

		// If target is a directory, move on
		if fi.IsDir() {
			return nil
		}

		data, err := ioutil.ReadFile(path)
		if err != nil {
			// TODO: handle this in the future
			return nil
		}

		results = regex.FindAllIndex(data, -1)
		if results != nil {
			for _, result := range results {

				br := bytes.NewReader(data)
				br.ReadByte()
				//fmt.Println(result)
				//fmt.Printf("Match: %v\n", result)
				fmt.Printf("File: %s, Position: %v\n", fi.Name(), result)
			}

		}

		return nil

	}

}
