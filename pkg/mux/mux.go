// package mux provides a mux for http handlers.
package mux

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hemozeetah/journi/pkg/logger"
)

// HandlerFunc represents a function that handles a http request.
type HandlerFunc func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

// MidFunc is a handler function designed to run code before and/or after
// another Handler.
type MidFunc func(handler HandlerFunc) HandlerFunc

// Mux is what configure each of our http handlers.
type Mux struct {
	log       *logger.Logger
	mux       *http.ServeMux
	globalMWs []MidFunc
}

// New creates a Mux to handle a set of routes.
// The global middlewares' handlers will be executed by requests in the order
// they are provided.
func New(log *logger.Logger, globalMWs ...MidFunc) *Mux {
	return &Mux{
		mux:       http.NewServeMux(),
		log:       log,
		globalMWs: globalMWs,
	}
}

// ServeHTTP implements the http.Handler interface.
func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.mux.ServeHTTP(w, r)
}

// HandlerFunc sets a handler function for a given HTTP method and path.
// The middlewares' handlers will be executed by requests in the order they
// are provided after the global middlewares.
func (m *Mux) HandlerFunc(method string, group string, path string, handlerFunc HandlerFunc, mws ...MidFunc) {
	handlerFunc = wrapMiddleware(mws, handlerFunc)
	handlerFunc = wrapMiddleware(m.globalMWs, handlerFunc)

	h := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		if err := handlerFunc(ctx, w, r); err != nil {
			m.log.Info(ctx).
				Attr("ERROR", err).
				Msg("handler")
			return
		}
	}

	finalPath := path
	if group != "" {
		finalPath = "/" + group + path
	}
	finalPath = fmt.Sprintf("%s %s", method, finalPath)

	m.mux.HandleFunc(finalPath, h)
}

func wrapMiddleware(mw []MidFunc, handler HandlerFunc) HandlerFunc {
	for i := len(mw) - 1; i >= 0; i-- {
		mwFunc := mw[i]
		if mwFunc != nil {
			handler = mwFunc(handler)
		}
	}

	return handler
}
