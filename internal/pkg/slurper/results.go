package slurper

import (
	"fmt"

	"github.com/dustin/go-humanize"
	"github.com/peterbooker/wpds/internal/pkg/context"
)

func printResults(ctx *context.Context) {

	fmt.Printf("\n== Command Results ==\n")
	fmt.Printf("Time Taken: %s\n", ctx.Stats.GetTimeTaken())

	switch ctx.ExtensionType {
	case "plugins":
		pluginResults(ctx)
	case "themes":
		themeResults(ctx)
	}

	failedDownloads := int64(ctx.Stats.GetTotalExtensionsFailed())
	totalFiles := int64(ctx.Stats.GetTotalFiles())
	totalFileSize := ctx.Stats.GetTotalSize()

	fmt.Printf("Failed Downloads: %s\n", humanize.Comma(failedDownloads))
	fmt.Printf("Total Files: %s\n", humanize.Comma(totalFiles))
	fmt.Printf("Total Disk Size: %s\n", humanize.Bytes(totalFileSize))

}

func pluginResults(ctx *context.Context) {

	totalPlugins := int64(ctx.Stats.GetTotalExtensions())
	closedPlugins := int64(ctx.Stats.GetTotalExtensionsClosed())

	fmt.Printf("Total Plugins: %s\n", humanize.Comma(totalPlugins))
	fmt.Printf("Closed/Disabled Plugins: %s\n", humanize.Comma(closedPlugins))

}

func themeResults(ctx *context.Context) {

	totalThemes := int64(ctx.Stats.GetTotalExtensions())
	closedThemes := int64(ctx.Stats.GetTotalExtensionsClosed())

	fmt.Printf("Total Themes: %s\n", humanize.Comma(totalThemes))
	fmt.Printf("Unapproved/Closed Themes: %s\n", humanize.Comma(closedThemes))

}
