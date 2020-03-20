package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/CESARBR/knot-cloud-storage/pkg/controllers"
	"github.com/CESARBR/knot-cloud-storage/pkg/logging"
	"github.com/urfave/negroni"

	"github.com/gorilla/mux"
)

// Health represents the service's health status
type Health struct {
	Status string `json:"status"`
}

// Server represents the HTTP server
type Server struct {
	port   int
	logger logging.Logger
}

// NewServer creates a new server instance
func NewServer(port int, logger logging.Logger) Server {
	return Server{port: port, logger: logger}
}

// Start starts the http server
func (s *Server) Start() {
	routers := s.createRouters()
	s.logger.Infof("Listening on %d", s.port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", s.port), s.logRequest(routers))
	if err != nil {
		s.logger.Error(err)
	}
}

func (s *Server) checkDeviceIdPermission() negroni.Handler {
	return negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

		url := strings.Split(r.URL.String(), "?")

		params := strings.Split(url[0], "/")

		if r.Method == "POST" || params[2] == "teste-12345" {
			next(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
		}
	})
}

func (s *Server) createRouters() *negroni.Negroni {
	dataController := controllers.NewDataController()
	n := negroni.Classic()

	n.Use(s.checkDeviceIdPermission())

	r := mux.NewRouter().StrictSlash(true)

	n.UseHandler(r)

	r.HandleFunc("/data/{deviceId}", dataController.GetAll).Methods("GET")
	r.HandleFunc("/data/{deviceId}/sensor/{id}", dataController.GetByID).Methods("GET")
	r.HandleFunc("/data", dataController.Create).Methods("POST")

	return n
}

func (s *Server) logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.logger.Infof("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func (s *Server) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response, _ := json.Marshal(&Health{Status: "online"})
	_, err := w.Write(response)
	if err != nil {
		s.logger.Errorf("Error sending response, %s\n", err)
	}
}
