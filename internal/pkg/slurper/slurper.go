package slurper

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/peterbooker/wpds/internal/pkg/connector"
	"github.com/peterbooker/wpds/internal/pkg/context"
	"github.com/peterbooker/wpds/internal/pkg/utils"
)

const (
	WPAllPluginsListURL        = "http://plugins.svn.wordpress.org/"
	WPAllThemesListURL         = "http://themes.svn.wordpress.org/"
	WPLatestPluginsRevisionURL = "http://plugins.trac.wordpress.org/log/?format=changelog&stop_rev=HEAD"
	WPLatestThemesRevisionURL  = "http://themes.trac.wordpress.org/log/?format=changelog&stop_rev=HEAD"
	WPPluginChangelogURL       = "https://plugins.trac.wordpress.org/log/?verbose=on&mode=follow_copy&format=changelog&rev=%d&limit=%d"
	WPThemeChangelogURL        = "https://themes.trac.wordpress.org/log/?verbose=on&mode=follow_copy&format=changelog&rev=%d&limit=%d"
	WPPluginDownloadURL        = "http://downloads.wordpress.org/plugin/%s.latest-stable.zip?nostats=1"
	WPThemeDownloadURL         = "http://downloads.wordpress.org/theme/%s.latest-stable.zip?nostats=1"
	WPPluginReadmeURL          = "https://plugins.svn.wordpress.org/%s/trunk/readme.txt"
	WPThemeReadmeURL           = "https://theme.svn.wordpress.org/%s/trunk/readme.txt"
	WPPluginInfoURL            = "https://api.wordpress.org/plugins/info/1.1/?action=plugin_information&request[slug]=%s&request[fields][active_installs]=1"
	WPThemeInfoURL             = "https://api.wordpress.org/themes/info/1.1/?action=plugin_information&request[slug]=%s&request[fields][active_installs]=1"
)

var (
	regexRevision     = regexp.MustCompile(`\[(\d*)\]`)
	regexHTMLRevision = regexp.MustCompile(`[0-9]+`)
	regexUpdatedItems = regexp.MustCompile(`\* ([^/A-Z ]+)`)
)

var (
	wd, _ = os.Getwd()
)

// StartUpdate begins the update 'plugin/theme' command.
// It begins by checking for an existing folder.
// TODO: Check for folder and .last-revision file as that is needed to update an existing slurp.
func StartUpdate(ctx *context.Context) {

	var fresh bool

	path := filepath.Join(wd, ctx.ExtensionType)

	if utils.DirExists(path) {

		// Dir exists, check if empty
		if utils.IsDirEmpty(path) {
			fresh = true
		} else {
			fresh = false
		}

	} else {
		// No existing slurp folder
		fresh = true
	}

	if fresh {

		// Begin fresh Directory Slurp
		err := newSlurp(ctx)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}

	} else {

		// Continue Existing Slurp Directory
		err := updateSlurp(ctx)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}

	}

	// Print Results of Command
	printResults(ctx)

}

func newSlurp(ctx *context.Context) error {

	var extensions []string
	var revision int
	var err error

	conn := connector.Init(ctx.Connector)

	extensions, err = conn.GetFullExtensionsList(ctx)
	if err != nil {
		return err
	}

	revision, err = conn.GetLatestRevision(ctx)
	if err != nil {
		return err
	}

	err = fetchExtensions(extensions, ctx)
	if err != nil {
		return err
	}

	err = setCurrentRevision(revision, ctx.ExtensionType)
	if err != nil {
		return err
	}

	return nil

}

func updateSlurp(ctx *context.Context) error {

	var extensions []string

	conn := connector.Init(ctx.Connector)

	currentRevision, err := getCurrentRevision(ctx.ExtensionType)
	if err != nil {
		return err
	}

	ctx.CurrentRevision = currentRevision

	latestRevision, err := conn.GetLatestRevision(ctx)
	if err != nil {
		return err
	}

	ctx.LatestRevision = latestRevision

	revisionDiff := latestRevision - currentRevision

	if revisionDiff <= 0 {
		fmt.Printf("No updates available. Revision: %d/%d.\n", currentRevision, latestRevision)
		os.Exit(1)
	}

	extensions, err = conn.GetUpdatedExtensionsList(ctx)
	if err != nil {
		return err
	}

	err = fetchExtensions(extensions, ctx)
	if err != nil {
		return err
	}

	err = setCurrentRevision(latestRevision, ctx.ExtensionType)
	if err != nil {
		return err
	}

	return nil

}
