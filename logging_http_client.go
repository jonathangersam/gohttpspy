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
	reqClone, err := CloneRequest(request) // clone it to preserve request.Body
	if err != nil {
		return nil, err
	}

	resp, err := l.RoundTripper.RoundTrip(request)

	config, err := NewConfig(l.logLabel, reqClone, resp)
	if err != nil {
		return nil, err
	}

	yaml := config.ToYAML()
	fmt.Println(WithEyecatcher(l.logLabel, yaml))

	return resp, err
}
