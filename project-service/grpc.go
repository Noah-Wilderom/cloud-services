package main

import (
	"fmt"
	"github.com/Noah-Wilderom/cloud-services/shared-grpc/projects"
	"google.golang.org/grpc"
	"log"
	"net"
)

type Config struct{}

func (app *Config) gRPCListen() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", gRPCPort))
	if err != nil {
		log.Fatalln("Failed to listen for gRPC:", err)
	}

	s := grpc.NewServer()

	projects.RegisterQueueListenerServiceServer(s, &QueueListenerServer{})
	log.Println("gRPC Server started on port", gRPCPort)

	if err = s.Serve(lis); err != nil {
		log.Fatalln("Failed to listen for gRPC:", err)
	}
}
