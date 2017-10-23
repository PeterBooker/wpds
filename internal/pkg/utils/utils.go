package utils

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/PuerkitoBio/goquery"
)

var (
	wd                string
	regexHTMLRevision = regexp.MustCompile(`[0-9]+`)
)

func init() {

	wd, _ = os.Getwd()

}

// parseItemListHTML parses HTML for the Lists of all Plugins or Themes. Identifies the latest revision
// number of the overall repository and a list containing the names of all items in that repository.
func ParseItemListHTML(r io.Reader) ([]string, int, error) {

	var items []string
	var revision int

	doc, err := goquery.NewDocumentFromReader(r)

	if err != nil {
		return []string{}, 0, nil
	}

	revText := doc.Find("h2").Text()

	rev := regexHTMLRevision.FindString(revText)

	revision, err = strconv.Atoi(rev)
	if err != nil {
		return []string{}, 0, err
	}

	doc.Find("ul").Each(func(i int, s *goquery.Selection) {

		doc.Find("a").Each(func(i int, s *goquery.Selection) {

			name := s.Text()
			name = strings.TrimSuffix(name, "/")
			items = append(items, name)

		})

	})

	return items, revision, nil

}

// IsConfirmationRequired checks if files already exist in the destination and the user
// needs to confirm beginning a fresh download.
func IsConfirmationRequired(dir string) bool {

	path := filepath.Join(wd, dir)

	// Check if directory exists
	if _, err := os.Stat(path); !os.IsNotExist(err) {

		// If exists, check if empty
		empty := IsDirEmpty(path)

		if empty == false {
			return true
		}

	}

	return false

}

// GetUserConfirmation prompts the user to confirm they are happy starting a fresh download
// to avoid accidentally polluting any current files.
func GetUserConfirmation() bool {

	// Ask for confirmation
	fmt.Println("Are you sure you want to begin a fresh download? (Y/N)")

	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		fmt.Println(err)
	}

	runes, _ := utf8.DecodeRuneInString(response)
	response = string(unicode.ToLower(runes))

	if string(response[0]) == "y" {
		return true
	}

	return false

}

// Properly encodes the URL for compatibility with special characters
// e.g. 新浪微博 and ЯндексФотки
func EncodeURL(rawURL string) string {

	u, _ := url.Parse(rawURL)

	URL := u.String()

	return URL

}
