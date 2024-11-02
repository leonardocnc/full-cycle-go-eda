package webserver

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

type WebServer struct {
	Router        chi.Router
	Handlers      map[string]http.HandlerFunc
	GetHandlers   map[string]http.HandlerFunc
	WebServerPort string
}

func NewWebServer(webServerPort string) *WebServer {
	return &WebServer{
		Router:        chi.NewRouter(),
		Handlers:      make(map[string]http.HandlerFunc),
		GetHandlers:   make(map[string]http.HandlerFunc),
		WebServerPort: webServerPort,
	}
}

func (s *WebServer) AddHandler(path string, handler http.HandlerFunc) {
	s.Handlers[path] = handler
}

func (s *WebServer) AddGetHandler(path string, handler http.HandlerFunc) error {
	if _, exists := s.GetHandlers[path]; exists {
		return fmt.Errorf("GET handler already exists for path: %s", path)
	}
	s.GetHandlers[path] = handler
	return nil
}

func (s *WebServer) Start() error {
	s.Router.Use(middleware.Logger)
	for path, handler := range s.Handlers {
		s.Router.Post(path, handler)
	}

	for path, handler := range s.GetHandlers {
		fmt.Printf("Registering GET handler for path: %s\n", path)
		s.Router.Get(path, handler)
	}

	if err := http.ListenAndServe(s.WebServerPort, s.Router); err != nil {
		return fmt.Errorf("failed to start web server: %w", err)
	}
	return nil
}
