package types

import "net/http"

type RouteInfo struct {
	HandlerFunc func(http.ResponseWriter, *http.Request)
	Path        string
	Description string
}

type AvailableVideos struct {
	Available []string `json:"available"`
}

//type Video struct {
//	VideoUID               string
//	EncodedRepresentations []*HLSManifest
//}

type HLSManifest struct {
	Resolution       string
	Version          int
	TargetDuration   int
	MediaSequence    int
	SegmentLocations map[int]string
}
