package gohttpspy

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

func CloneRequest(r *http.Request) *http.Request {
	clone := *r

	if r.Body != nil {
		ps, err := io.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}

		r.Body = ioutil.NopCloser(bytes.NewReader(ps))
		clone.Body = ioutil.NopCloser(bytes.NewReader(ps))
	}

	return &clone
}

func WithEyecatcher(label, s string) string {
	//var buffer strings.Builder
	buffer := new(strings.Builder)

	fmt.Fprintf(buffer, "\n--- [ S T A R T ] --- %s --- [ S T A R T ] ---\n", label)
	fmt.Fprint(buffer, s)
	fmt.Fprintf(buffer, "--- [   F I N   ] --- %s --- [   F I N   ] ---\n\n", label)

	return buffer.String()
}
