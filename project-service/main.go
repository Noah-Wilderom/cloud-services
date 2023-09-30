package main

import (
	"sync"
)

const (
	gRPCPort = 5004
)

func main() {

	app := Config{}

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		defer wg.Done()
		app.gRPCListen()
	}()

	wg.Wait()
}
