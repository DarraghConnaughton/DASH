package statistics

import (
	"dash/pkg/types"
	"flag"
	"fmt"
	"log"
	"net/rpc"
	"runtime"
	"time"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

type StatBot struct {
	KillSwitch    chan int
	AccessAttempt int
	ErrorChan     chan error
	ServiceUID    string
	Ticker        time.Ticker
	MonitorHost   string
	RPCClient     *rpc.Client
}

func (sb *StatBot) Start() {

	client, err := rpc.Dial("tcp", sb.MonitorHost)
	if err != nil {
		log.Println("Error connecting to RPC server:", err)
		delay := time.Duration(1<<uint(sb.AccessAttempt)) * time.Second
		time.Sleep(delay)
		sb.AccessAttempt += 1
		sb.Start()
	}

	defer func(client *rpc.Client) {
		err := client.Close()
		if err != nil {
			log.Println("[-] failed to establish connection with monitor service.")
		}
	}(client)

	for {
		select {
		case <-sb.Ticker.C:
			var memStats runtime.MemStats
			runtime.ReadMemStats(&memStats)
			hb := types.RPCHeartBeat{
				UID:                sb.ServiceUID,
				NumberOfGoroutines: runtime.NumGoroutine(),
				AllocatedMemory:    memStats.Alloc,
				TotalAlloc:         memStats.TotalAlloc,
				SysMem:             memStats.Sys,
			}
			var response int
			err = client.Call("Monitor.HandleHeartBeat", hb, &response)
			if err != nil {
				log.Printf("[statbot] error communicating with monitor server: %s. Attempting to establish a new connection", err)
				sb.Start()
			} else {
				log.Println("[statbot] heartbeat.")
			}
		case c := <-sb.KillSwitch:
			log.Println(fmt.Sprintf("[statbot] shutting down. exit code: %d", c))
			return
		}
	}

}

func New(uid string, mmgthost string, ticker time.Ticker, ks chan int, errChan chan error) StatBot {
	return StatBot{
		Ticker:      ticker,
		KillSwitch:  ks,
		ServiceUID:  uid,
		MonitorHost: mmgthost,
		ErrorChan:   errChan,
	}
}
