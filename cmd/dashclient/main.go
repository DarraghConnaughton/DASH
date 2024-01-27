package main

import (
	"dash/internal/dashclient"
	"dash/pkg/helper"
	"dash/pkg/statistics"
	"log"
	"os"
	"time"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	// ======= Ideally use flags for each of these...
	// ======= Ideally use flags for each of these...
	// ======= Ideally use flags for each of these...
	errChan := make(chan error, 1)
	statKS := make(chan int, 1)
	statTicker := time.NewTicker(5 * time.Second)
	// ======= Ideally use flags for each of these...
	// ======= Ideally use flags for each of these...
	// ======= Ideally use flags for each of these...
	statbot := statistics.New("dashclient", "localhost:1234", *statTicker, statKS, errChan)
	go statbot.Start()

	//=============
	//=============
	//=============
	//=============
	// Channel for sharing bytes better DASH retriever and function feeding ffplay.

	v := dashclient.New([]string{"1080p", "720p", "480p", "360p"}, make(chan []byte, 5), errChan)

	//i := 0
	//for i < 5 {
	v.Watch("video2")
	// Monitor all channels for errors.
	if err := helper.MonitorErrorChannel(errChan, true); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	//	i++
	//}
	//
	//=============
	//=============
	//=============
	//=============

}
