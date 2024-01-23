package server

import (
	"dash/pkg/types"
	"dash/pkg/video"
	"github.com/gorilla/mux"
	"net/http"
)

type Server struct {
	BindAndServe        func(string, http.Handler) error
	errChan             chan error
	fileSource          string
	Videos              []video.Video
	SupportedResolution []string
	VideoIDs            []string
}

func (ds *Server) LoadRoutes(routes []types.RouteInfo) {
	router := mux.NewRouter()
	for _, route := range routes {
		localRoute := route

		router.HandleFunc(localRoute.Path, func(w http.ResponseWriter, r *http.Request) {
			localRoute.HandlerFunc(w, r)
		}).Methods("GET").Name(localRoute.Description)
	}
	http.Handle("/", router)
}

func (ds *Server) Start(port string) {
	//initialise go routine which listens out of docker changes
	// TODO
	if err := ds.Server(port, http.ListenAndServe); err != nil {
		ds.errChan <- err
	}
}

func (ds *Server) Server(bind string, listenAndServe func(string, http.Handler) error) error {
	if err := listenAndServe(bind, nil); err != nil {
		return err
	}
	return nil
}
