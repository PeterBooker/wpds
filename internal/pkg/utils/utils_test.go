package utils

import (
	"net/url"
	"testing"
)

func TestEncodeURL(t *testing.T) {

	cases := map[string]struct{ text, expected string }{
		"Plain":   {"wp-seo", "wp-seo"},
		"Special": {"µmint", "%C2%B5mint"},
		"Symbols": {"★-wpsymbols-★", "%E2%98%85-wpsymbols-%E2%98%85"},
		"Russian": {"ЯндексФотки", "%D0%AF%D0%BD%D0%B4%D0%B5%D0%BA%D1%81%D0%A4%D0%BE%D1%82%D0%BA%D0%B8"},
		"Arabic":  {"لينوكس-ويكى", "%D9%84%D9%8A%D9%86%D9%88%D9%83%D8%B3-%D9%88%D9%8A%D9%83%D9%89"},
		"Chinese": {"豆瓣秀-for-wordpress", "%E8%B1%86%E7%93%A3%E7%A7%80-for-wordpress"},
	}

	for k, v := range cases {
		u, _ := url.Parse(v.text)
		actual := u.String()
		if actual != v.expected {
			t.Errorf("%s - Raw: %s Encoded: %s Expected: %s\n", k, v.text, u.String(), v.expected)
		}

	}

}
