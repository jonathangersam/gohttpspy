package gohttpspy

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/url"
	"testing"
)

func Test_loggingRoundTripper_RoundTrip(t *testing.T) {
	type fields struct {
		RoundTripper http.RoundTripper
		logLabel     string
	}

	tests := []struct {
		name         string
		fields       fields
		request      *http.Request
		wantResponse *http.Response
		wantLog      string
		wantErr      error
	}{
		{
			name: "success",
			fields: fields{
				RoundTripper: noopRoundTripper{},
				logLabel:     "my log",
			},
			request: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					Scheme: "http",
					Host:   "foo.com",
				},
			},
			wantResponse: &http.Response{
				Status:     "ok",
				StatusCode: http.StatusOK,
			},
			wantLog: `
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
  statusCode: 200
  headers: {}
  body: ""
--- [   F I N   ] --- my log --- [   F I N   ] ---


`,
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := bytes.NewBuffer(nil)

			lrt := &loggingRoundTripper{
				RoundTripper: tt.fields.RoundTripper,
				logLabel:     tt.fields.logLabel,
				output:       buf,
			}

			resp, err := lrt.RoundTrip(tt.request)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.wantResponse, resp)
			assert.Equal(t, tt.wantLog, buf.String())
		})
	}
}

// a no-op round tripper for mocks
type noopRoundTripper struct{}

func (noopRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	return &http.Response{
		Status:     "ok",
		StatusCode: http.StatusOK,
	}, nil
}
