package olcode

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Service is the HTTP service
type Service struct {
	server *http.Server
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

// NewService cretes a service
func NewService(port, staticPath string) *Service {
	h := newHandler(staticPath)

	r := mux.NewRouter()
	r.HandleFunc("/api/login", h.login)
	r.HandleFunc("/api/create_room", h.create)
	r.HandleFunc("/api/ws/attend", h.attend)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(staticPath))))
	r.HandleFunc("/", h.serverHome)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
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
