package connector

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"

	"github.com/peterbooker/wpds/internal/pkg/context"
)

var (
	regexSVNUpdatedExtensionsList = regexp.MustCompile(`.{1,} \/(.+?)\/`)
	regexSVNFullExtensionsList    = regexp.MustCompile(`(.+?)\/`)
	regexSVNLatestRevision        = regexp.MustCompile(`r([0-9]+)`)
)

// SVN implements the Repository inferface.
// It uses a local SVN client to communicate with the WordPress Directory SVN Repositories.
type SVN struct {
	currentRevision int
	latestRevision  int
	extensions      []string
}

func newSVN(ctx *context.Context) *SVN {

	return &SVN{}

}

// GetLatestRevision gets the latest revision of the target directory.
func (svn *SVN) GetLatestRevision(ctx *context.Context) (int, error) {

	URL := fmt.Sprintf("https://%s.svn.wordpress.org/", ctx.ExtensionType)

	args := []string{"log", "-v", "-q", URL, "-r", "HEAD"}

	out, _ := exec.Command("svn", args...).Output()

	matches := regexSVNLatestRevision.FindAllStringSubmatch(string(out), -1)

	var err error

	svn.latestRevision, err = strconv.Atoi(matches[0][1])
	if err != nil {
		return 0, err
	}

	return svn.latestRevision, nil

}

// GetFullExtensionsList gets the fill list of all Extensions.
func (svn *SVN) GetFullExtensionsList(ctx *context.Context) ([]string, error) {

	URL := fmt.Sprintf("https://%s.svn.wordpress.org/", ctx.ExtensionType)

	args := []string{"list", URL}

	out, _ := exec.Command("svn", args...).Output()

	matches := regexSVNFullExtensionsList.FindAllStringSubmatch(string(out), -1)

	// Add all matches to extension list
	for _, extension := range matches {

		svn.extensions = append(svn.extensions, extension[1])

	}

	return svn.extensions, nil

}

// GetUpdatedExtensionsList gets an updated list of Extensions.
func (svn *SVN) GetUpdatedExtensionsList(ctx *context.Context) ([]string, error) {

	diff := fmt.Sprintf("%d:%d", ctx.CurrentRevision, ctx.LatestRevision)

	URL := fmt.Sprintf("https://%s.svn.wordpress.org/", ctx.ExtensionType)

	args := []string{"log", "-v", "-q", URL, "-r", diff}

	out, _ := exec.Command("svn", args...).Output()

	groups := regexSVNUpdatedExtensionsList.FindAllStringSubmatch(string(out), -1)

	found := make(map[string]bool)

	// Get the desired substring match and remove duplicates
	for _, extension := range groups {

		if !found[extension[1]] {
			found[extension[1]] = true
			svn.extensions = append(svn.extensions, extension[1])
		}

	}

	return svn.extensions, nil

}
