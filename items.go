package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

func getAllItems(dir string, baseURL string) {

}

func getUpdatedItems(dir string, baseURL string) error {

	rev, err := getCurrentRevision(dir)
	if err != nil {
		fmt.Println("No revision found, cannot continue updating. Make sure the .last-revision file exists.")
		os.Exit(1)
	}

	lrev, err := getLatestRevision(dir)
	if err != nil {
		fmt.Println("Cannot get the latest revision, updating cancelled.")
		os.Exit(1)
	}

	irev, err := strconv.Atoi(rev)
	if err != nil {
		fmt.Println("Current revision does not appear to be a valid number.")
		os.Exit(1)
	}

	ilrev, err := strconv.Atoi(lrev)
	if err != nil {
		fmt.Println("Latest revision does not appear to be a valid number.")
		os.Exit(1)
	}

	rdiff := ilrev - irev

	c := &http.Client{
		Timeout: 30 * time.Second,
	}

	var eURL string

	switch dir {
	case "plugins":
		eURL, err = encodeURL(fmt.Sprintf(wpPluginChangelogURL, lrev, string(rdiff)))
		if err != nil {
			return err
		}
	case "themes":
		eURL, err = encodeURL(fmt.Sprintf(wpThemeChangelogURL, lrev, string(rdiff)))
		if err != nil {
			return err
		}
	}

	resp, err := c.Get(eURL)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("Invalid HTTP Response")
	}

	defer resp.Body.Close()
	bBytes, err := ioutil.ReadAll(resp.Body)
	bString := string(bBytes)

	items := regexUpdatedItems.FindAllString(bString, 1)

	removeDuplicates(&items)

	return nil

}
