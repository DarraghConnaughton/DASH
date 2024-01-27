package types

import (
	"net/http"
)

type RouteInfo struct {
	HandlerFunc func(http.ResponseWriter, *http.Request)
	Path        string
	Description string
}

type AvailableVideos struct {
	Available []string `json:"available"`
}

type Segment struct {
	Location string
	Duration string
}

type HLSManifest struct {
	Resolution       string
	Version          int
	TargetDuration   int
	MediaSequence    int
	SegmentLocations map[int]Segment
}

type NoiseData struct {
	FileContents []byte `json:"contents"`
	Noise        []byte `json:"noise"`
}

type NetworkTraceData struct {
	Timestamp  string
	Bytes      string
	Sequence   string
	Resolution string
}

type RPCHeartBeat struct {
	UID                string
	CPU                int
	NumberOfGoroutines int
	AllocatedMemory    uint64
	TotalAlloc         uint64
	SysMem             uint64
}
