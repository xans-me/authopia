package http

import (
	"net"
	"net/http"
	"time"
)

var netTransport = &http.Transport{
	Dial: (&net.Dialer{
		Timeout:   10 * time.Second,
		KeepAlive: 30 * time.Second,
	}).Dial,
	TLSHandshakeTimeout: 10 * time.Second,
	MaxIdleConns:        100,
	MaxIdleConnsPerHost: 100,
}
var netClient = &http.Client{
	Timeout:   30 * time.Second,
	Transport: netTransport,
}

// Request function to do http request, default 30 second timeout
func Request(request *http.Request) (*http.Response, error) {
	netClient.Timeout = 30 * time.Second
	resp, err := netClient.Do(request)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
