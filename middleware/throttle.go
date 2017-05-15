package main

import (
	"net/http"
	"time"

	cache "github.com/patrickmn/go-cache"
)

/*
TODO: Similar to the LogRequests middleware function, define a
ThrottleRequests middleware function here that accepts two parameters:
- maxRequests int = max number of requests a client can make during the duration
- duration time.Duration = a duration during which a client can make up to maxRequests

Like LogRequests, this ThrottleRequests function should return an Adapter function.
The Adapeter function accepts an http.Handler function and returns an http.Handler
function. The returned handler should check how many requests the client has made
already, and if the client has already exceeded maxRequests, respond with an
http.StatusTooManyRequests. If not, call the original handler.

Since we don't have an authentication context here, use the request struct's
RemoteAddr field to identify the client. This should contain the IP address
of the client who made the request.

To track how many requests each client has made so far, use the go-cache
package to create a new in-memory TTL cache. https://github.com/patrickmn/go-cache
This cache is safe for concurrent access, so you can share it amongst
multiple requests (which are processed concurrently in Go)

Or if you're feeling adventurous, spin up a redis server using Docker,
connect to it in your main() function, and pass a pointer to the redis client
as a third parameter to your ThrottleRequests function.
*/

func throttleRequests(maxRequests int, duration time.Duration) Adapter {
	c := cache.New(duration, time.Second)
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ipAddr := r.RemoteAddr
			value, expTime, found := c.GetWithExpiration(ipAddr)
			var numReqs int
			//value, found := c.Get(ipAddr)
			if !found {
				numReqs = 0
			} else {
				numReqs = value.(int)
			}
			numReqs++

			if numReqs == 1 {
				c.Set(ipAddr, numReqs, duration)
				handler.ServeHTTP(w, r)
			} else if numReqs <= maxRequests {
				c.Set(ipAddr, numReqs, time.Until(expTime))
				handler.ServeHTTP(w, r)
			} else {
				http.Error(w, "Too many", http.StatusTooManyRequests)
			}

		})
	}
}
