package server

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/Igoraamc/knot-cloud-storage/pkg/controllers"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

type Server struct {
	port           string
	dataController *controllers.DataController
}

func NewServer(port string, dataController *controllers.DataController) Server {
	return Server{port, dataController}
}

func (s *Server) Start() {
	routers := s.createRouters()
	fmt.Println("Server running in port:", s.port)
	log.Fatal(http.ListenAndServe(s.port, routers))
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
	n := negroni.Classic()

	n.Use(s.checkDeviceIdPermission())

	r := mux.NewRouter().StrictSlash(true)

	n.UseHandler(r)

	r.HandleFunc("/data/{deviceId}", s.dataController.GetAll).Methods("GET")
	r.HandleFunc("/data/{deviceId}/sensor/{id}", s.dataController.GetByID).Methods("GET")
	r.HandleFunc("/data", s.dataController.Create).Methods("POST")

	return n
}
