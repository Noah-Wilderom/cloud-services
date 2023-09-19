package main

import "sync"

const (
	gRPCPort = 5002
)

type Config struct {
}

func main() {

	app := Config{}

	// Create a WaitGroup
	var wg sync.WaitGroup

	// Increment the WaitGroup to indicate a goroutine is starting
	wg.Add(1)

	// Start your gRPC server in a goroutine
	go func() {
		defer wg.Done() // Decrement the WaitGroup when done
		app.gRPCListen()
	}()

	// Wait for the gRPC server to start
	wg.Wait()
}
