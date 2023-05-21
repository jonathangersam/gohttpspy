package gohttpspy

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/url"
	"testing"
)

func TestReport_Parse(t *testing.T) {
	tests := []struct {
		name     string
		request  *http.Request
		response *http.Response
		want     Content
		wantErr  error
	}{
		{
			name: "normal case",
			request: &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					Scheme: "http",
					Path:   "foo.com",
				},
			},
			response: &http.Response{
				Status:     "ok",
				StatusCode: http.StatusOK,
			},
			want: Content{
				Description: "my description",
				Request: Request{
					Scheme:                "http",
					Path:                  "foo.com",
					Method:                "GET",
					QueryStringParameters: QueryStringParameters{},
					Headers:               nil,
					Body:                  "",
				},
				Response: Response{
					StatusCode: http.StatusOK,
					Headers:    nil,
					Body:       "",
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			report := &Report{
				Description: "my description",
			}

			err := report.Parse(tt.request, tt.response)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, report.content, "Parse()")
		})
	}
}

func TestReport_ToYAML(t *testing.T) {
	tests := []struct {
		name    string
		content Content
		want    string
	}{
		{
			name: "normal case",
			content: Content{
				Description: "my description",
				Request: Request{
					Scheme:                "https",
					Path:                  "foo.com",
					Method:                "GET",
					QueryStringParameters: nil,
					Headers:               nil,
					Body:                  "",
				},
				Response: Response{
					StatusCode: http.StatusOK,
					Headers:    nil,
					Body:       "",
				},
			},
			want: `
description: my description
request:
  scheme: https
  path: foo.com
  method: GET
  queryStringParameters: {}
  headers: {}
  body: ""
response:
  statusCode: 200
  headers: {}
  body: ""`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			report := &Report{
				content: tt.content,
			}
			assert.YAMLEq(t, tt.want, report.ToYAML(), "ToYAML()")
		})
	}
}
