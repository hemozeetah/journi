// package muxer provides a mux for http handlers.
package muxer

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
	origins   []string
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

	if m.origins != nil {
		handlerFunc = wrapMiddleware([]MidFunc{m.corsHandler}, handlerFunc)
	}

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

// EnableCORS enables CORS preflight requests to work in the middleware. It
// prevents the MethodNotAllowedHandler from being called. This must be enabled
// for the CORS middleware to work.
func (m *Mux) EnableCORS(origins []string) {
	m.origins = origins

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		return nil
	}

	m.HandlerFuncWithoutGlobalMid(http.MethodOptions, "", "/", handler, m.corsHandler)
}

func (m *Mux) corsHandler(handler HandlerFunc) HandlerFunc {
	h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Origin
		// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Origin
		//
		// Limiting the possible Access-Control-Allow-Origin values to a set of
		// allowed origins requires code on the server side to check the value of
		// the Origin request header, compare that to a list of allowed origins, and
		// then if the Origin value is in the list, set the
		// Access-Control-Allow-Origin value to the same value as the Origin.

		reqOrigin := r.Header.Get("Origin")
		for _, origin := range m.origins {
			if origin == "*" || origin == reqOrigin {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				break
			}
		}

		w.Header().Set("Access-Control-Allow-Methods", "POST, PATCH, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		w.Header().Set("Access-Control-Max-Age", "86400")

		return handler(ctx, w, r)
	}

	return h
}

// HandlerFuncWithoutGlobalMid sets a handler function for a given HTTP method and path
// pair to the application server mux. Does not include the application global
// middlewares.
func (m *Mux) HandlerFuncWithoutGlobalMid(method string, group string, path string, handlerFunc HandlerFunc, mws ...MidFunc) {
	handlerFunc = wrapMiddleware(mws, handlerFunc)

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

// FileServer starts a file server based on the specified file system and
// directory inside that file system.
func (m *Mux) FileServer(path string) {
	fileServer := http.StripPrefix("/static", http.FileServer(http.Dir(path)))
	m.mux.Handle("GET /static/", fileServer)
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
