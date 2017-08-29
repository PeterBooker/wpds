package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

func getLatestRevision(dir string) (string, error) {

	var revision string
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
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("Invalid HTTP Response")
	}

	defer resp.Body.Close()
	bBytes, err := ioutil.ReadAll(resp.Body)
	bString := string(bBytes)

	regex := regexp.MustCompile("[0-9]+")
	revs := regex.FindAllString(bString, 1)

	revision = revs[0]

	return revision, nil

}

func getCurrentRevision(dir string) (string, error) {

	var revision string

	fname := ".last-revision"

	path := filepath.Join(wd, dir, fname)

	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		return "", err
	}

	return revision, nil

}

func setCurrentRevision(revision string, dir string) error {

	fname := ".last-revision"

	path := filepath.Join(wd, dir, fname)

	f, err := os.Create(path)
	defer f.Close()
	if err != nil {
		return err
	}

	_, err = io.WriteString(f, revision)
	if err != nil {
		return err
	}

	return nil

}
