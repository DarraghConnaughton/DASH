package dashserver

import (
	"dash/pkg/helper"
	"dash/pkg/types"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func (ds *DASHServer) getVideoContents(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	for _, v := range ds.Videos {
		if v.VideoUID == vars["video_uid"] {
			w.Header().Set("Content-Type", "hlsmanifest/mp4")
			fp := fmt.Sprintf("%s/%s/%s/%s%s.ts",
				ds.fileSource,
				vars["video_uid"],
				vars["resolution"],
				vars["resolution"],
				vars["segment"])
			if helper.Exists(fp) {
				http.ServeFile(w, r, fp)
				return
			}
		}
	}
	http.Error(w, fmt.Sprintf("%s not found", vars["video_uid"]), http.StatusNotFound)
}

func (ds *DASHServer) availableRoute(w http.ResponseWriter, _ *http.Request) {
	jsonResponse, err := json.Marshal(types.AvailableVideos{
		Available: ds.VideoIDs,
	})
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(jsonResponse); err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusInternalServerError)
	}
}

func (ds *DASHServer) retrieveManifest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	for _, v := range ds.Videos {
		if v.VideoUID == vars["video_uid"] {
			jsonResponse, err := json.Marshal(v.EncodedRepresentations)
			if err != nil {
				http.Error(w, "", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write(jsonResponse); err != nil {
				http.Error(w, "", http.StatusInternalServerError)
			}
			return
		}
	}
	http.Error(w, "", http.StatusNotFound)
	return
}
