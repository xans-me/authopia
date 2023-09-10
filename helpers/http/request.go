package http

import (
	"github.com/felixge/httpsnoop"
	log "github.com/sirupsen/logrus"
	"net"
	"net/http"
	"strings"
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

func WithLogger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		m := httpsnoop.CaptureMetrics(handler, writer, request)
		log.Printf("http[%d]-- %s -- %s\n", m.Code, m.Duration, request.URL.Path)
	})
}

var allowedHeaders = map[string]struct{}{
	"x-request-id": {},
}

func IsHeaderAllowed(s string) (string, bool) {
	// check if allowedHeaders contain the header
	if _, isAllowed := allowedHeaders[s]; isAllowed {
		// send uppercase header
		return strings.ToUpper(s), true
	}
	// if not in allowed header, don't send the header
	return s, false
}
