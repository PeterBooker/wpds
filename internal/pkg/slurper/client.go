package slurper

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"runtime"
	"time"
)

// NewClient returns a new HTTP client configured for large numbers of requests.
func NewClient(timeout int, max int) *http.Client {

	var netTransport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		MaxIdleConnsPerHost:   runtime.GOMAXPROCS(0) + 1,
	}

	var netClient = &http.Client{
		Timeout:   time.Second * time.Duration(timeout),
		Transport: netTransport,
	}

	return netClient

}

// NewRequest sets up and creates a new HTTP request to the given URL
func NewRequest(URL string, timeout int, concurrent int) ([]byte, error) {

	client := NewClient(timeout, concurrent)

	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		log.Println(err)
	}

	req.Header.Set("User-Agent", "wpds/0.5.0")

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Failed HTTP Request")
		os.Exit(1)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return []byte{}, fmt.Errorf("Invalid HTTP Status Code: %d", resp.StatusCode)

	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	return content, nil

}
