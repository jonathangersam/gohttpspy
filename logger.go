package gohttpspy

import (
	"io"
	"net/http"
	"os"
)

// Logger is a decorator that modifies a given http.Client to provide request/response cycle logging to stdout
type Logger struct {
	Output io.Writer
	Label  string
}

// Wrap returns a new http client that logs request/response cycle
func (l *Logger) Wrap(client *http.Client) *http.Client {
	// default log output to stdout
	output := l.Output
	if output == nil {
		output = os.Stdout
	}

	// wrap the existing client transport
	transportWithLogging := loggingRoundTripper{
		RoundTripper: client.Transport,
		logLabel:     l.Label,
		output:       l.Output,
	}

	client.Transport = &transportWithLogging

	return client
}
