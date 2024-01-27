package dashclient

import (
	"bytes"
	"dash/pkg/helper"
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
	HTTPS                https.HTTP
	PlaybackBuffer       playbackbuffer.PlaybackBuffer
	SupportedResolutions []string
	Video                video.Video
	PreviousRTT          int64
}

func (c *DASHClient) Watch(videouid string) {
	c.retrieveVideoManifest(videouid)
	startsegment := c.PlaybackBufferHydration()
	go c.serveVideoViaTCP()
	c.gatherDASHSegments(startsegment)
}

// ============ FFPlayHandler
// ============ FFPlayHandler
// ============ FFPlayHandler
// ============ FFPlayHandler
// ============ FFPlayHandler
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
		log.Println("do we make it ehre??? ")
		log.Println("do we make it ehre??? ")
		log.Println("ideal place to cleanup ")
		log.Println("ideal place to cleanup ")
		log.Println("ideal place to cleanup ")
		err := process.Kill()
		if err != nil {
			c.ErrorChan <- err
		}

		helper.WriteCSV(c.Video.VideoUID, c.gatherNetworkTraces())
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

func (c *DASHClient) handleConnection(conn net.Conn) {
	log.Printf("Accepted connection from %s\n", conn.RemoteAddr())

	for {
		select {
		case segment, ok := <-c.ByteChan:
			if !ok {
				return
			}
			if err := binary.Write(conn, binary.BigEndian, uint64(len(segment))); err != nil {
				c.ErrorChan <- err
				return
			}
			_, err := io.Copy(conn, bytes.NewReader(segment))
			if err != nil {
				c.ErrorChan <- err
				return
			}
		}
	}
}

func (c *DASHClient) acceptConnections(listener net.Listener) net.Conn {
	conn, err := listener.Accept()
	if err != nil {
		c.ErrorChan <- err
	}
	return conn
}

func triggerFFPlayChildProcess(address string) (*os.Process, error) {
	cmd := exec.Command("ffplay", fmt.Sprintf("tcp://localhost%s", address))
	//cmd := exec.Command("ffplay", fmt.Sprintf("tcp://localhost%s", address))
	err := cmd.Start()
	if err != nil {
		return nil, err
	}
	return cmd.Process, nil
}

//============ FFPlayHandler
//============ FFPlayHandler
//============ FFPlayHandler
//============ FFPlayHandler
//============ FFPlayHandler

func (c *DASHClient) determineQuality() string {
	quality := "1080p"
	switch {
	case c.PreviousRTT >= int64(150*time.Millisecond) && c.PreviousRTT < int64(250*time.Millisecond):
		quality = "720p"
	case c.PreviousRTT >= int64(250*time.Millisecond) && c.PreviousRTT < int64(500*time.Millisecond):
		quality = "480p"
	case c.PreviousRTT >= int64(500*time.Millisecond):
		quality = "380p"
	}
	return quality
}

func (c *DASHClient) loadDashSegment(segmentindex int) {
	quality := c.determineQuality()
	url := fmt.Sprintf(
		"http://127.0.0.1:8888/delay/hlsmanifest/%s/%s/%d",
		c.Video.VideoUID,
		quality,
		segmentindex)

	resp, err := c.HTTPS.GenericMethod(url)
	if err != nil {
		c.ErrorChan <- err
	}

	c.PreviousRTT = int64(resp.RTT)

	respParts := bytes.Split(resp.Bytes, []byte("!!thisIsTheDelimiter!!"))
	if len(respParts) <= 1 {
		c.ErrorChan <- errors.New("expected noise, but encountered none.")
	}

	log.Println(fmt.Sprintf(
		"[/] RTT: %d;  [/] quality: %s; [/] bytes retrieved: [%d];",
		c.PreviousRTT,
		quality,
		len(resp.Bytes)))

	log.Println("_________________________")
	log.Println("video data:: --> ", len(respParts[0]))
	log.Println("noise     :: --> ", len(respParts[1]))
	log.Println("_________________________")

	c.ByteChan <- respParts[0]
	c.PlaybackBuffer.UpdateTotalTimeLoaded(
		c.Video.EncodedRepresentations[0].SegmentLocations[segmentindex].Duration,
		c.ErrorChan)
}

func (c *DASHClient) gatherDASHSegments(startsegment int) {

	defer close(c.ByteChan)
	for i := startsegment; i < len(c.Video.EncodedRepresentations[0].SegmentLocations); i++ {
		for c.PlaybackBuffer.GetTotalTimeLoaded()-time.Since(c.PlaybackBuffer.StartTime) > c.PlaybackBuffer.Buffersize {
			time.Sleep(100 * time.Millisecond)
		}

		log.Println(fmt.Sprintf(
			"buffer hydration in process. Retrieving segment %d, duration: %s",
			i,
			c.Video.EncodedRepresentations[0].SegmentLocations[i].Duration))

		c.loadDashSegment(i)
	}

}

func (c *DASHClient) retrieveVideoManifest(videouid string) {
	log.Println("we make it here?", videouid)
	resp, err := c.HTTPS.GenericMethod(fmt.Sprintf("http://127.0.0.1:8888/delay/manifest/%s", videouid))
	log.Println(resp)
	log.Println(err)
	if err != nil {
		log.Println(err)
		c.ErrorChan <- err
	}
	fmt.Println("here?")
	fmt.Println("here?", resp)

	if resp.StatusCode == 404 || resp.StatusCode == 502 {
		c.ErrorChan <- errors.New(fmt.Sprintf("unable to find %s", videouid))
	}

	var manifest []*types.HLSManifest
	err = json.Unmarshal(resp.Bytes, &manifest)
	if err != nil {
		c.ErrorChan <- err
	}

	c.Video.EncodedRepresentations = manifest
	c.Video.VideoUID = videouid
}

func (c *DASHClient) PlaybackBufferHydration() int {
	segmentcount := 0
	for c.PlaybackBuffer.GetTotalTimeLoaded() < c.PlaybackBuffer.Buffersize {
		log.Println(fmt.Sprintf(
			"[/] buffer hydration in process. Retrieving segment %d, duration: %s",
			segmentcount,
			c.Video.EncodedRepresentations[0].SegmentLocations[segmentcount].Duration))

		c.loadDashSegment(segmentcount)
		segmentcount += 1
	}
	return segmentcount
}

func (c *DASHClient) gatherNetworkTraces() []types.NetworkTraceData {
	dataCapture, err := c.HTTPS.GenericMethod("http://127.0.0.1:8889/process")
	if err != nil {
		c.ErrorChan <- err
	}
	var traces []types.NetworkTraceData
	err = json.Unmarshal(dataCapture.Bytes, &traces)
	if err != nil {
		c.ErrorChan <- err
	}
	return traces
}

func New(resolutions []string, byteChan chan []byte, errChan chan error) DASHClient {
	return DASHClient{
		PlaybackBuffer:       playbackbuffer.New(),
		ErrorChan:            errChan,
		ByteChan:             byteChan,
		SupportedResolutions: resolutions,
		HTTPS:                https.HTTP{},
	}
}
