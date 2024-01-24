package main

import (
	"dash/internal/dashclient"
	"dash/pkg/helper"
	"log"
	"os"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	// Shared error channel
	errChan := make(chan error, 1)
	// Channel for sharing bytes better DASH retriever and function feeding ffplay.
	byteChan := make(chan []byte, 5)

	v := dashclient.New([]string{"1080p", "720p", "480p", "360p"}, byteChan, errChan)
	v.Watch("video2")

	// Monitor all channels for errors.
	if err := helper.MonitorErrorChannel(errChan, true); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
