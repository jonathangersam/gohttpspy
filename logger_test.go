package gohttpspy

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"testing"
)

func TestLogger_Wrap(t *testing.T) {
	logger := &Logger{
		Output: os.Stdout,
		Label:  "my log",
	}

	wrappedClient := logger.Wrap(http.DefaultClient)

	_, ok := wrappedClient.Transport.(*loggingRoundTripper)
	assert.True(t, ok)
}
