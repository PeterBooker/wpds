package context

import (
	"github.com/peterbooker/wpds2/internal/pkg/stats"
)

// Context contains the data required for Slurping
type Context struct {
	Name              string
	Version           string
	ConcurrentActions int
	ExtensionType     string
	FileType          string
	CurrentRevision   int
	LatestRevision    int
	SVN               bool
	Stats             *stats.Stats
}
