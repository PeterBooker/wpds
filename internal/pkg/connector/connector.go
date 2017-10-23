package connector

import (
	"github.com/peterbooker/wpds2/internal/pkg/context"
)

// DirectoryConnector implements the methods required to interact with the WordPress Directories
// Implemented via external HTTP API and local SVN client
type DirectoryConnector interface {
	GetLatestRevision(ctx *context.Context) (int, error)
	GetFullExtensionsList(ctx *context.Context) ([]string, error)
	GetUpdatedExtensionsList(ctx *context.Context) ([]string, error)
}

// GetConnector returns a connector used to communicate with the WordPress Directory SVN repositories.
// Implemented via an external HTTP API or local SVN client
func GetConnector(ctx *context.Context) DirectoryConnector {

	var connector DirectoryConnector

	if ctx.SVN {

		// If SVN is available use it.
		connector = newSVN(ctx)

	} else {

		// If SVN is not available, fallback to the HTTP API
		connector = newAPI(ctx)

	}

	return connector

}
