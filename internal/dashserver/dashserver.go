package dashserver

import (
	"dash/pkg/helper"
	"dash/pkg/video"
	"net/http"
	"strings"
)

type DASHServer struct {
	BindAndServe        func(string, http.Handler) error
	errChan             chan error
	fileSource          string
	Videos              []video.Video
	SupportedResolution []string
	VideoIDs            []string
}

func (ds *DASHServer) Start() {
	//initialise go routine which listens out of directory changes
	// TODO
	ds.loadRoutes()
	if err := ds.startServer(":8080", http.ListenAndServe); err != nil {
		ds.errChan <- err
	}
}

func (ds *DASHServer) startServer(bind string, listenAndServe func(string, http.Handler) error) error {
	if err := listenAndServe(bind, nil); err != nil {
		return err
	}
	return nil
}

func (ds *DASHServer) hydrateVideosObject() error {
	var videos []string
	videoDirs, err := helper.ListDirectory(ds.fileSource)
	if err != nil {
		return nil
	}
	for _, videoPath := range videoDirs {
		tmp := strings.Split(videoPath, "/")
		tVideo := video.Video{
			VideoUID: tmp[len(tmp)-1],
		}
		if err := tVideo.LoadManifestsFromFS(ds.SupportedResolution); err != nil {
			return err
		}
		ds.Videos = append(ds.Videos, tVideo)
		videos = append(videos, tmp[len(tmp)-1])
	}
	ds.VideoIDs = videos
	return nil
}

func New(source string, resolutions []string, errChan chan error) DASHServer {
	ds := DASHServer{
		fileSource:          source,
		errChan:             errChan,
		SupportedResolution: resolutions,
	}
	err := ds.hydrateVideosObject()
	if err != nil {
		ds.errChan <- err
	}
	return ds
}
