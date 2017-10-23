package connector

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"

	"github.com/peterbooker/wpds2/internal/pkg/context"
)

var (
	regexSVNUpdatedExtensionsList = regexp.MustCompile(`.{1,} \/(.+?)\/`)
	regexSVNFullExtensionsList    = regexp.MustCompile(`(.+?)\/`)
	regexSVNLatestRevision        = regexp.MustCompile(`r([0-9]+)`)
)

// SVN implements the Repository inferface.
// It uses a local SVN client to communicate with the WordPress Directory SVN Repositories.
type SVN struct{}

func newSVN(ctx *context.Context) *SVN {

	return &SVN{}

}

// GetLatestRevision gets the latest revision of the target directory
func (svn *SVN) GetLatestRevision(ctx *context.Context) (int, error) {

	URL := fmt.Sprintf("https://%s.svn.wordpress.org/", ctx.ExtensionType)

	args := []string{"log", "-v", "-q", URL, "-r", "HEAD"}

	out, _ := exec.Command("svn", args...).Output()

	itemsGroups := regexSVNLatestRevision.FindAllStringSubmatch(string(out), -1)

	revision, err := strconv.Atoi(itemsGroups[0][1])

	if err != nil {
		return 0, err
	}

	return revision, nil

}

// GetFullExtensionsList gets a full list
// TODO: Finish this, needs regex and output
func (svn *SVN) GetFullExtensionsList(ctx *context.Context) ([]string, error) {

	URL := fmt.Sprintf("https://%s.svn.wordpress.org/", ctx.ExtensionType)

	args := []string{"list", URL}

	out, _ := exec.Command("svn", args...).Output()

	groups := regexSVNFullExtensionsList.FindAllStringSubmatch(string(out), -1)

	var extensions []string

	// Add all matches to extension list
	for _, extension := range groups {

		extensions = append(extensions, extension[1])

	}

	return extensions, nil

}

// GetUpdatedExtensionsList gets an updated list
func (svn *SVN) GetUpdatedExtensionsList(ctx *context.Context) ([]string, error) {

	diff := fmt.Sprintf("%d:%d", ctx.CurrentRevision, ctx.LatestRevision)

	URL := fmt.Sprintf("https://%s.svn.wordpress.org/", ctx.ExtensionType)

	args := []string{"log", "-v", "-q", URL, "-r", diff}

	out, _ := exec.Command("svn", args...).Output()

	groups := regexSVNUpdatedExtensionsList.FindAllStringSubmatch(string(out), -1)

	var extensions []string

	found := make(map[string]bool)

	// Get the desired substring match and remove duplicates
	for _, extension := range groups {

		if !found[extension[1]] {
			found[extension[1]] = true
			extensions = append(extensions, extension[1])
		}

	}

	return extensions, nil

}
