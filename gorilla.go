package endpoint

import (
	"fmt"
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

// Gorilla-specific implementations of the RouteParams and RequestContext
// interfaces.

type GorillaRouteParams map[string]string

func (grp GorillaRouteParams) Get(key string) (string, bool) {
	value, ok := grp[key]
	return value, ok
}

func (grp GorillaRouteParams) Require(key string) string {
	value, ok := grp.Get(key)
	if !ok {
		panic(fmt.Errorf("Key missing from GorillaRouteParams: %s", key))
	}
	return value
}

type GorillaContext struct {
	R *http.Request
}

func (gc *GorillaContext) Get(key interface{}) (interface{}, bool) {
	return context.GetOk(gc.R, key)
}
func (gc *GorillaContext) Require(key interface{}) interface{} {
	result, ok := gc.Get(key)
	if !ok {
		panic(fmt.Errorf("Key missing from GorillaContext: %#v", key))
	}
	return result
}
func (gc *GorillaContext) Set(key interface{}, value interface{}) {
	context.Set(gc.R, key, value)
}

func GorillaInitialize(input *Input) {
	input.RouteParams = GorillaRouteParams(mux.Vars(input.R))
	input.Context = &GorillaContext{R: input.R}
}
