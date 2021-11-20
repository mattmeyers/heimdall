package http

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mattmeyers/heimdall/logger"
)

// Registerer represents a type that can register its routes. Generally, this is a controller
// that registers all of the routes for its domain.
type Registerer interface {
	Register(*httprouter.Router)
}

// Middleware is a function that wraps a handler with additional behavior.
type Middleware func(http.Handler) http.Handler

// Chain wraps a handler with multiple middleware. The provided middleware are iterated through
// in reverse order. This means that the first provided middleware will be called first when a
// request is handled. For example, given the following function call
//
//		Chain(baseHandler, loggingMiddleware, authMiddleware)
//
// a request would go through the loggingMiddleware, then the authMiddleware, then be handled
// by the base handler.
func Chain(h http.Handler, ms ...Middleware) http.Handler {
	for i := len(ms) - 1; i >= 0; i-- {
		h = ms[i](h)
	}

	return h
}

// Server is a wrapper around net/http.Server that allows for route registration. The server
// provides the standard loggingMiddleware which will be called before all requests.
type Server struct {
	s      *http.Server
	router *httprouter.Router
	logger logger.Logger

	loggingMiddleware Middleware
}

// NewServer constructs a new server object with all of the default values.
func NewServer(addr string, logger logger.Logger) (*Server, error) {
	return &Server{
		s:                 &http.Server{Addr: addr},
		router:            httprouter.New(),
		logger:            logger,
		loggingMiddleware: newLoggingMiddleware(logger),
	}, nil
}

// RegisterRoutes registers all of the provided routes. If a route is registered multiple times,
// this function will panic.
func (s *Server) RegisterRoutes(rs ...Registerer) {
	for _, r := range rs {
		r.Register(s.router)
	}
}

// ListenAndServe begins listening for requests, and handles them using the registered handlers.
// This method will always return a non-nil error.
func (s *Server) ListenAndServe() error {
	s.logger.Info("Server starting on %s...", s.s.Addr)

	s.s.Handler = Chain(s.router, s.loggingMiddleware)
	return s.s.ListenAndServe()
}
