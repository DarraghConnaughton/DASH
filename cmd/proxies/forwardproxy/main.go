package main

import (
	"dash/internal/proxies/forwardproxy"
	"dash/pkg/helper"
	"log"
	"os"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	// Shared error channel
	errChan := make(chan error, 1)

	proxy := forwardproxy.New(errChan)
	proxy.Start(":8889", proxy.GetRoutes())

	if err := helper.MonitorErrorChannel(errChan, true); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
