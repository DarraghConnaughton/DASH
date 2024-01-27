package server

import (
	"dash/pkg/types"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Server struct {
	BindAndServe func(string, http.Handler) error
	errChan      chan error
}

func (ds *Server) LoadRoutes(routes []types.RouteInfo) {
	log.Println("[+] loading routes.")
	router := mux.NewRouter()
	for _, route := range routes {
		localRoute := route
		log.Println(
			fmt.Sprintf("[+] path: %s; description: %s",
				localRoute.Path,
				localRoute.Description))

		router.HandleFunc(localRoute.Path, func(w http.ResponseWriter, r *http.Request) {
			localRoute.HandlerFunc(w, r)
		}).Methods("GET").Name(localRoute.Description)
	}
	http.Handle("/", router)
}

func (ds *Server) Start(port string) {
	log.Println(fmt.Sprintf("[+] starting webserver 127.0.0.1%s", port))
	ds.Server(port, http.ListenAndServe)
}

func (ds *Server) Server(bind string, listenAndServe func(string, http.Handler) error) {
	if err := listenAndServe(bind, nil); err != nil {
		ds.errChan <- err
	}
}
