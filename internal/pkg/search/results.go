package search

import (
	"fmt"

	"github.com/dustin/go-humanize"
)

func printSummary(ctx *Context) {

	switch ctx.ExtensionType {
	case "plugins":
		fmt.Printf("\n== Plugin Search Summary ==\n")
	case "themes":
		fmt.Printf("\n== Theme Search Summary ==\n")
	}

	fmt.Printf("Time Taken: %s\n", ctx.Stats.GetTimeTaken())

	totalExtensions := int64(ctx.Stats.GetTotalExtensionsFailed())
	totalFiles := int64(ctx.Stats.GetTotalFiles())
	totalMatches := int64(ctx.Stats.GetTotalFiles())

	fmt.Printf("Total Extensions: %s\n", humanize.Comma(totalExtensions))
	fmt.Printf("Total Files: %s\n", humanize.Comma(totalFiles))
	fmt.Printf("Total Matches: %s\n", humanize.Comma(totalMatches))

}
