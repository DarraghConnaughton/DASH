package dashclient

import (
	"bytes"
	"dash/pkg/https"
	"dash/pkg/playbackbuffer"
	"dash/pkg/types"
	"dash/pkg/video"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"
)

type DASHClient struct {
	ByteChan             chan []byte
	ErrorChan            chan error
	HTTPS                https.HTTPS
	PlaybackBuffer       playbackbuffer.PlaybackBuffer
	SupportedResolutions []string
	Video                video.Video
	PreviousRTT          time.Duration
}

func (c *DASHClient) Watch(videouid string) {

	err := c.retrieveVideoManifest(videouid)
	if err != nil {
		c.ErrorChan <- err
	}
	// Preload the buffer.
	startsegment := c.initialSequentialBufferHydration()

	// Gather the remainder of the DASH segments in accordance with playback buffer state.
	go c.gatherDASHSegments(startsegment)
	go c.serveVideoViaTCP()
}

func (c *DASHClient) establishListener() net.Listener {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		c.ErrorChan <- err
	}
	return listener
}

func (c *DASHClient) launchVideoWithFFPlay(port string) *os.Process {
	process, err := triggerFFPlayChildProcess(port)
	if err != nil {
		c.ErrorChan <- err
	}
	return process
}

func (c *DASHClient) serveVideoViaTCP() {

	// Set up listener
	listener := c.establishListener()
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			c.ErrorChan <- err
		}
	}(listener)

	// Launch video which reads from previous established TCP port
	port := strings.Split(listener.Addr().String(), "[::]")[1]
	process := c.launchVideoWithFFPlay(port)
	defer func(process *os.Process) {
		err := process.Kill()
		if err != nil {
			c.ErrorChan <- err
		}
	}(process)

	// Accept connection and proceed to handler function.
	conn := c.acceptConnections(listener)
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			c.ErrorChan <- err
		}
	}(conn)

	c.handleConnection(conn)

}

func (c *DASHClient) determineQuality() interface{} {
	quality := "1080p"
	durationNanos := int64(c.PreviousRTT)
	switch {
	case durationNanos >= int64(1*time.Second) && durationNanos < int64(2*time.Second):
		quality = "720p"
	case durationNanos >= int64(2*time.Second) && durationNanos < int64(3*time.Second):
		quality = "480p"
	case durationNanos >= int64(3*time.Second) && durationNanos < int64(4*time.Second):
		quality = "380p"
	}
	return quality
}

func (c *DASHClient) loadDashSegment(segmentindex int) {
	// Convert time.Duration to int64 representing nanoseconds
	quality := c.determineQuality()
	resp, err := c.HTTPS.GenericMethod(
		fmt.Sprintf("http://127.0.0.1:8080/hlsmanifest/%s/%s/%d", c.Video.VideoUID, quality, segmentindex))
	if err != nil {
		c.ErrorChan <- err
	}
	c.ByteChan <- resp.Bytes
	c.PlaybackBuffer.UpdateTotalTimeLoaded(c.Video.EncodedRepresentations[0].SegmentLocations[segmentindex].Duration)
}

func (c *DASHClient) gatherDASHSegments(startsegment int) {

	defer close(c.ByteChan)
	for i := startsegment; i < len(c.Video.EncodedRepresentations[0].SegmentLocations); i++ {
		for c.PlaybackBuffer.GetTotalTimeLoaded()-time.Since(c.PlaybackBuffer.StartTime) > c.PlaybackBuffer.Buffersize {
			time.Sleep(2)
		}

		log.Println(fmt.Sprintf(
			"buffer hydration in process. Retrieving segment %d, duration: %s",
			i,
			c.Video.EncodedRepresentations[0].SegmentLocations[i].Duration))

		c.loadDashSegment(i)
	}

}

func (c *DASHClient) handleConnection(conn net.Conn) {
	log.Printf("Accepted connection from %s\n", conn.RemoteAddr())
	for {
		fmt.Println("Inside the client retrieval loop!!")
		segment, ok := <-c.ByteChan
		if !ok {
			break
		}
		if err := binary.Write(conn, binary.BigEndian, uint64(len(segment))); err != nil {
			c.ErrorChan <- err
		}
		_, err := io.Copy(conn, bytes.NewReader(segment))
		if err != nil {
			c.ErrorChan <- err
		}
	}
	//	forcefully exit;
	c.ErrorChan <- errors.New("viewing completed.")
}

func (c *DASHClient) acceptConnections(listener net.Listener) net.Conn {
	conn, err := listener.Accept()
	if err != nil {
		c.ErrorChan <- err
	}
	return conn
}

func (c *DASHClient) retrieveVideoManifest(videouid string) error {
	resp, err := c.HTTPS.GenericMethod(fmt.Sprintf("http://127.0.0.1:8080/manifest/%s", videouid))
	if err != nil {
		c.ErrorChan <- err
	}

	if resp.StatusCode == 404 {
		c.ErrorChan <- errors.New(fmt.Sprintf("unable to find %s", videouid))
	}

	var manifest []*types.HLSManifest
	err = json.Unmarshal(resp.Bytes, &manifest)
	if err != nil {
		return err
	}

	c.Video.EncodedRepresentations = manifest
	c.Video.VideoUID = videouid
	return nil
}

func (c *DASHClient) initialSequentialBufferHydration() int {
	segmentcount := 0
	start := time.Now()
	for c.PlaybackBuffer.GetTotalTimeLoaded() < c.PlaybackBuffer.Buffersize {

		log.Println(fmt.Sprintf(
			"buffer hydration in process. Retrieving segment %d, duration: %s",
			segmentcount,
			c.Video.EncodedRepresentations[0].SegmentLocations[segmentcount].Duration))

		c.loadDashSegment(segmentcount)
		segmentcount += 1
	}

	log.Println("hydration time: ", time.Since(start))
	return segmentcount
}

func triggerFFPlayChildProcess(address string) (*os.Process, error) {
	cmd := exec.Command("ffplay", fmt.Sprintf("tcp://localhost%s", address))
	err := cmd.Start()
	if err != nil {
		return nil, err
	}
	return cmd.Process, nil
}

func New(resolutions []string, byteChan chan []byte, errChan chan error) DASHClient {
	return DASHClient{
		PlaybackBuffer: playbackbuffer.New(),
		//SourceURL:            videouid,
		ErrorChan:            errChan,
		ByteChan:             byteChan,
		SupportedResolutions: resolutions,
		HTTPS:                https.HTTPS{},
	}
}
