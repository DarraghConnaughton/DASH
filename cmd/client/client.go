package client

import (
	"bytes"
	"dash/cmd/https"
	"dash/cmd/types"
	"dash/cmd/video"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
)

type VideoPlayer struct {
	ByteChan             chan []byte
	ErrorChan            chan error
	HTTPS                https.HTTPS
	SupportedResolutions []string
	SourceURL            string
	Video                video.Video
}

func (v *VideoPlayer) Watch() {
	err := v.Video.LoadManifests(v.SupportedResolutions)
	if err != nil {
		v.ErrorChan <- err
	}
	go v.gatherDASHSegments()
	go v.serveVideoViaTCP()
}

func (v *VideoPlayer) establishListener() net.Listener {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		v.ErrorChan <- err
	}
	return listener
}

func (v *VideoPlayer) launchVideoWithFFPlay(port string) *os.Process {
	process, err := triggerFFPlayChildProcess(port)
	if err != nil {
		v.ErrorChan <- err
	}
	return process
}

func (v *VideoPlayer) serveVideoViaTCP() {

	// Set up listener
	listener := v.establishListener()
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			v.ErrorChan <- err
		}
	}(listener)

	//=====
	//=====
	//Could add this to separate TCP class.

	// Launch video which reads from previous established TCP port
	port := strings.Split(listener.Addr().String(), "[::]")[1]
	process := v.launchVideoWithFFPlay(port)
	defer func(process *os.Process) {
		err := process.Kill()
		if err != nil {
			v.ErrorChan <- err
		}
	}(process)
	//=====
	//=====

	// Accept connection and proceed to handler function.
	conn := v.acceptConnections(listener)
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			v.ErrorChan <- err
		}
	}(conn)

	v.handleConnection(conn)

}

func (v *VideoPlayer) gatherDASHSegments() {
	resp, err := v.HTTPS.GenericMethod(fmt.Sprintf("http://127.0.0.1:8080/manifest/%s", v.SourceURL))
	if err != nil {
		v.ErrorChan <- err
	}

	var manifest []types.HLSManifest
	err = json.Unmarshal(resp.Bytes, &manifest)
	if err != nil {
		return
	}

	defer close(v.ByteChan)
	for i := 0; i < len(manifest[0].SegmentLocations); i++ {
		resp, err = v.HTTPS.GenericMethod(
			fmt.Sprintf("http://127.0.0.1:8080/hlsmanifest/%s/1080p/%d", v.SourceURL, i))
		if err != nil {
			v.ErrorChan <- err
		}
		v.ByteChan <- resp.Bytes
	}

}

func (v *VideoPlayer) handleConnection(conn net.Conn) {
	log.Printf("Accepted connection from %s\n", conn.RemoteAddr())
	for {
		segment, ok := <-v.ByteChan
		if !ok {
			break
		}
		if err := binary.Write(conn, binary.BigEndian, uint64(len(segment))); err != nil {
			v.ErrorChan <- err
		}
		_, err := io.Copy(conn, bytes.NewReader(segment))
		if err != nil {
			v.ErrorChan <- err
		}
	}
}

func (v *VideoPlayer) acceptConnections(listener net.Listener) net.Conn {
	conn, err := listener.Accept()
	if err != nil {
		v.ErrorChan <- err
	}
	return conn
}

func triggerFFPlayChildProcess(address string) (*os.Process, error) {
	cmd := exec.Command("ffplay", fmt.Sprintf("tcp://localhost%s", address))
	err := cmd.Start()
	if err != nil {
		return nil, err
	}
	return cmd.Process, nil
}

func New(videouid string, errChan chan error, resolutions []string, byteChan chan []byte) VideoPlayer {
	return VideoPlayer{
		SourceURL:            videouid,
		ErrorChan:            errChan,
		ByteChan:             byteChan,
		SupportedResolutions: resolutions,
		HTTPS:                https.HTTPS{},
	}
}
