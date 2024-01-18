package main

import (
	"dash/cmd/client"
	"dash/cmd/fileserver"
	"log"
	"os"
	"time"
)

func MonitorErrorChannel(errChan chan error, hardfail bool) error {
	for {
		select {
		case err := <-errChan:
			if err != nil {
				log.Println("[-]received an error from the goroutine:", err)
				if hardfail {
					log.Println("[-] hard fail mode enabled, exiting main goroutine.")
					return err
				}
			}
		}
	}
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	// Shared error channel
	errChan := make(chan error, 1)
	// Channel for sharing bytes better DASH retriever and function feeding ffplay.
	byteChan := make(chan []byte)

	// Launch DASH server
	vs := fileserver.New("./data", []string{"1080p", "720p", "480p", "360p"}, errChan)
	go vs.Start()

	// Give time for the fileserver to start.
	time.Sleep(1 * time.Second)

	v := client.New("video1", errChan, []string{"1080p", "720p", "480p", "360p"}, byteChan)
	v.Watch()

	// Monitor all channels for errors.
	if err := MonitorErrorChannel(errChan, true); err != nil {
		log.Println(err)
		os.Exit(1)
	}

}
