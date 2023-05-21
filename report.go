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

type Report struct {
	Description string
	content     Content
}

func (r *Report) Parse(req *http.Request, resp *http.Response) error {
	// parse request body
	var reqBody []byte
	if req.Body != nil {
		var err error

		if req.Body != nil {
			reqBody, err = io.ReadAll(req.Body)
			if err != nil {
				return err
			}

			req.Body = ioutil.NopCloser(bytes.NewReader(reqBody))
		}
	}

	// parse response body
	var respBody []byte
	if resp.Body != nil {
		if resp.Body != nil {
			ps, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			resp.Body = ioutil.NopCloser(bytes.NewReader(ps))

			// indent
			buf := bytes.NewBuffer(nil)
			if err := json.Indent(buf, ps, "", "  "); err != nil {
				return err
			}

			respBody, err = io.ReadAll(buf)
			if err != nil {
				return err
			}
		}
	}

	// fill struct that will be used to generate YAML
	r.content = Content{
		Description: r.Description,
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
	}

	return nil
}

func (r *Report) ToYAML() string {
	//ps, err := yaml.Marshal(m.content)
	buf := bytes.NewBuffer([]byte{})
	enc := yaml.NewEncoder(buf)
	enc.SetIndent(2)
	err := enc.Encode(r.content)
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
