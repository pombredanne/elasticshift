package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Methods
const (
	GET  = "GET"
	POST = "POST"
)

var (
	pathParamIndicator = ":"
	urlSeparator       = "/"

	errPathParamConflift          = errors.New("Path parameter placeholder and actual vaues doesn't match")
	errMethodUnknown              = errors.New("Unknown request METHOD")
	errCannotSetBodyOnPostRequest = errors.New("Body can only be set for GET request type")
)

// RequestMaker ...
type RequestMaker struct {
	method      string
	url         string
	headers     map[string]string
	pathParams  []string
	queryParams map[string]string
	body        interface{}
	response    interface{}
}

// NewGetRequestMaker ..
// Create new request maker of GET request
func NewGetRequestMaker(url string) *RequestMaker {

	return &RequestMaker{
		method:      GET,
		url:         url,
		headers:     make(map[string]string),
		queryParams: make(map[string]string),
	}
}

// NewPostRequestMaker ..
// Create new request maker of GET request
func NewPostRequestMaker(url string) *RequestMaker {

	return &RequestMaker{
		method:      POST,
		url:         url,
		headers:     make(map[string]string),
		queryParams: make(map[string]string),
	}
}

// PathParams ..
// Inject params to the url path described with :
// Ex: http://elasticshift.com/api/users/:name
func (r *RequestMaker) PathParams(params ...string) *RequestMaker {
	r.pathParams = params
	return r
}

// QueryParam ..
// Set a query paramter to a request
func (r *RequestMaker) QueryParam(key, value string) *RequestMaker {
	r.queryParams[key] = value
	return r
}

// Header ..
// Set a header value to a request
func (r *RequestMaker) Header(key, value string) *RequestMaker {
	r.headers[key] = value
	return r
}

// Body ..
// Set the request struct which will be converted to json during request
func (r *RequestMaker) Body(request interface{}) *RequestMaker {
	r.body = request
	return r
}

// Scan ..
// Maps the response to response struct
func (r *RequestMaker) Scan(response interface{}) *RequestMaker {
	r.response = response
	return r
}

// Dispatch ..
// This is where actuall request made to destination
func (r *RequestMaker) Dispatch() error {

	// Set the path params
	splits := strings.Split(r.url, urlSeparator)
	var idx int
	for i, s := range splits {

		if strings.HasPrefix(s, pathParamIndicator) {
			splits[i] = r.pathParams[idx]
			idx++
		}
	}

	// Verify all the path params are set
	if idx != len(r.pathParams) {
		return errPathParamConflift
	}

	// sets the final url after injecting path params
	r.url = strings.Join(splits, urlSeparator)

	if r.method == "" {
		return errMethodUnknown
	}

	// create a request
	req, err := http.NewRequest(r.method, r.url, nil)
	if err != nil {
		return err
	}

	// Sets the header
	if len(r.headers) > 0 {

		for k, v := range r.headers {
			req.Header.Add(k, v)
		}
	}

	// Set the query params
	if len(r.queryParams) > 0 {

		q := req.URL.Query()
		for k, v := range r.queryParams {
			q.Add(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}

	// Set the body
	if r.body != nil {

		if r.method != GET {
			return errCannotSetBodyOnPostRequest
		}

		bits, err := json.Marshal(r.body)
		if err != nil {
			return err
		}
		_, err = req.Body.Read(bits)
		if err != nil {
			return err
		}
	}

	fmt.Println("Making request to = ", req.URL.String())

	// dispatch the request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	// Scans the response
	if err != nil {
		if res != nil {
			res.Body.Close()
		}
		return err
	}
	defer res.Body.Close()

	// read the response body
	bits, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	fmt.Println("Response = ", string(bits[:]))
	// decode to response type
	err = json.NewDecoder(bytes.NewReader(bits)).Decode(r.response)
	if err != nil {
		return err
	}
	return nil
}
