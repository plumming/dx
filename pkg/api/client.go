package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"strings"

	"github.com/plumming/chilly/pkg/util"

	"github.com/jenkins-x/jx-logging/pkg/log"
	"github.com/plumming/chilly/pkg/version"
)

// ClientOption represents an argument to NewClient.
type ClientOption = func(http.RoundTripper) http.RoundTripper

// NewClient initializes a Client.
func NewClient(opts ...ClientOption) *Client {
	tr := http.DefaultTransport
	for _, opt := range opts {
		tr = opt(tr)
	}

	h := &http.Client{Transport: tr}
	client := &Client{http: h}

	return client
}

// AddHeader turns a RoundTripper into one that adds a request header.
func AddHeader(name, value string) ClientOption {
	return func(tr http.RoundTripper) http.RoundTripper {
		return &funcTripper{roundTrip: func(req *http.Request) (*http.Response, error) {
			host := req.Host
			log.Logger().Debugf("sending request to host '%s'", host)
			if name == "Authorization" {
				if host == "api.github.com" {
					log.Logger().Debugf("Adding Authorization Header %s=%s", name, value)
					req.Header.Add(name, value)
				}
			} else {
				log.Logger().Debugf("Adding Header %s=%s", name, value)
				req.Header.Add(name, value)
			}

			return tr.RoundTrip(req)
		}}
	}
}

func enableTracing() ClientOption {
	return func(rt http.RoundTripper) http.RoundTripper {
		return &util.Tracer{RoundTripper: rt}
	}
}

// ReplaceTripper substitutes the underlying RoundTripper with a custom one.
func ReplaceTripper(tr http.RoundTripper) ClientOption {
	return func(http.RoundTripper) http.RoundTripper {
		return tr
	}
}

type funcTripper struct {
	roundTrip func(*http.Request) (*http.Response, error)
}

func (tr funcTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return tr.roundTrip(req)
}

// Client facilitates making HTTP requests to the GitHub API.
type Client struct {
	http *http.Client
}

type graphQLResponse struct {
	Data   interface{}
	Errors []GraphQLError
}

// GraphQLError is a single error returned in a GraphQL response.
type GraphQLError struct {
	Type    string
	Path    []string
	Message string
}

// GraphQLErrorResponse contains errors returned in a GraphQL response.
type GraphQLErrorResponse struct {
	Errors []GraphQLError
}

func (gr GraphQLErrorResponse) Error() string {
	errorMessages := make([]string, 0, len(gr.Errors))
	for _, e := range gr.Errors {
		errorMessages = append(errorMessages, e.Message)
	}
	return fmt.Sprintf("graphql error: '%s'", strings.Join(errorMessages, ", "))
}

// GraphQL performs a GraphQL request and parses the response.
func (c Client) GraphQL(query string, variables map[string]interface{}, data interface{}) error {
	url := "https://api.github.com/graphql"
	reqBody, err := json.Marshal(map[string]interface{}{"query": query, "variables": variables})
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/vnd.github.antiope-preview+json")

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return handleResponse(resp, data)
}

func handleResponse(resp *http.Response, data interface{}) error {
	success := resp.StatusCode >= 200 && resp.StatusCode < 300

	if !success {
		return handleHTTPError(resp)
	}

	body, err := ioutil.ReadAll(io.LimitReader(resp.Body, 10000000))
	if err != nil {
		return err
	}

	log.Logger().Debugf("got response %s", body)
	gr := &graphQLResponse{Data: data}
	err = json.Unmarshal(body, &gr)
	if err != nil {
		return err
	}

	if len(gr.Errors) > 0 {
		return &GraphQLErrorResponse{Errors: gr.Errors}
	}
	return nil
}

// REST performs a REST request and parses the response.
func (c Client) REST(method string, p string, body io.Reader, data interface{}) error {
	url := "https://api.github.com/" + p
	log.Logger().Debugf("checking url: %s", url)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	success := resp.StatusCode >= 200 && resp.StatusCode < 300
	if !success {
		log.Logger().Debugf("failed with resp code %d", resp.StatusCode)
		return handleHTTPError(resp)
	}

	if resp.StatusCode == http.StatusNoContent {
		return nil
	}

	b, err := ioutil.ReadAll(io.LimitReader(resp.Body, 10000000))
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, &data)
	if err != nil {
		return err
	}

	return nil
}

// Download performs a REST request and parses the response.
func (c Client) Download(url string) (*http.Response, error) {
	log.Logger().Debugf("getting url: %s", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/octet-stream")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	//defer resp.Body.Close()

	success := resp.StatusCode >= 200 && resp.StatusCode < 300
	if !success {
		log.Logger().Debugf("failed with resp code %d", resp.StatusCode)
		return nil, handleHTTPError(resp)
	}

	if resp.StatusCode == http.StatusNoContent {
		return nil, errors.New("no content")
	}

	return resp, nil
}

func (c Client) Do(req *http.Request) (*http.Response, error) {
	return c.http.Do(req)
}

func handleHTTPError(resp *http.Response) error {
	var message string
	var parsedBody struct {
		Message string `json:"message"`
	}
	body, err := ioutil.ReadAll(io.LimitReader(resp.Body, 10000000))
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &parsedBody)
	if err != nil {
		message = string(body)
	} else {
		message = parsedBody.Message
	}

	return fmt.Errorf("http error, '%s' failed (%d): '%s'", resp.Request.URL, resp.StatusCode, message)
}

func BasicClient() (*Client, error) {
	var opts []ClientOption

	// for testing purposes one can enable tracing of GH REST API calls.
	traceGithubAPI := os.Getenv("TRACE_GITHUB_API")

	if traceGithubAPI == "1" || traceGithubAPI == "on" || traceGithubAPI == "true" {
		opts = append(opts, enableTracing())
	}

	opts = append(opts, AddHeader("User-Agent", fmt.Sprintf("Chilly CLI %s", version.Version)))

	if c, err := ParseDefaultConfig(); err == nil {
		if token := c.GetToken(defaultHostname); token != "" {
			log.Logger().Debugf("Using Auth %s", token)
			opts = append(opts, AddHeader("Authorization", fmt.Sprintf("token %s", token)))
		}
	}
	return NewClient(opts...), nil
}
