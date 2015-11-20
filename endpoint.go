package endpoint

import (
	"net/http"
)

type RouteParams interface {
	Get(key string) (string, bool)
	Require(key string) string
}

type RequestContext interface {
	Get(key interface{}) (interface{}, bool)
	Require(key interface{}) interface{}
	Set(key interface{}, value interface{})
}

type Input struct {
	// The response.
	W http.ResponseWriter
	R *http.Request
	// Holds URL-routing parameters, like "id" in the route "/clips/:id'.
	RouteParams RouteParams
	// Key/value interface for passing data between middleware and the endpoint.
	Context RequestContext
}

type Output struct {
	StatusCode int
	Headers    map[string]string
	Data       interface{}
}

type Endpoint func(*Input) *Output

type Initializer func(input *Input)
type Finalizer func(output *Output)
type Byter func(output *Output) []byte
type Control struct {
	Initialize Initializer
	Finalize   Finalizer
	Bytes      Byter
}

func NewControl(i Initializer, f Finalizer, b Byter) *Control {
	return &Control{i, f, b}
}

func NullInitialize(input *Input) {}
func NullFinalize(output *Output) {}
func NullBytes(output *Output) []byte {
	return []byte{}
}

func (control Control) Handler(endpoint Endpoint) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		input := &Input{W: w, R: r}
		control.Initialize(input) // set up
		output := endpoint(input)
		if output == nil {
			return
		}
		control.Finalize(output) // set headers, etc.
		data := control.Bytes(output)
		// The Header map must be updated before the call to WriteHeader, see
		// http://golang.org/pkg/net/http/#ResponseWriter. This cannot be checked
		// for in tests due to the way `httptest.ResponseRecorder` mocks the
		// `WriteHeader` call.
		if output.Headers != nil {
			for headerKey, headerValue := range output.Headers {
				w.Header().Set(headerKey, headerValue)
			}
		}
		w.WriteHeader(output.StatusCode)
		w.Write(data)
	})
}
