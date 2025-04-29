package webserver

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

type WebServer struct {
	Router                chi.Router
	Handlers              map[string]http.HandlerFunc
	CustomizedMiddlewares []func(http.Handler) http.Handler
	WebServerPort         string
}

func NewWebServer(serverPort string) *WebServer {
	return &WebServer{
		Router:                chi.NewRouter(),
		Handlers:              make(map[string]http.HandlerFunc),
		CustomizedMiddlewares: []func(http.Handler) http.Handler{},
		WebServerPort:         serverPort,
	}
}

func (s *WebServer) AddHandler(path string, handler http.HandlerFunc) {
	s.Handlers[path] = handler
}

func (s *WebServer) AddCustomizedMiddleware(middleware func(http.Handler) http.Handler) {
	s.CustomizedMiddlewares = append(s.CustomizedMiddlewares, middleware)
}

// loop through the handlers and add them to the router
// register middeleware logger
// loop through the customized middlewares and add them to the router
// start the server
func (s *WebServer) Start() {
	s.Router.Use(middleware.Logger)

	// loop to add middlewares
	for _, customizedMiddleware := range s.CustomizedMiddlewares {
		s.Router.Use(customizedMiddleware)
	}

	for path, handler := range s.Handlers {
		s.Router.Handle(path, handler)
	}
	http.ListenAndServe(s.WebServerPort, s.Router)
}
