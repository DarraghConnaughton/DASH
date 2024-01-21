package playbackbuffer

import (
	"strconv"
	"sync"
	"time"
)

type PlaybackBuffer struct {
	mu              sync.Mutex
	StartTime       time.Time
	totalTimeLoaded time.Duration
	Buffersize      time.Duration
}

func (pb *PlaybackBuffer) UpdateTotalTimeLoaded(duration string) error {
	seconds, err := strconv.ParseFloat(duration, 64)
	if err != nil {
		return err
	}

	pb.mu.Lock()
	defer pb.mu.Unlock()
	pb.totalTimeLoaded += time.Duration(seconds * float64(time.Second))
	return nil
}

func (pb *PlaybackBuffer) GetTotalTimeLoaded() time.Duration {
	pb.mu.Lock()
	defer pb.mu.Unlock()
	return pb.totalTimeLoaded
}

func New() PlaybackBuffer {
	return PlaybackBuffer{
		Buffersize:      time.Second * 20,
		StartTime:       time.Now(),
		totalTimeLoaded: time.Second * 0,
	}
}
