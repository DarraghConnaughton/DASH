package main

import (
	"dash/internal/dashclient"
	"dash/pkg/helper"
	"log"
	"os"
	"time"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	// Shared error channel
	errChan := make(chan error, 1)
	// Channel for sharing bytes better DASH retriever and function feeding ffplay.
	byteChan := make(chan []byte, 50)

	resolutions := []string{"1080p", "720p", "480p", "360p"}

	// Give time for the dashserver to start.
	time.Sleep(10 * time.Second)

	v := dashclient.New(resolutions, byteChan, errChan)
	v.Watch("video2")

	// Monitor all channels for errors.
	if err := helper.MonitorErrorChannel(errChan, true); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
