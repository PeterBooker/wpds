package commands

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(summaryCmd)
}

var summaryCmd = &cobra.Command{
	Use:   "summary",
	Short: "Print the version number of WPDS",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("View Search Summary...")
	},
}
