package slurper

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
	"unicode/utf8"

	retry "github.com/giantswarm/retry-go"
	"github.com/peterbooker/wpds/internal/pkg/config"
	"github.com/peterbooker/wpds/internal/pkg/context"
	"github.com/peterbooker/wpds/internal/pkg/utils"
	pb "gopkg.in/cheggaaa/pb.v2"
)

const tmpl = `{{counters .}} {{bar . "[" "=" ">" "-" "]"}} {{rtime .}} {{percent .}}`

// fetchExtensions uses a list of extensions (themes or plugins) to download and extract their archives.
func fetchExtensions(extensions []string, ctx *context.Context) error {

	// Use WaitGroup to ensure all Gorountines have finished downloading/extracting.
	var wg sync.WaitGroup

	bar := pb.ProgressBarTemplate(tmpl).Start(len(extensions))

	limiter := make(chan struct{}, ctx.ConcurrentActions)

	// Make Plugins Dir ready for extracting plugins
	path := filepath.Join(wd, ctx.ExtensionType)
	err := utils.CreateDir(path)
	if err != nil {
		return err
	}

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

	bar.Finish()

	return nil

}

// getExtension fetches the relevant data for the extension e.g. All files, readme.txt, etc.
func getExtension(name string, ctx *context.Context, wg *sync.WaitGroup) {

	var data []byte
	var err error
	var size uint64

	// isExtensionNameValid?
	if !isValidName(name) {
		ctx.Stats.IncrementTotalExtensionsClosed()
		return
	}

	switch ctx.FileType {
	case "all":

		// Gets the data of the archive file.
		fetch := func() error {
			data, err = getExtensionZip(name, ctx)
			return err
		}
		err := retry.Do(fetch, retry.Timeout(600*time.Second), retry.MaxTries(3), retry.Sleep(5*time.Second))
		if err != nil {
			extensionFailure(name, ctx)
			return
		}

		// Received 404 response, not an error but we have no data so no more actions to take.
		if len(data) == 0 {
			ctx.Stats.IncrementTotalExtensionsClosed()
			return
		}

		// Extracts the archive data to disk.
		size, err = ExtractZip(data, int64(len(data)), name, ctx)
		if err != nil {
			extensionFailure(name, ctx)
			return
		}

	case "readme":

		// Gets the data of the readme file.
		fetch := func() error {
			data, err = getExtensionReadme(name, ctx)
			return err
		}
		err := retry.Do(fetch, retry.Timeout(600*time.Second), retry.MaxTries(3), retry.Sleep(5*time.Second))
		if err != nil {
			extensionFailure(name, ctx)
			return
		}

		// Received 404 response, not an error but we have no data so no more actions to take.
		if len(data) == 0 {
			ctx.Stats.IncrementTotalExtensionsClosed()
			return
		}

		// Writes the readme file to disk.
		size, err = writeReadme(data, name, ctx)
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

	client := NewClient(180, ctx.ConcurrentActions)

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

		// Code 404 is acceptable, it means the plugin/theme is no longer available.
		if resp.StatusCode == 404 {
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

	client := NewClient(180, ctx.ConcurrentActions)

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

// isValidName checks if the extension name is utf8 encoded, anything not will have been closed in the repository.
func isValidName(name string) bool {
	return utf8.Valid([]byte(name))
}
