package main

import (
	"dash/internal/dashserver"
	"dash/pkg/helper"
	"log"
	"os"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	// Shared error channel
	errChan := make(chan error, 1)

	// Read supported resolutions from configuration file.
	resolutions := []string{"1080p", "720p", "480p", "360p"}

	// Launch DASH servertr
	vs := dashserver.New("./data", resolutions, errChan)
	go vs.Start()

	// Monitor all channels for errors.
	if err := helper.MonitorErrorChannel(errChan, true); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
