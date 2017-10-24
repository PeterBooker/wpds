package connector

import (
	"fmt"
	"os"
	"strings"

	"github.com/peterbooker/wpds/internal/pkg/context"
)

// DirectoryConnector implements the methods required to interact with the WordPress Directories
// Implemented via external HTTP API and local SVN client
type DirectoryConnector interface {
	GetLatestRevision(ctx *context.Context) (int, error)
	GetFullExtensionsList(ctx *context.Context) ([]string, error)
	GetUpdatedExtensionsList(ctx *context.Context) ([]string, error)
}

// Init returns a connector used to communicate with the WordPress Directory SVN repositories.
// Implemented via an external HTTP API or local SVN client
func Init(cType string) DirectoryConnector {

	switch cType {

	case "svn":
		return &SVN{}

	case "api":
		return &API{}

	default:
		// No supported storage found.
		fmt.Printf("The defined connector '%s' is not supported.", strings.ToUpper(cType))
		os.Exit(1)

	}

	return nil

}
