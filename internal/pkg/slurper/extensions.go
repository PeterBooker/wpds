package slurper

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"

	"github.com/peterbooker/wpds/internal/pkg/config"
	"github.com/peterbooker/wpds/internal/pkg/context"
	"github.com/peterbooker/wpds/internal/pkg/utils"
)

// fetchExtensions uses a list of extensions (themes or plugins) to download and extract their archives.
func fetchExtensions(extensions []string, ctx *context.Context) error {

	// Setup Progress Bar
	p := mpb.New(
		mpb.WithWidth(100),
	)
	bar := p.AddBar(int64(len(extensions)),
		mpb.PrependDecorators(
			decor.CountersNoUnit("%d / %d", 4, 0),
		),
		mpb.AppendDecorators(
			decor.Percentage(5, decor.DSyncSpace),
			decor.ETA(4, decor.DSyncSpace),
		),
	)

	limiter := make(chan struct{}, ctx.ConcurrentActions)

	// Make Plugins Dir ready for extracting plugins
	path := filepath.Join(wd, ctx.ExtensionType)
	err := utils.CreateDir(path)
	if err != nil {
		return err
	}

	// Use WaitGroup to ensure all Gorountines have finished downloading/extracting.
	var wg sync.WaitGroup

	// Look through extensions and start a Goroutine to download and extract the files.
	for _, name := range extensions {

		// Will block if more than max Goroutines already running.
		limiter <- struct{}{}

		wg.Add(1)

		go func(name string, ctx *context.Context, wg *sync.WaitGroup) {
			defer wg.Done()
			defer bar.Increment()

			getExtension(name, ctx, wg)
			<-limiter
		}(name, ctx, &wg)

	}

	wg.Wait()

	p.Wait()

	return nil

}

// getExtension fetches the relevant data for the extension e.g. All files, readme.txt, etc.
func getExtension(name string, ctx *context.Context, wg *sync.WaitGroup) {

	var file []byte
	var err error
	var size uint64

	switch ctx.FileType {
	case "all":

		// Gets the data of the archive file.
		file, err = getExtensionZip(name, ctx)
		if err != nil {
			extensionFailure(name, ctx)
			return
		}

		// Extracts the archive data to disk.
		size, err = ExtractZip(file, int64(len(file)), name, ctx)
		if err != nil {
			extensionFailure(name, ctx)
			return
		}

	case "readme":

		// Gets the data of the readme file.
		file, err = getExtensionReadme(name, ctx)
		if err != nil {
			extensionFailure(name, ctx)
			return
		}

		// Writes the readme file to disk.
		size, err = writeReadme(file, name, ctx)
		if err != nil {
			extensionFailure(name, ctx)
			return
		}

	}

	ctx.Stats.IncrementTotalExtensions()
	ctx.Stats.IncreaseTotalSize(size)

}

// getExtensionZip gets the extension archive.
func getExtensionZip(name string, ctx *context.Context) ([]byte, error) {

	client := NewClient(120, ctx.ConcurrentActions)

	var URL string
	var content []byte

	switch ctx.ExtensionType {
	case "plugins":
		URL = utils.EncodeURL(fmt.Sprintf(wpPluginDownloadURL, name))
	case "themes":
		URL = utils.EncodeURL(fmt.Sprintf(wpThemeDownloadURL, name))
	}

	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		log.Println(err)
		return content, err
	}

	// Dynamically set User-Agent from config
	req.Header.Set("User-Agent", config.GetName()+"/"+config.GetVersion())

	resp, err := client.Do(req)
	if err != nil {
		return content, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {

		// Code 404 is acceptable, means the plugin/theme is no longer available.
		if resp.StatusCode == 404 {
			ctx.Stats.IncrementTotalExtensionsClosed()
			return content, nil
		}

		log.Printf("Downloading the extension '%s' failed. Response code: %d\n", name, resp.StatusCode)

		return content, err

	}

	content, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return content, err
	}

	return content, err

}

// getExtensionReadme gets the extension readme.
func getExtensionReadme(name string, ctx *context.Context) ([]byte, error) {

	client := NewClient(120, ctx.ConcurrentActions)

	var URL string
	var content []byte

	switch ctx.ExtensionType {
	case "plugins":
		URL = utils.EncodeURL(fmt.Sprintf(wpPluginReadmeURL, name))
	case "themes":
		URL = utils.EncodeURL(fmt.Sprintf(wpThemeReadmeURL, name))
	}

	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		log.Println(err)
		return content, err
	}

	// Dynamically set User-Agent from config
	req.Header.Set("User-Agent", config.GetName()+"/"+config.GetVersion())

	resp, err := client.Do(req)
	if err != nil {
		return content, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {

		// Code 404 is acceptable, means the plugin/theme is no longer available.
		if resp.StatusCode == 404 {
			ctx.Stats.IncrementTotalExtensionsClosed()
			return content, nil
		}

		log.Printf("Downloading the extension '%s' failed. Response code: %d\n", name, resp.StatusCode)

		return content, err

	}

	content, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return content, err
	}

	return content, err

}

// writeReadme writes the readme file to disk.
func writeReadme(content []byte, name string, ctx *context.Context) (uint64, error) {

	base := filepath.Join(wd, ctx.ExtensionType, name)

	if utils.DirExists(base) {
		err := utils.RemoveDir(base)
		if err != nil {
			log.Printf("Cannot delete extension folder: %s\n", base)
		}
	}

	// Create base dir
	err := utils.CreateDir(base)
	if err != nil {
		log.Printf("Cannot create extension folder: %s\n", base)
	}

	path := filepath.Join(wd, ctx.ExtensionType, name, "readme.txt")

	if _, err := os.Stat(path); os.IsNotExist(err) {

		f, err := os.Create(path)
		defer f.Close()
		if err != nil {
			return 0, err
		}

	}

	f, err := os.OpenFile(path, os.O_RDWR, 0777)
	if err != nil {
		return 0, err
	}

	defer f.Close()

	_, err = f.Write(content)
	if err != nil {
		return 0, err
	}

	fi, err := f.Stat()
	if err != nil {
		return 0, err
	}

	return uint64(fi.Size()), nil

}
