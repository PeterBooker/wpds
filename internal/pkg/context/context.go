package context

import (
	"net/http"

	"github.com/peterbooker/wpds/internal/pkg/stats"
)

// Context contains the data required for Slurping
type Context struct {
	Name              string
	Version           string
	ConcurrentActions int
	ExtensionType     string
	FileType          string
	Connector         string
	CurrentRevision   int
	LatestRevision    int
	WorkingDirectory  string
	FailedList        []string
	Stats             *stats.Stats
	Client            *http.Client
}

// SearchContext contains data required for Searching
type SearchContext struct {
	ExtensionType    string
	FileType         string
	ExtWhitelist     []string
	WorkingDirectory string
	Stats            *stats.Stats
}
