package main

import (
	"fmt"
	"os"
	"regexp"
	"runtime"
	"runtime/pprof"
	"sort"
	"time"

	"github.com/peterbooker/wpds/search"
	"github.com/urfave/cli"
)

const (
	version = "0.3.0"
	userAgent = "wpds/" + version
)

const (
	wpAllPluginsListURL        = "http://plugins.svn.wordpress.org/"
	wpAllThemesListURL         = "http://themes.svn.wordpress.org/"
	wpLatestPluginsRevisionURL = "http://plugins.trac.wordpress.org/log/?format=changelog&stop_rev=HEAD"
	wpLatestThemesRevisionURL  = "http://themes.trac.wordpress.org/log/?format=changelog&stop_rev=HEAD"
	wpPluginChangelogURL       = "https://plugins.trac.wordpress.org/log/?verbose=on&mode=follow_copy&format=changelog&rev=%d&limit=%d"
	wpThemeChangelogURL        = "https://themes.trac.wordpress.org/log/?verbose=on&mode=follow_copy&format=changelog&rev=%d&limit=%d"
	wpPluginDownloadURL        = "http://downloads.wordpress.org/plugin/%s.latest-stable.zip?nostats=1"
	wpThemeDownloadURL         = "http://downloads.wordpress.org/theme/%s.latest-stable.zip?nostats=1"
	wpPluginReadmeURL          = "https://plugins.svn.wordpress.org/%s/trunk/readme.txt"
	wpThemeReadmeURL           = "https://theme.svn.wordpress.org/%s/trunk/readme.txt"
	wpPluginInfoURL            = "https://api.wordpress.org/plugins/info/1.1/?action=plugin_information&request[slug]=%s&request[fields][active_installs]=1"
	wpThemeInfoURL             = "https://api.wordpress.org/themes/info/1.1/?action=plugin_information&request[slug]=%s&request[fields][active_installs]=1"
)

var (
	regexRevision     = regexp.MustCompile(`\[(\d*)\]`)
	regexHTMLRevision = regexp.MustCompile(`[0-9]+`)
	regexUpdatedItems = regexp.MustCompile(`\* ([^/A-Z ]+)`)
)

var (
	wd string
)

func main() {

	app := cli.NewApp()
	app.Name = "WPDS"
	app.Usage = "WPDS is a CLI tool for downloading and searching the WordPress Plugin/Theme Directories."
	app.Version = version
	app.Compiled = time.Now()
	app.Authors = []cli.Author{
		{
			Name:  "Peter Booker",
			Email: "mail@peterbooker.com",
		},
	}

	wd, _ = os.Getwd()

	// Default Action - No Command
	app.Action = func(c *cli.Context) error {

		fmt.Printf("Name: %s Version: %s\n", c.App.Name, c.App.Version)
		fmt.Printf("Description: %s\n", c.App.Usage)
		fmt.Println("Type \"wpds -help\" for more information.")

		return nil

	}

	// Support for inbuilt performance monitoring
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "cpuprofile",
			Value: "",
			Usage: "Write CPU profiling to file.",
		},
		cli.StringFlag{
			Name:  "memprofile",
			Value: "",
			Usage: "Write Memory profiling to file.",
		},
	}

	// Setup Commands and Sub Commands
	app.Commands = []cli.Command{
		{
			Name:    "download",
			Aliases: []string{"d"},
			Usage:   "Download and update all WordPress Plugins or Themes.",
			Subcommands: []cli.Command{
				{
					Name:  "plugins",
					Usage: "Download all WordPress Plugins.",
					Flags: []cli.Flag{
						cli.IntFlag{
							Name:  "limit, l",
							Value: 10,
							Usage: "Number of simultaneous downloads.",
						},
						cli.StringFlag{
							Name:  "type, t",
							Value: "all",
							Usage: "Type of files to download.",
						},
					},
					Action: func(c *cli.Context) error {

						if isConfirmationRequired("plugins") {

							confirm := getUserConfirmation()
							if !confirm {
								os.Exit(0)
							}

						}

						getAllItems(c, "plugins")

						return nil

					},
				},
				{
					Name:  "themes",
					Usage: "Download all WordPress Themes.",
					Flags: []cli.Flag{
						cli.IntFlag{
							Name:  "limit, l",
							Value: 10,
							Usage: "Number of simultaneous downloads.",
						},
						cli.StringFlag{
							Name:  "type, t",
							Value: "all",
							Usage: "Type of files to download.",
						},
					},
					Action: func(c *cli.Context) error {

						if isConfirmationRequired("themes") {

							confirm := getUserConfirmation()
							if !confirm {
								os.Exit(0)
							}

						}

						getAllItems(c, "themes")

						return nil

					},
				},
			},
			Before: func(c *cli.Context) error {

				if cprof := c.GlobalString("cpuprofile"); cprof != "" {

					go func() {

						time.Sleep(60 * time.Second)

						f, err := os.Create(cprof)
						if err != nil {
							panic(err)
						}
						pprof.StartCPUProfile(f)
						defer pprof.StopCPUProfile()

					}()

				}

				if mprof := c.GlobalString("memprofile"); mprof != "" {

					go func() {

						time.Sleep(60 * time.Second)

						f, err := os.Create(mprof)
						if err != nil {
							panic(err)
						}
						runtime.GC()
						if err := pprof.WriteHeapProfile(f); err != nil {
							panic(err)
						}
						f.Close()

					}()

				}

				started := time.Now()

				c.App.Metadata["started"] = started

				return nil
			},
			After: func(c *cli.Context) error {

				elapsed := time.Since(c.App.Metadata["started"].(time.Time))
				fmt.Printf("Command took %s\n", elapsed)

				return nil
			},
			Action: func(c *cli.Context) error {
				return nil
			},
		},
		{
			Name:    "update",
			Aliases: []string{"u"},
			Usage:   "Update all WordPress Plugins or Themes.",
			Subcommands: []cli.Command{
				{
					Name:  "plugins",
					Usage: "Update all WordPress Plugins.",
					Flags: []cli.Flag{
						cli.IntFlag{
							Name:  "limit, l",
							Value: 10,
							Usage: "Number of simultaneous downloads.",
						},
						cli.StringFlag{
							Name:  "type, t",
							Value: "all",
							Usage: "Type of files to download.",
						},
					},
					Action: func(c *cli.Context) error {

						rev, err := getCurrentRevision("plugins")
						if err != nil {
							return cli.NewExitError(err.Error(), 1)
						}

						getUpdatedItems(c, "plugins", rev)

						return nil

					},
				},
				{
					Name:  "themes",
					Usage: "Update all WordPress Themes.",
					Flags: []cli.Flag{
						cli.IntFlag{
							Name:  "limit, l",
							Value: 10,
							Usage: "Number of simultaneous downloads.",
						},
						cli.StringFlag{
							Name:  "type, t",
							Value: "all",
							Usage: "Type of files to download.",
						},
					},
					Action: func(c *cli.Context) error {

						rev, err := getCurrentRevision("themes")
						if err != nil {
							return cli.NewExitError(err.Error(), 1)
						}

						getUpdatedItems(c, "themes", rev)

						return nil

					},
				},
			},
			Before: func(c *cli.Context) error {

				if cprof := c.GlobalString("cpuprofile"); cprof != "" {

					f, err := os.Create(cprof)
					if err != nil {
						panic(err)
					}
					pprof.StartCPUProfile(f)
					defer pprof.StopCPUProfile()

				}

				if mprof := c.GlobalString("memprofile"); mprof != "" {

					f, err := os.Create(mprof)
					if err != nil {
						panic(err)
					}
					runtime.GC()
					if err := pprof.WriteHeapProfile(f); err != nil {
						panic(err)
					}
					f.Close()

				}

				started := time.Now()

				c.App.Metadata["started"] = started

				return nil
			},
			After: func(c *cli.Context) error {

				elapsed := time.Since(c.App.Metadata["started"].(time.Time))
				fmt.Printf("Command took %s\n", elapsed)

				return nil
			},
			Action: func(c *cli.Context) error {
				return nil
			},
		},
		{
			Name:    "search",
			Aliases: []string{"s"},
			Usage:   "Search all WordPress Plugins.",
			Action: func(c *cli.Context) error {
				return nil
			},
			Subcommands: []cli.Command{
				{
					Name:  "plugins",
					Usage: "Search all WordPress Plugins.",
					Action: func(c *cli.Context) error {

						pattern := c.Args().Get(0)
						if pattern == "" {
							return cli.NewExitError("Please specify a search pattern.", 20)
						}

						search.NewStringSearch(pattern, "plugins")

						//results := startSearch(pattern)

						//outputResults(results, pattern, "stdout")

						return nil
					},
				},
				{
					Name:  "themes",
					Usage: "Search all WordPress Themes.",
					Action: func(c *cli.Context) error {

						pattern := c.Args().Get(0)
						if pattern == "" {
							return cli.NewExitError("Please specify a search pattern.", 20)
						}

						//results := startSearch(pattern)

						//outputResults(results, pattern, "stdout")

						return nil
					},
				},
			},
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	app.Run(os.Args)

}
