package main

import (
	p "dash/internal/proxies/controlproxy"
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
	statbot := statistics.New("controlproxy", "monitor:1234", *statTicker, statKS, errChan)
	go statbot.Start()

	proxy := p.New(errChan)
	proxy.Start(":8887", proxy.GetRoutes())
	//proxy.LaunchNGINX()
	// Monitor all channels for errors.
	if err := helper.MonitorErrorChannel(errChan, true); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
