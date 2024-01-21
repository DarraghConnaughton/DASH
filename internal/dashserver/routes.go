package dashserver

import (
	"dash/pkg/types"
	"github.com/gorilla/mux"
	"net/http"
)

func (ds *DASHServer) loadRoutes() {
	router := mux.NewRouter()
	routes := []types.RouteInfo{
		{
			HandlerFunc: ds.getVideoContents,
			Path:        "/hlsmanifest/{video_uid}/{resolution}/{segment}",
			Description: "get hlsmanifest as specified in the HTTP Request.",
		},
		{
			HandlerFunc: ds.availableRoute,
			Path:        "/available",
			Description: "get list of available videos on webserver",
		},
		{
			HandlerFunc: ds.retrieveManifest,
			Path:        "/manifest/{video_uid}",
			Description: "retrieve manifest associated with video_uid",
		},
	}
	for _, route := range routes {
		localRoute := route

		router.HandleFunc(localRoute.Path, func(w http.ResponseWriter, r *http.Request) {
			localRoute.HandlerFunc(w, r)
		}).Methods("GET").Name(localRoute.Description)
	}
	http.Handle("/", router)
}
