package utils

import (
	"net/url"
)

// EncodeURL properly encodes the URL for compatibility with special characters
// e.g. 新浪微博 and ЯндексФотки
func EncodeURL(rawURL string) string {

	u, _ := url.Parse(rawURL)

	URL := u.String()

	return URL

}
