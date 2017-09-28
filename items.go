package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/urfave/cli"
	"gopkg.in/cheggaaa/pb.v1"
)

func getAllItems(ctx *cli.Context, dir string) {

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

	fetchItems(items, dir, ctx.Int("concurrent-actions"))

	err = setCurrentRevision(revision, dir)
	if err != nil {
		fmt.Println("The current revision could not be saved, updating will not work.")
	}

}

func getUpdatedItems(ctx *cli.Context, dir string, rev int) {

	var items []string

	lastRev, err := getLatestRevision(dir)
	if err != nil {
		fmt.Println("Cannot get the latest revision, updating cancelled.")
		os.Exit(1)
	}

	if rev == lastRev {
		fmt.Printf("You are currently at the latest revision: %d. No update needed.\n", rev)
		os.Exit(1)
	}

	items = getItemsList(dir, rev, lastRev)

	fetchItems(items, dir, ctx.Int("concurrent-actions"))

	err = setCurrentRevision(lastRev, dir)
	if err != nil {
		fmt.Println("The current revision could not be saved, updating will not work.")
	}

}

func getItemsList(dir string, rev int, lastRev int) []string {

	var items []string

	revDiff := lastRev - rev

	if revDiff > 500 {
		items = getBatchedItemsList(dir, rev, lastRev)
		return items
	}

	var eURL string

	switch dir {
	case "plugins":
		eURL = encodeURL(fmt.Sprintf(wpPluginChangelogURL, lastRev, revDiff))
	case "themes":
		eURL = encodeURL(fmt.Sprintf(wpThemeChangelogURL, lastRev, revDiff))
	}

	client := newClient(60, 500)

	req, err := http.NewRequest("GET", eURL, nil)
	if err != nil {
		fmt.Println(err)
	}

	req.Header.Set("User-Agent", userAgent)

	resp, err := client.Do(req)
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

	itemsGroups := regexUpdatedItems.FindAllStringSubmatch(bString, -1)

	found := make(map[string]bool)

	// Get the desired substring match and remove duplicates
	for _, item := range itemsGroups {

		if !found[item[1]] {
			found[item[1]] = true
			items = append(items, item[1])
		}

	}

	return items

}

func getBatchedItemsList(dir string, rev int, lastRev int) []string {

	var items []string

	client := newClient(60, 500)

	var eURL string

	var curRev int
	curRev = rev

	var batchSize int = 400

	fmt.Printf("Current Rev: %d\n", rev)
	fmt.Printf("Last Rev: %d\n", lastRev)

	found := make(map[string]bool)

	for curRev < lastRev {

		curRev += batchSize
		if curRev > lastRev {
			curRev = lastRev
		}

		switch dir {
		case "plugins":
			eURL = encodeURL(fmt.Sprintf(wpPluginChangelogURL, curRev, batchSize))
		case "themes":
			eURL = encodeURL(fmt.Sprintf(wpThemeChangelogURL, curRev, batchSize))
		}

		fmt.Println(eURL)

		req, err := http.NewRequest("GET", eURL, nil)
		if err != nil {
			fmt.Println(err)
		}

		req.Header.Set("User-Agent", userAgent)

		resp, err := client.Do(req)
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

		itemsGroups := regexUpdatedItems.FindAllStringSubmatch(bString, -1)

		// Get the desired substring match and remove duplicates
		for _, item := range itemsGroups {

			if !found[item[1]] {
				found[item[1]] = true
				items = append(items, item[1])
			}

		}

	}

	return items

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

		go func(name string, dir string) {
			getItem(name, dir)
			<-limiter
		}(v, dir)

	}

	bar.FinishPrint(fmt.Sprintf("Completed download of %d Items.", iCount))

	return nil

}

func getItem(item string, dir string) {

	c := &http.Client{
		Timeout: 60 * time.Second,
	}

	var eURL string

	switch dir {
	case "plugins":
		eURL = encodeURL(fmt.Sprintf(wpPluginDownloadURL, item))
	case "themes":
		eURL = encodeURL(fmt.Sprintf(wpThemeDownloadURL, item))
	}

	resp, err := c.Get(eURL)
	if err != nil {
		itemFetchFailure(item, dir)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {

		// Code 404 is acceptable, means the plugin/theme is no longer available.
		if resp.StatusCode == 404 {
			return
		}

		itemFetchFailure(item, dir)
		return

	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		itemFetchFailure(item, dir)
		return
	}

	err = extract(content, resp.ContentLength, item, dir)
	if err != nil {
		itemFetchFailure(item, dir)
		return
	}

}
