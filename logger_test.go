package gohttpspy

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
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

func TestLogger_Integration(t *testing.T) {
	// setup mock
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		hdr := w.Header()
		hdr.Set("Date", "xyz") // overwrite header Date because its non-deterministic
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(`{"foo":"bar"}`))
		assert.NoError(t, err)
	}))
	defer testServer.Close()

	// create SUT
	buf := bytes.NewBuffer(nil)
	logger := &Logger{
		Output: buf,
		Label:  "my log",
	}
	wrappedClient := logger.Wrap(http.DefaultClient)

	// execute test
	resp, err := wrappedClient.Get(testServer.URL)
	assert.NoError(t, err)
	defer resp.Body.Close()

	// assert
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	ps, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.JSONEq(t, `{"foo": "bar"}`, string(ps))

	wantLogs := `
--- [ S T A R T ] --- my log --- [ S T A R T ] ---
description: my log
request:
  scheme: http
  path: ""
  method: GET
  queryStringParameters: {}
  headers: {}
  body: ""
response:
  statusCode: 400
  headers:
    Content-Length:
      - "13"
    Content-Type:
      - text/plain; charset=utf-8
    Date:
      - xyz
  body: |-
    {
      "foo": "bar"
    }
--- [   F I N   ] --- my log --- [   F I N   ] ---


`

	assert.Equal(t, wantLogs, buf.String())
}
