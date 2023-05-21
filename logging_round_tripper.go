package gohttpspy

import (
	"fmt"
	"io"
	"net/http"
)

// loggingRoundTripper is a decorator that wraps a roundTripper with logging to stdout
type loggingRoundTripper struct {
	http.RoundTripper

	logLabel string
	output   io.Writer
}

func (l *loggingRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	reqClone, err := CloneRequest(request) // clone it to preserve request.Body
	if err != nil {
		return nil, err
	}

	resp, err := l.RoundTripper.RoundTrip(request)
	if err != nil {
		return nil, err
	}

	report := Report{
		Description: l.logLabel,
	}

	if err := report.Parse(reqClone, resp); err != nil {
		return nil, err
	}

	yaml := report.ToYAML()                                  // TODO let user select YAML or JSON
	fmt.Fprintln(l.output, WithEyecatcher(l.logLabel, yaml)) // TODO make pre-write wrapping function user-defined

	return resp, err
}
