package main

import (
	"dash/internal/monitor"
	"dash/pkg/helper"
	"log"
	"os"
)

//// Monitor is the struct that implements the remote methods.
//type Monitor struct{}
//
//// HandleHeartBeat is a method of Monitor to handle incoming requests.
//func (s *Monitor) HandleHeartBeat(hb *types.RPCHeartBeat, reply *int) error { // Perform some processing with the received data
//	log.Printf("[%s statbot report] # goroutines: %d; AllocMemory: %d; TotalAlloc: %d; SysMem: %d\n",
//		hb.UID,
//		hb.NumberOfGoroutines,
//		hb.AllocatedMemory,
//		hb.TotalAlloc,
//		hb.SysMem)
//
//	*reply = 1
//
//	currentTimestamp := time.Now().Unix()
//
//	// Sample data with dynamic timestamp
//	data := fmt.Sprintf(`[{
//		"metric": "dash_service_monitor",
//		"timestamp": %d,
//		"value": %d,
//		"tags": {
//			"service": "%s",
//			"type": "active_goroutines"
//		}
//	}]`, currentTimestamp, hb.NumberOfGoroutines, hb.UID)
//
//	resp, err := http.Post("http://dash-opentsdb-1:4242/api/put", "application/json", bytes.NewBuffer([]byte(data)))
//	if err != nil {
//		fmt.Println("Error sending data to OpenTSDB:", err)
//		return err
//	}
//	defer resp.Body.Close()
//
//	// Check the response status
//	if resp.StatusCode == http.StatusOK {
//		fmt.Println("Data sent successfully.")
//	} else {
//		fmt.Printf("Error: %s\n", resp.Status)
//	}
//	return nil
//}

func main() {
	errChan := make(chan error, 1)
	m := monitor.New(errChan)
	m.Start("tcp", ":1234")

	if err := helper.MonitorErrorChannel(errChan, false); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
