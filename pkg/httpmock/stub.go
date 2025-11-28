package httpmock

import (
	"bytes"
	"io"
	"net/http"
)

type Matcher func(req *http.Request) bool
type Responder func(req *http.Request) (*http.Response, error)

type Stub struct {
	matched   bool
	Matcher   Matcher
	Responder Responder
}

func MatchAny(*http.Request) bool {
	return true
}

func StringResponse(body string) Responder {
	return func(*http.Request) (*http.Response, error) {
		return httpResponse(200, bytes.NewBufferString(body)), nil
	}
}

func httpResponse(status int, body io.Reader) *http.Response {
	return &http.Response{
		StatusCode: status,
		Body:       io.NopCloser(body),
	}
}
