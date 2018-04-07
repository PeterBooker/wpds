package slurper

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"sync"

	"github.com/peterbooker/wpds/internal/pkg/config"
	"github.com/peterbooker/wpds/internal/pkg/context"
	"github.com/peterbooker/wpds/internal/pkg/utils"
	"gopkg.in/cheggaaa/pb.v2"
)

const (
	barTemplate = `{{counters .}} {{bar .}} {{rtime .}} {{percent . }}`
)

func fetchExtensions(extensions []string, ctx *context.Context) error {

	bar := pb.ProgressBarTemplate(barTemplate).Start(len(extensions))

	limiter := make(chan struct{}, ctx.ConcurrentActions)

	// Make Plugins Dir ready for extracting plugins
	path := filepath.Join(wd, ctx.ExtensionType)
	err := utils.CreateDir(path)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup

	for _, name := range extensions {

		// Will block if more than max Goroutines already running.
		limiter <- struct{}{}

		wg.Add(1)

		go func(name string, ctx *context.Context, wg *sync.WaitGroup) {
			defer wg.Done()
			getExtension(name, ctx, wg)
			bar.Increment()
			<-limiter
		}(name, ctx, &wg)

	}

	wg.Wait()

	bar.Finish()

	return nil

}

func getExtension(name string, ctx *context.Context, wg *sync.WaitGroup) {

	client := NewClient(120, ctx.ConcurrentActions)

	var eURL string

	switch ctx.ExtensionType {
	case "plugins":
		eURL = utils.EncodeURL(fmt.Sprintf(wpPluginDownloadURL, name))
	case "themes":
		eURL = utils.EncodeURL(fmt.Sprintf(wpThemeDownloadURL, name))
	}

	req, err := http.NewRequest("GET", eURL, nil)
	if err != nil {
		log.Println(err)
	}

	// Dynamically set User-Agent from config
	req.Header.Set("User-Agent", config.GetName()+"/"+config.GetVersion())

	resp, err := client.Do(req)
	if err != nil {
		extensionFailure(name, ctx)
		return
	}

	if resp.StatusCode != 200 {

		// Code 404 is acceptable, means the plugin/theme is no longer available.
		if resp.StatusCode == 404 {
			ctx.Stats.IncrementTotalExtensionsClosed()
			return
		}

		log.Printf("Downloading the extension '%s' failed. Response code: %d\n", name, resp.StatusCode)

		extensionFailure(name, ctx)
		return

	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		extensionFailure(name, ctx)
		return
	}

	resp.Body.Close()

	size, err := ExtractZip(content, resp.ContentLength, name, ctx)
	if err != nil {
		extensionFailure(name, ctx)
		return
	}

	ctx.Stats.IncrementTotalExtensions()
	ctx.Stats.IncreaseTotalSize(size)

}

// GetExtension decides which extension data to fetch.
func GetExtension(name string, ctx *context.Context) {

	var file []byte
	var err error

	switch ctx.FileType {
	case "all":

		file, err = GetExtensionZip(name, ctx)
		if err != nil {
			// TODO: Handle Error
		}

		ExtractZip(file, int64(len(file)), "", ctx)

	case "readme":

		file, err = GetExtensionReadme(name, ctx)
		if err != nil {
			// TODO: Handle Error
		}

		//utils.WriteReadmeFile(file, name)

	}

}

// GetExtensionZip gets the extension archive.
func GetExtensionZip(name string, ctx *context.Context) ([]byte, error) {

	var URL string

	switch ctx.ExtensionType {
	case "plugins":
		URL = utils.EncodeURL(fmt.Sprintf(wpPluginDownloadURL, name))
	case "themes":
		URL = utils.EncodeURL(fmt.Sprintf(wpThemeDownloadURL, name))
	}

	response, err := NewRequest(URL, 60, ctx.ConcurrentActions)
	if err != nil {
		// TODO: Handle Failed Download
		return response, err
	}

	return response, nil

}

// GetExtensionReadme gets the extension readme.
func GetExtensionReadme(name string, ctx *context.Context) ([]byte, error) {

	var URL string

	switch ctx.ExtensionType {
	case "plugins":
		URL = utils.EncodeURL(fmt.Sprintf(wpPluginReadmeURL, name))
	case "themes":
		URL = utils.EncodeURL(fmt.Sprintf(wpThemeReadmeURL, name))
	}

	response, err := NewRequest(URL, 60, ctx.ConcurrentActions)
	if err != nil {
		// TODO: Handle Failed Download
		return response, err
	}

	return response, nil

}
