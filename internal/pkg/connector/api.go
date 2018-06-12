package connector

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/peterbooker/wpds/internal/pkg/context"
)

const (
	wpAllPluginsListURL        = "http://plugins.svn.wordpress.org/"
	wpAllThemesListURL         = "http://themes.svn.wordpress.org/"
	wpLatestPluginsRevisionURL = "http://plugins.trac.wordpress.org/log/?format=changelog&stop_rev=HEAD"
	wpLatestThemesRevisionURL  = "http://themes.trac.wordpress.org/log/?format=changelog&stop_rev=HEAD"
	wpPluginChangelogURL       = "https://plugins.trac.wordpress.org/log/?verbose=on&mode=follow_copy&format=changelog&rev=%d&limit=%d"
	wpThemeChangelogURL        = "https://themes.trac.wordpress.org/log/?verbose=on&mode=follow_copy&format=changelog&rev=%d&limit=%d"
)

var (
	regexAPILatestRevision     = regexp.MustCompile(`\[(\d*)\]`)
	regexAPIFullExtensionsList = regexp.MustCompile(`.+?\>(\S+?)\/\<`)
)

// API implements the DirectoryConnector inferface.
// It uses an HTTP API to communicate with the WordPress Directory SVN Repositories.
type API struct {
	currentRevision int
	latestRevision  int
	extensions      []string
}

// GetLatestRevision gets the latest revision of the target directory.
func (api API) GetLatestRevision(ctx *context.Context) (int, error) {

	var URL string

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	URL = fmt.Sprintf("http://%s.trac.wordpress.org/log/?format=changelog&stop_rev=HEAD", ctx.ExtensionType)

	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return 0, err
	}

	// Set User Agent e.g. wpds/1.1.3
	userAgent := fmt.Sprintf("%s/%s", ctx.Name, ctx.Version)
	req.Header.Set("User-Agent", userAgent)

	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}

	if resp.StatusCode != 200 {
		return 0, fmt.Errorf("Invalid HTTP response code: %d", resp.StatusCode)
	}

	defer resp.Body.Close()
	bBytes, err := ioutil.ReadAll(resp.Body)
	bString := string(bBytes)

	revs := regexAPILatestRevision.FindAllStringSubmatch(bString, 1)

	api.latestRevision, err = strconv.Atoi(revs[0][1])
	if err != nil {
		return 0, err
	}

	return api.latestRevision, nil

}

// GetFullExtensionsList gets the fill list of all Extensions.
func (api *API) GetFullExtensionsList(ctx *context.Context) ([]string, error) {

	client := &http.Client{
		Timeout: 60 * time.Second,
	}

	URL := fmt.Sprintf("https://%s.svn.wordpress.org/", ctx.ExtensionType)

	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return api.extensions, err
	}

	// Set User Agent e.g. wpds/1.1.3
	userAgent := fmt.Sprintf("%s/%s", ctx.Name, ctx.Version)
	req.Header.Set("User-Agent", userAgent)

	resp, err := client.Do(req)
	if err != nil {
		return api.extensions, err
	}

	if resp.StatusCode != 200 {
		return api.extensions, fmt.Errorf("Invalid HTTP response code: %d", resp.StatusCode)
	}

	defer resp.Body.Close()
	bBytes, err := ioutil.ReadAll(resp.Body)
	bString := string(bBytes)

	matches := regexAPIFullExtensionsList.FindAllStringSubmatch(bString, -1)

	// Add all matches to extension list
	for _, extension := range matches {

		api.extensions = append(api.extensions, extension[1])

	}

	return api.extensions, nil

}

// GetUpdatedExtensionsList an updated list of Extensions.
func (api *API) GetUpdatedExtensionsList(ctx *context.Context) ([]string, error) {

	client := &http.Client{
		Timeout: 60 * time.Second,
	}

	URL := fmt.Sprintf("https://%s.trac.wordpress.org/log/?verbose=on&mode=follow_copy&format=changelog&rev=%d&limit=%d", ctx.ExtensionType, ctx.CurrentRevision, ctx.LatestRevision)

	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return api.extensions, err
	}

	// Set User Agent e.g. wpds/1.1.3
	userAgent := fmt.Sprintf("%s/%s", ctx.Name, ctx.Version)
	req.Header.Set("User-Agent", userAgent)

	resp, err := client.Do(req)
	if err != nil {
		return api.extensions, err
	}

	if resp.StatusCode != 200 {
		return api.extensions, fmt.Errorf("Invalid HTTP response code: %d", resp.StatusCode)
	}

	defer resp.Body.Close()
	bBytes, err := ioutil.ReadAll(resp.Body)
	bString := string(bBytes)

	matches := regexAPIFullExtensionsList.FindAllStringSubmatch(bString, 1)

	found := make(map[string]bool)

	// Get the desired substring match and remove duplicates
	for _, extension := range matches {

		if !found[extension[1]] {
			found[extension[1]] = true
			api.extensions = append(api.extensions, extension[1])
		}

	}

	return api.extensions, nil

}
