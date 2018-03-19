package olcode

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// Service is the HTTP service
type Service struct {
	server *http.Server
}

// NewService cretes a service
func NewService(port int16) *Service {
	h := newHandler()

	r := mux.NewRouter()
	r.HandleFunc("/login", h.login)
	r.HandleFunc("/create", h.create)
	r.HandleFunc("/attend", h.attend)
	r.HandleFunc("/leave", h.leave)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: r,
	}

	return &Service{
		server: server,
	}
}

// Serve starts the service
func (s *Service) Serve() error {
	return s.server.ListenAndServe()
}
