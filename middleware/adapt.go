package main

import (
	"net/http"
)

//Adapter is the type for adapter functions.
//An adapter function accepts an http.Handler
//and returns a new http.Handler that wraps the
//input handler, providing some pre- and/or
//post-processing.
type Adapter func(http.Handler) http.Handler

//TODO: write an Adapt() function that accepts:
// - handler http.Handler the handler to adapt
// - a variadic slice of Adapter functions
//iterate the slice of Adapter functions in
//reverse order, passing the `handler` to
//each, and resetting `handler` to the
//handler returned from the Adapter func
func Adapt(handler http.Handler, adapters ...Adapter) http.Handler {
	for idx := len(adapters) - 1; idx >= 0; idx-- {
		handler = adapters[idx](handler)
	}
	return handler
}
