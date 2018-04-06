package slurper

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
)

var (
	regexSVNRevisions = regexp.MustCompile(`.{1,} \/(.+?)\/`)
	regexSVNLatest    = regexp.MustCompile(`r([0-9]+)`)
)

// CheckForSVN checks if the SVN CLI tool is available.
func CheckForSVN() bool {

	_, err := exec.LookPath("svn")
	if err != nil {
		return false
	}

	return true

}

// getSVNUpdatedExtensions gets a list of extensions which were updated between the given revisions.
func getSVNUpdatedExtensions(cRev, lRev int, extType string) []string {

	diff := fmt.Sprintf("%d:%d", cRev, lRev)

	URL := fmt.Sprintf("https://%s.svn.wordpress.org/", extType)

	args := []string{"log", "-v", "-q", URL, "-r", diff}

	out, _ := exec.Command("svn", args...).Output()

	var items []string

	itemsGroups := regexSVNRevisions.FindAllStringSubmatch(string(out), -1)

	found := make(map[string]bool)

	// Get the desired substring match and remove duplicates
	for _, item := range itemsGroups {

		if !found[item[1]] {
			found[item[1]] = true
			items = append(items, item[1])
		}

	}

	return items

}

// getSVNLatestRevision gets the latest revision from the target repository.
func getSVNLatestRevision(extType string) (int, error) {

	URL := fmt.Sprintf("https://%s.svn.wordpress.org/", extType)

	args := []string{"log", "-v", "-q", URL, "-r", "HEAD"}

	out, _ := exec.Command("svn", args...).Output()

	itemsGroups := regexSVNLatest.FindAllStringSubmatch(string(out), -1)

	revision, err := strconv.Atoi(itemsGroups[0][1])

	if err != nil {
		return 0, err
	}

	return revision, nil

}
