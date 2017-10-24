package storage

import (
	"fmt"
	"os"
)

// Controller implements the methods required to interact with the WordPress Directories
// Implemented via external HTTP API and local SVN client
type Controller interface {
	Write(path string) ([]byte, error)
	Read(path string) ([]byte, error)
}

// Init returns a connector used to communicate with the WordPress Directory SVN repositories.
// Implemented via an external HTTP API or local SVN client
func Init(cType string) Controller {

	switch cType {

	case "local":
		return &Local{}

	// FUTURE: Support more storage systems, e.g. AWS S3, Remote Filesystem, etc.

	default:
		// No supported storage found.
		fmt.Printf("The defined storage '%s' is not supported.", cType)
		os.Exit(1)
	}

	return nil

}
