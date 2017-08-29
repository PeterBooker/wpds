package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/cheggaaa/pb.v1"
)

func getAllItems(dir string) {

	c := &http.Client{
		Timeout: 30 * time.Second,
	}

	var eURL string

	switch dir {
	case "plugins":
		eURL = wpAllPluginsListURL
	case "themes":
		eURL = wpAllThemesListURL
	}

	resp, err := c.Get(eURL)
	if err != nil {
		fmt.Printf("Failed HTTP GET of updated %s.\n", dir)
		os.Exit(1)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Println("Invalid HTTP Response")
		os.Exit(1)
	}

	items, revision, err := parseItemListHTML(resp.Body)
	if err != nil {
		fmt.Println("Failed parsing HTML response. Could not get items list.")
		os.Exit(1)
	}

	fetchItems(items, dir, 10)

	err = setCurrentRevision(revision, "plugins")
	if err != nil {
		fmt.Println("The current revision could not be saved, updating will not work.")
	}

}

func getUpdatedItems(dir string, rev int) {

	lrev, err := getLatestRevision(dir)
	if err != nil {
		fmt.Println("Cannot get the latest revision, updating cancelled.")
		os.Exit(1)
	}

	rdiff := lrev - rev

	c := &http.Client{
		Timeout: 30 * time.Second,
	}

	var eURL string

	switch dir {
	case "plugins":
		eURL = encodeURL(fmt.Sprintf(wpPluginChangelogURL, lrev, rdiff))
	case "themes":
		eURL = encodeURL(fmt.Sprintf(wpThemeChangelogURL, lrev, rdiff))
	}

	resp, err := c.Get(eURL)
	if err != nil {
		fmt.Printf("Failed HTTP GET of updated %s.\n", dir)
		os.Exit(1)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Println("Invalid HTTP Response")
		os.Exit(1)
	}

	bBytes, err := ioutil.ReadAll(resp.Body)
	bString := string(bBytes)

	items := regexUpdatedItems.FindAllString(bString, -1)

	removeDuplicates(&items)

	fetchItems(items, dir, 10)

	err = setCurrentRevision(rev, "plugins")
	if err != nil {
		fmt.Println("The current revision could not be saved, updating will not work.")
	}

}

func fetchItems(items []string, dir string, limit int) error {

	iCount := len(items)

	bar := pb.StartNew(iCount)

	limiter := make(chan struct{}, limit)

	// Make Plugins Dir ready for extracting plugins
	path := filepath.Join(wd, dir)
	err := mkdir(path)
	if err != nil {
		return err
	}

	for _, v := range items {

		// Will block if more than max Goroutines already running.
		limiter <- struct{}{}
		bar.Increment()

		go func(name string) {
			getItem(name)
			<-limiter
		}(v)

	}

	bar.FinishPrint(fmt.Sprintf("Completed download of %d Plugins.", iCount))

	return nil

}

func getItem(item string) {

	c := &http.Client{
		Timeout: 60 * time.Second,
	}

	eURL := encodeURL(fmt.Sprintf(wpPluginDownloadURL, item))

	resp, err := c.Get(eURL)
	if err != nil {
		fmt.Printf("Error Downloading Plugin: %s\n", item)
		return
	}

	if resp.StatusCode != 200 {

		if resp.StatusCode == 404 {
			return
		}

		fmt.Printf("Error Downloading Plugin: %s Status Code: %d\n", item, resp.StatusCode)
		return
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading Get request body for Plugin: %s\n", item)
		fmt.Println(err)
		return
	}

	err = extract(content, item)

}
