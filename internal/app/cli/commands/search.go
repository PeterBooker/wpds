package commands

import (
	"log"
	"os"

	"github.com/peterbooker/wpds/internal/pkg/search"
	"github.com/peterbooker/wpds/internal/pkg/stats"
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
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		// Search Input
		input := args[0]

		// Get Working Directory
		wd, _ := os.Getwd()

		// Create new Stats
		stats := stats.New()

		// Setup Whitelist
		whitelist := []string{}

		// Setup Context
		ctx := &search.Context{
			ExtensionType:    cmd.Use,
			FileType:         F,
			ExtWhitelist:     whitelist,
			WorkingDirectory: wd,
			Stats:            stats,
		}

		log.Println("Search All Plugin files...")

		s := search.Setup(input, ctx)

		err := s.Run()
		if err != nil {
			log.Fatal(err)
		}

	},
}

var themesSearchCmd = &cobra.Command{
	Use:   "themes",
	Short: "Search the Theme files downloaded from the WordPress Directory.",
	Long:  `Search the Theme files downloaded from the WordPress Directory.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		// Search Input
		input := args[0]

		// Get Working Directory
		wd, _ := os.Getwd()

		// Create new Stats
		stats := stats.New()

		// Setup Whitelist
		whitelist := []string{}

		// Setup Context
		ctx := &search.Context{
			ExtensionType:    cmd.Use,
			FileType:         F,
			ExtWhitelist:     whitelist,
			WorkingDirectory: wd,
			Stats:            stats,
		}

		log.Println("Search All Theme files...")

		s := search.Setup(input, ctx)

		err := s.Run()
		if err != nil {
			log.Fatal(err)
		}

	},
}
