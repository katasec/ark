package web

import "net/http"

// Middleware is a type alias
type Middleware func(http.HandlerFunc) http.HandlerFunc

// MultipleMiddleware Returns a chain of handlers that were provided as a list in a slice
func MultipleMiddleware(h http.HandlerFunc, m ...Middleware) http.HandlerFunc {

	if len(m) < 1 {
		return h
	}

	chained := h

	// loop in reverse to preserve middleware order
	for i := len(m) - 1; i >= 0; i-- {
		currentHandler := m[i]
		chained = currentHandler(chained)
	}

	return chained

}
