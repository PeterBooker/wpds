package main

import (
	"net/http"
	"net"
	"time"
)

func newClient(timeout int, max int) *http.Client {

	var netTransport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 10 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
		MaxIdleConns: max,
		MaxIdleConnsPerHost: max,
	}

	var netClient = &http.Client{
		Timeout: time.Second * time.Duration(timeout),
		Transport: netTransport,
	}

	return netClient

}