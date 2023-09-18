package main

import (
	"context"
	"fmt"
	"github.com/Noah-Wilderom/cloud-services/logger/data"
	"github.com/Noah-Wilderom/cloud-services/logger/logs"
	"google.golang.org/grpc"
	"log"
	"net"
)

type LogServer struct {
	logs.UnimplementedLogServiceServer
	Models data.Models
}

func (l *LogServer) WriteLog(ctx context.Context, req *logs.LogRequest) (*logs.LogResponse, error) {
	input := req.GetLogEntry()

	logEntry := data.LogEntry{
		Name: input.Name,
		Data: input.Data,
	}

	err := l.Models.LogEntry.Insert(logEntry)
	if err != nil {
		res := &logs.LogResponse{
			Result: "Failed",
		}
		return res, err
	}

	res := &logs.LogResponse{
		Result: "Logged",
	}

	return res, nil
}

func (app *Config) gRPCListen() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", grpcPort))
	if err != nil {
		log.Fatalln("Failed to listen for gRPC:", err)
	}

	s := grpc.NewServer()

	logs.RegisterLogServiceServer(s, &LogServer{Models: app.Models})
	log.Println("gRPC Server started on port", grpcPort)

	if err = s.Serve(lis); err != nil {
		log.Fatalln("Failed to listen for gRPC:", err)
	}
}
