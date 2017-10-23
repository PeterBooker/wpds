package cli

import (
	"fmt"
	"os"

	"github.com/peterbooker/wpds/internal/app/cli/commands"
)

// Execute ...
func Execute() {

	if err := commands.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
