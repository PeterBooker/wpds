package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func getLatestRevision(dir string) (int, error) {

	var revision int
	var rURL string

	c := &http.Client{
		Timeout: 30 * time.Second,
	}

	switch dir {
	case "plugins":
		rURL = wpLatestPluginsRevisionURL
	case "themes":
		rURL = wpLatestThemesRevisionURL
	}

	resp, err := c.Get(rURL)
	if err != nil {
		return 0, err
	}

	if resp.StatusCode != 200 {
		return 0, fmt.Errorf("Invalid HTTP Response")
	}

	defer resp.Body.Close()
	bBytes, err := ioutil.ReadAll(resp.Body)
	bString := string(bBytes)

	revs := regexRevision.FindAllStringSubmatch(bString, 1)

	revision, err = strconv.Atoi(revs[0][1])
	if err != nil {
		return 0, err
	}

	return revision, nil

}

func getCurrentRevision(dir string) (int, error) {

	var revision int

	fname := ".last-revision"

	path := filepath.Join(wd, dir, fname)

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return 0, err
	}

	revision, err = strconv.Atoi(string(data))
	if err != nil {
		return 0, err
	}

	return revision, nil

}

func setCurrentRevision(revision int, dir string) error {

	fname := ".last-revision"

	path := filepath.Join(wd, dir, fname)

	f, err := os.Create(path)
	defer f.Close()
	if err != nil {
		return err
	}

	_, err = io.WriteString(f, string(revision))
	if err != nil {
		return err
	}

	return nil

}
