package main

import (
	"dash/internal/monitor"
	"dash/pkg/helper"
	"log"
	"os"
)

func main() {
	errChan := make(chan error, 1)
	m := monitor.New(errChan)
	m.Start("tcp", ":1234")

	if err := helper.MonitorErrorChannel(errChan, false); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
