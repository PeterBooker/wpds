package main

import (
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/url"
	"regexp"
	"strings"
)

// parseItemListHTML parses HTML for the Lists of all Plugins or Themes. Identifies the latest revision
// number of the overall repository and a list containing the names of all items in that repository.
func parseItemListHTML(r io.Reader) ([]string, string, error) {

	var items []string
	var revision string

	doc, err := goquery.NewDocumentFromReader(r)

	if err != nil {
		return []string{}, "", nil
	}

	revText := doc.Find("h2").Text()

	regex := regexp.MustCompile("[0-9]+")
	rev := regex.FindAllString(revText, 1)

	revision = rev[0]

	doc.Find("ul").Each(func(i int, s *goquery.Selection) {

		doc.Find("a").Each(func(i int, s *goquery.Selection) {

			name := s.Text()
			name = strings.TrimSuffix(name, "/")
			items = append(items, name)

		})

	})

	return items, revision, nil

}

// Properly encodes the URL for compatibility with special characters
// e.g. 新浪微博 and ЯндексФотки
func encodeURL(rawURL string) (string, error) {

	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	URL := u.String()

	return URL, nil

}
