package fileserver

import (
	"dash/cmd/helper"
	"dash/cmd/types"
	"dash/cmd/video"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strings"
)

type VideoServer struct {
	BindAndServe        func(string, http.Handler) error
	errChan             chan error
	fileSource          string
	routeInfo           [][]types.RouteInfo
	Videos              []video.Video
	SupportedResolution []string
	VideoIDs            []string
}

func (fs *VideoServer) Start() {

	err := fs.hydrateVideosObject()
	if err != nil {
		log.Println(err)
	}

	var videos []string
	for _, v := range fs.Videos {
		log.Println("WE323423")
		log.Println(v.VideoUID)
		videos = append(videos, v.VideoUID)
	}
	fs.VideoIDs = videos

	//initialise go routine which listens out of directory changes
	// TODO
	fs.loadRoutes()
	if err := fs.startServer(":8080", http.ListenAndServe); err != nil {
		fs.errChan <- err
	}
}

//func (fs *VideoServer) retrieveMP4Files() ([]types.Video, error) {
//	var files []types.Video
//	dir, err := os.Open(fs.fileSource)
//	if err != nil {
//		return files, err
//	}
//
//	defer func() {
//		if err := dir.Close(); err != nil {
//			log.Println(fmt.Sprintf(
//				"[-]warning: error encountered when attempting to close file: %s", err.Error()))
//		}
//	}()
//
//	fileInfos, err := dir.Readdir(-1)
//	if err != nil {
//		return files, err
//	}
//	for _, fileInfo := range fileInfos {
//		if strings.Contains(fileInfo.Name(), "mp4") {
//			files = append(files, fileInfo.Name())
//		}
//	}
//	return files, nil
//}

// NEEDS TO BE EXPANDED; WE SHOULD TAKE AN OFFSET AS WELL, NOT SOLELY THE VIDEO_UID
// NEEDS TO BE EXPANDED; WE SHOULD TAKE AN OFFSET AS WELL, NOT SOLELY THE VIDEO_UID
// NEEDS TO BE EXPANDED; WE SHOULD TAKE AN OFFSET AS WELL, NOT SOLELY THE VIDEO_UID
// NEEDS TO BE EXPANDED; WE SHOULD TAKE AN OFFSET AS WELL, NOT SOLELY THE VIDEO_UID
func (fs *VideoServer) getVideoContents(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	for _, v := range fs.Videos {
		fmt.Println("**********")
		fmt.Println(v.VideoUID)
		fmt.Println(vars["video_uid"])

		if v.VideoUID == vars["video_uid"] {
			w.Header().Set("Content-Type", "hlsmanifest/mp4")
			fmt.Println(fs.fileSource)
			log.Println(fmt.Sprintf("%s/%s/%sp/%sp_00%s.ts",
				fs.fileSource,
				vars["video_uid"],
				vars["resolution"],
				vars["resolution"],
				vars["segment"]))

			fmt.Println("fs.fileSourcefs.fileSourcefs.fileSource")

			//Likely bug here, what happens when segment number rolls over to
			//-best to use a different naming convention if possible.
			http.ServeFile(w, r,
				fmt.Sprintf("%s/%s/%s/%s_00%s.ts",
					fs.fileSource,
					vars["video_uid"],
					vars["resolution"],
					vars["resolution"],
					vars["segment"]))
			return
		}
	}
	http.Error(w, fmt.Sprintf("%s not found", vars["video_uid"]), http.StatusNotFound)
}

//func contains(files []string, file string) bool {
//	for _, f := range files {
//		if f == file {
//			return true
//		}
//	}
//	return false
//}

func (fs *VideoServer) loadRoutes() {
	router := mux.NewRouter()
	routes := []types.RouteInfo{
		{
			HandlerFunc: fs.getVideoContents,
			Path:        "/hlsmanifest/{video_uid}/{resolution}/{segment}",
			Description: "get hlsmanifest as specified in the HTTP Request.",
		},
		{
			HandlerFunc: fs.availableRoute,
			Path:        "/available",
			Description: "get list of available videos on webserver",
		},
		{
			HandlerFunc: fs.retrieveManifest,
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

func (fs *VideoServer) availableRoute(w http.ResponseWriter, r *http.Request) {
	jsonResponse, err := json.Marshal(types.AvailableVideos{
		Available: fs.VideoIDs,
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

func (fs *VideoServer) startServer(bind string, listenAndServe func(string, http.Handler) error) error {
	log.Println("BINDING SERVER!! ")
	log.Println(bind)
	log.Println("BINDING SERVER!! ")
	if err := listenAndServe(bind, nil); err != nil {
		log.Println("Error encountered during binding....")
		log.Println(err)
		return err
	}
	return nil
}

func (fs *VideoServer) retrieveManifest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	fmt.Println(vars)
	fmt.Println("****88*****")
	for _, v := range fs.Videos {
		if v.VideoUID == vars["video_uid"] {
			fmt.Println(v.EncodedRepresentations)
			jsonResponse, err := json.Marshal(v.EncodedRepresentations)
			fmt.Println(jsonResponse)
			if err != nil {
				http.Error(w, "", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write(jsonResponse); err != nil {
				log.Println(err)
				http.Error(w, "", http.StatusInternalServerError)
			}
			return
		}
	}

	http.Error(w, "", http.StatusNotFound)
	return
}

func (fs *VideoServer) hydrateVideosObject() error {
	videoDirs, err := helper.ListDirectory(fs.fileSource)
	if err != nil {
		return nil
	}
	fmt.Println("234234324324234")
	fmt.Println(videoDirs)
	fmt.Println("234234324324234")
	for _, videoPath := range videoDirs {
		tmp := strings.Split(videoPath, "/")
		log.Println("tmp ")
		log.Println(tmp[len(tmp)-1])
		log.Println("tmp ")
		tVideo := video.Video{
			VideoUID: tmp[len(tmp)-1],
		}
		if err := tVideo.LoadManifests(fs.SupportedResolution); err != nil {
			return err
		}
		fs.Videos = append(fs.Videos, tVideo)
	}
	return nil
}

func New(source string, resolutions []string, errChan chan error) VideoServer {
	return VideoServer{
		fileSource:          source,
		errChan:             errChan,
		SupportedResolution: resolutions,
	}
}
