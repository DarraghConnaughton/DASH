package main

import (
	"dash/internal/proxies/forwardproxy"
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
	statbot := statistics.New("forwardproxy", "monitor:1234", *statTicker, statKS, errChan)
	go statbot.Start()

	proxy := forwardproxy.New(errChan)
	proxy.Start(":8889", proxy.GetRoutes())

	if err := helper.MonitorErrorChannel(errChan, true); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
