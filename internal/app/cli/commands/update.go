package commands

import (
	"log"
	"os"
	"runtime/pprof"

	"github.com/peterbooker/wpds2/internal/pkg/config"
	"github.com/peterbooker/wpds2/internal/pkg/context"
	"github.com/peterbooker/wpds2/internal/pkg/slurper"
	"github.com/peterbooker/wpds2/internal/pkg/stats"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.AddCommand(pluginsUpdateCmd)
	updateCmd.AddCommand(themesUpdateCmd)
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update files from the WordPress Directory.",
	Long:  `Update Plugin or Theme files from their WordPress Directory.`,
}

var pluginsUpdateCmd = &cobra.Command{
	Use:     "plugins",
	Short:   "Update Plugin files.",
	Long:    ``,
	Example: `wpds update plugins -c 250`,
	Run: func(cmd *cobra.Command, args []string) {

		if CPUProf != "" {
			f, err := os.Create(CPUProf)
			if err != nil {
				log.Fatal(err)
			}
			pprof.StartCPUProfile(f)
			defer pprof.StopCPUProfile()
		}

		if (C < 10) || (C > 10000) {
			log.Printf("Flag (concurrent-actions, c) out of permitted range (10-10000).\n")
			os.Exit(1)
		}

		log.Println("Updating Plugins...")

		// Get Config Details
		name := config.GetName()
		version := config.GetVersion()

		stats := stats.New()

		// Check if SVN is installed
		// Used if available, as it is more reliable than the HTTP API
		svn := slurper.CheckForSVN()

		ctx := &context.Context{
			Name:              name,
			Version:           version,
			ConcurrentActions: C,
			ExtensionType:     "plugins",
			FileType:          F,
			CurrentRevision:   0,
			LatestRevision:    0,
			SVN:               svn,
			Stats:             stats,
		}

		slurper.StartUpdate(ctx)

	},
}

var themesUpdateCmd = &cobra.Command{
	Use:     "themes",
	Short:   "Update Theme files.",
	Long:    ``,
	Example: `wpds update themes -c 250`,
	Run: func(cmd *cobra.Command, args []string) {

		if CPUProf != "" {
			f, err := os.Create(CPUProf)
			if err != nil {
				log.Fatal(err)
			}
			pprof.StartCPUProfile(f)
			defer pprof.StopCPUProfile()
		}

		if (C < 10) || (C > 1000) {
			log.Printf("Flag (concurrent-actions, c) out of permitted range (10-1000).\n")
			os.Exit(1)
		}

		log.Println("Updating Themes...")

		// Get Config Details
		name := config.GetName()
		version := config.GetVersion()

		stats := stats.New()

		// Check if SVN is installed
		// Used if available, as it is more reliable than the HTTP API
		svn := slurper.CheckForSVN()

		ctx := &context.Context{
			Name:              name,
			Version:           version,
			ConcurrentActions: C,
			ExtensionType:     "themes",
			FileType:          F,
			CurrentRevision:   0,
			LatestRevision:    0,
			SVN:               svn,
			Stats:             stats,
		}

		slurper.StartUpdate(ctx)

	},
}
