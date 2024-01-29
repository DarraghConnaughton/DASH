package monitor

import (
	"bytes"
	"dash/pkg/helper"
	"dash/pkg/types"
	"io"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

type Monitor struct {
	ErrorChan chan error
}

func (s *Monitor) Start(network string, bind string) {
	listener, err := net.Listen(network, bind)
	if err != nil {
		s.ErrorChan <- err
	}
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {

		}
	}(listener)

	log.Printf("%s server - listening on port %s.", network, bind)
	for {
		conn, err := listener.Accept()
		if err != nil {
			s.ErrorChan <- err
		}

		go rpc.ServeConn(conn)
	}
}

// HandleHeartBeat is a method of Monitor to handle incoming requests.
func (s *Monitor) HandleHeartBeat(hb *types.RPCHeartBeat, reply *int) error { // Perform some processing with the received data
	log.Printf("[%s statbot report] # goroutines: %d; AllocMemory: %d; TotalAlloc: %d; SysMem: %d\n",
		hb.UID,
		hb.NumberOfGoroutines,
		hb.AllocatedMemory,
		hb.TotalAlloc,
		hb.SysMem)

	resp, err := http.Post(
		"http://dash-opentsdb-1:4242/api/put",
		"application/json",
		bytes.NewBuffer([]byte(helper.FormatMetric(hb))))
	log.Println(resp)
	log.Println(err)
	log.Println("-------")
	if err != nil {
		s.ErrorChan <- err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			s.ErrorChan <- err
		}
	}(resp.Body)

	log.Printf("opentsdb response: status_code: %s ", resp.StatusCode)
	*reply = 1
	return nil
}

func New(errChan chan error) *Monitor {
	monitor := new(Monitor)

	err := rpc.Register(monitor)
	if err != nil {
		errChan <- err
	}

	return monitor
}
