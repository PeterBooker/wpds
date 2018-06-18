package utils

import (
	"io"
	"io/ioutil"
	"net/url"
)

// EncodeURL properly encodes the URL for compatibility with special characters
// e.g. 新浪微博 and ЯндексФотки
func EncodeURL(rawURL string) string {

	u, _ := url.Parse(rawURL)

	URL := u.String()

	return URL

}

// checkClose is used to check the return from Close in a defer statement.
func checkClose(c io.Closer, err *error) {
	cerr := c.Close()
	if *err == nil {
		*err = cerr
	}
}

// drainAndClose discards all data from rd and closes it.
func drainAndClose(rd io.ReadCloser, err *error) {
	if rd == nil {
		return
	}

	_, _ = io.Copy(ioutil.Discard, rd)
	cerr := rd.Close()
	if err != nil && *err == nil {
		*err = cerr
	}
}
