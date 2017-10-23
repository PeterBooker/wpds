package commands

import (
	"os"
	"runtime/pprof"

	"github.com/peterbooker/wpds2/internal/app/cli/log"
	"github.com/spf13/cobra"
)

var (
	CPUProf string
	MemProf string
	C       int
	V       bool
	F       string
	L       string
)

var rootCmd = &cobra.Command{
	Use:   "wpds",
	Short: "WPDS is a tool for downloading and searching the WordPress Plugin/Theme Directories.",
	Long:  `WPDS is a tool for downloading and searching the WordPress Plugin/Theme Directories.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {

		// Start Memory Profile
		if MemProf != "" {
			f, err := os.Create("start_" + MemProf)
			if err != nil {
				panic(err)
			}
			pprof.WriteHeapProfile(f)
			f.Close()
			return
		}

		// Setup the global logger using flag values
		// Runs before every command
		log.Setup(V, L)

	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {

		// End Memory Profile
		if MemProf != "" {
			f, err := os.Create("end_" + MemProf)
			if err != nil {
				panic(err)
			}
			pprof.WriteHeapProfile(f)
			f.Close()
			return
		}

	},
}

func init() {

	// Debug / Profiling Flags
	rootCmd.PersistentFlags().StringVar(&CPUProf, "cpuprof", "", "Filename of CPU profiling file.")
	rootCmd.PersistentFlags().StringVar(&MemProf, "memprof", "", "Filename of Memory profiling file.")

	// General Flags
	rootCmd.PersistentFlags().IntVarP(&C, "concurrent-actions", "c", 50, "Maximum number of concurrent actions (valid between 10-1000).")
	rootCmd.PersistentFlags().StringVarP(&L, "log", "l", "", "Destination of file to log output to.")
	rootCmd.PersistentFlags().BoolVarP(&V, "verbose", "v", false, "Verbose mode changed output from log to stdout.")
	rootCmd.PersistentFlags().StringVarP(&F, "files", "f", "all", "Extension files to download e.g. all or readme.")

}

// Execute processes CLI commands
func Execute() error {
	return rootCmd.Execute()
}
