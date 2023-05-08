package gohttpspy

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"gopkg.in/yaml.v3"
)

// Config defines fields to be used for generating valid mmock configs.
// This is a helper for writers of e2e tests
type Config struct {
	content Content
}

func NewConfig(desc string, req *http.Request, resp *http.Response) *Config {
	// parse request body
	var reqBody []byte
	if req.Body != nil {
		var err error
		reqBody, err = io.ReadAll(req.Body)
		if err != nil {
			panic(err)
		}

		req.Body = ioutil.NopCloser(bytes.NewReader(reqBody))
	}

	// parse response body
	var respBody []byte
	if resp.Body != nil {
		ps, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		resp.Body = ioutil.NopCloser(bytes.NewReader(ps))

		// indent
		buf := bytes.NewBuffer(nil)
		if err := json.Indent(buf, ps, "", "  "); err != nil {
			panic(err)
		}

		respBody, err = io.ReadAll(buf)
		if err != nil {
			panic(err)
		}
	}

	// fill mmock struct that will be used to generate YAML
	config := Config{
		content: Content{
			Description: desc,
			Request: Request{
				Scheme: "http",
				//Path:                  req.URL.String(), // TODO
				Path:                  req.URL.Path,
				Method:                req.Method,
				QueryStringParameters: QueryStringParameters(req.URL.Query()),
				Headers:               HttpHeaders(req.Header),
				Body:                  string(reqBody),
			},
			Response: Response{
				StatusCode: resp.StatusCode,
				Headers:    HttpHeaders(resp.Header),
				Body:       string(respBody),
			},
		},
	}

	return &config
}

func (c *Config) ToYAML() string {
	//ps, err := yaml.Marshal(m.content)
	buf := bytes.NewBuffer([]byte{})
	enc := yaml.NewEncoder(buf)
	enc.SetIndent(2)
	err := enc.Encode(c.content)
	if err != nil {
		log.Fatal(err)
	}

	ps, err := io.ReadAll(buf)
	if err != nil {
		log.Fatal(err)
	}

	return string(ps)
}

type Content struct {
	Description string
	Request     Request
	Response    Response
}

type Request struct {
	Scheme                string
	Path                  string
	Method                string
	QueryStringParameters QueryStringParameters `yaml:"queryStringParameters"`
	Headers               HttpHeaders
	Body                  string
}

type QueryStringParameters map[string][]string // list of qParamValues by qParamKey

type Response struct {
	StatusCode int `yaml:"statusCode"`
	Headers    HttpHeaders
	Body       string
}

type HttpHeaders map[string][]string // list of headerValues by headerKey
