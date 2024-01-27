package main

import (
	"dash/internal/dashserver"
	"dash/pkg/helper"
	"dash/pkg/statistics"
	"log"
	"os"
	"time"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	// Shared error channel

	// ======= Ideally use flags for each of these...
	// ======= Ideally use flags for each of these...
	// ======= Ideally use flags for each of these...
	errChan := make(chan error, 1)
	statKS := make(chan int, 1)
	statTicker := time.NewTicker(5 * time.Second)
	// ======= Ideally use flags for each of these...
	// ======= Ideally use flags for each of these...
	// ======= Ideally use flags for each of these...
	statbot := statistics.New("dashserver", "monitor:1234", *statTicker, statKS, errChan)
	go statbot.Start()

	// Launch DASH server
	vs := dashserver.New("./data", []string{"1080p", "720p", "480p", "360p"}, errChan)
	go vs.Start(":8080")

	// Monitor all channels for errors.
	if err := helper.MonitorErrorChannel(errChan, true); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
