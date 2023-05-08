package gohttpspy

import (
	"fmt"
	"net/http"
)

// NewLoggingHttpClient is a decorator that modifies a given http.Client to provide request/response cycle logging to stdout
func NewLoggingHttpClient(client *http.Client, logLabel string) *http.Client {
	client.Transport = newLoggingRoundTripper(client.Transport, logLabel)
	return client
}

// loggingRoundTripper is a decorator that wraps a roundTripper with logging to stdout
type loggingRoundTripper struct {
	logLabel string
	http.RoundTripper
}

func newLoggingRoundTripper(wrapped http.RoundTripper, logLabel string) http.RoundTripper {
	return &loggingRoundTripper{
		logLabel:     logLabel,
		RoundTripper: wrapped,
	}
}

func (l *loggingRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	reqClone := CloneRequest(request) // clone it to preserve request.Body

	resp, err := l.RoundTripper.RoundTrip(request)

	yaml := NewConfig(l.logLabel, reqClone, resp).ToYAML()
	fmt.Println(WithEyecatcher(l.logLabel, yaml))

	return resp, err
}
