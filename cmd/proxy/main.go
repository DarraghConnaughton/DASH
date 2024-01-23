package main

import (
	p "dash/internal/proxy"
	"dash/pkg/helper"
	"log"
	"os"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	// Shared error channel
	errChan := make(chan error, 1)

	proxy := p.New(errChan)
	go proxy.Start(":8889")
	proxy.LaunchNGINX()
	// Monitor all channels for errors.
	if err := helper.MonitorErrorChannel(errChan, true); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
