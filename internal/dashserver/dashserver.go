package dashserver

import (
	"dash/pkg/helper"
	"dash/pkg/server"
	"dash/pkg/types"
	"dash/pkg/video"
	"strings"
)

type DASHServer struct {
	server.Server
	errChan             chan error
	fileSource          string
	Videos              []video.Video
	SupportedResolution []string
	VideoIDs            []string
}

func (ds *DASHServer) getRouteInfo() []types.RouteInfo {
	return []types.RouteInfo{
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

func (ds *DASHServer) Start(port string) {
	ds.Server.LoadRoutes(ds.getRouteInfo())
	go ds.Server.Start(port)
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
