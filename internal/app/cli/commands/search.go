package commands

import (
	"log"

	"github.com/peterbooker/wpds2/internal/pkg/search"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(searchCmd)
	searchCmd.AddCommand(pluginsSearchCmd)
	searchCmd.AddCommand(themesSearchCmd)
}

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search files downloaded from the WordPress Directory.",
	Long:  `Search files downloaded from the WordPress Directory.`,
}

var pluginsSearchCmd = &cobra.Command{
	Use:   "plugins",
	Short: "Search the Plugin files downloaded from the WordPress Directory.",
	Long:  `Search the Plugin files downloaded from the WordPress Directory.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		s := search.NewString(args[0])

		s = s

		log.Println("Search All Plugin files...")

	},
}

var themesSearchCmd = &cobra.Command{
	Use:   "themes",
	Short: "Search the Theme files downloaded from the WordPress Directory.",
	Long:  `Search the Theme files downloaded from the WordPress Directory.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		log.Printf("Search Input: %s\n", args[0])

		log.Println("Search All Theme files...")

	},
}
