package main

import (
	p "dash/internal/proxies/controlproxy"
	"dash/pkg/helper"
	"log"
	"os"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	// Shared error channel
	errChan := make(chan error, 1)

	proxy := p.New(errChan)
	proxy.Start(":8887", proxy.GetRoutes())
	//proxy.LaunchNGINX()
	// Monitor all channels for errors.
	if err := helper.MonitorErrorChannel(errChan, true); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
