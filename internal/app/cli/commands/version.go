package commands

import (
	"fmt"

	"github.com/peterbooker/wpds/internal/pkg/config"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{

	Use:   "version",
	Short: "Print the version number of WPDS",
	Run: func(cmd *cobra.Command, args []string) {

		command := cmd.Use
		name := config.GetName()
		version := config.GetVersion()

		fmt.Printf("%s %s %s\n", name, command, version)

	},
}
