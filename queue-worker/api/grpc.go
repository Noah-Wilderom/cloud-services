package main

import (
	"context"
	"fmt"
	"github.com/Noah-Wilderom/cloud-services/queue-worker/data"
	"github.com/Noah-Wilderom/cloud-services/queue-worker/handlers"
	"github.com/Noah-Wilderom/cloud-services/queue-worker/queue"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

type QueueWorkerServer struct {
	queue.UnimplementedQueueWorkerServiceServer
}

func (q *QueueWorkerServer) HandleJob(ctx context.Context, req *queue.JobRequest) (*queue.JobResponse, error) {
	fmt.Println("Job received!")
	input := req.GetJob()
	fmt.Println(req.GetJob().Id)

	timeNow := time.Now()
	job := data.Job{
		Id: input.GetId(),
		Payload: data.JobPayload{
			Service: input.GetPayload().GetService(),
			Data:    input.GetPayload().GetData(),
		},
		ReservedAt: &timeNow,
		UpdatedAt:  input.GetUpdatedAt().AsTime(),
		CreatedAt:  input.GetCreatedAt().AsTime(),
	}

	if !input.GetReservedAt().IsValid() {
		fmt.Println("Job is already reserved ERROR")
		return &queue.JobResponse{
			Error:        true,
			ErrorPayload: "Job is already reserved.",
		}, nil
	}

	if err := job.SetReserved(); err != nil {
		fmt.Println("error on setreserved", err)
	}

	switch job.Payload.Service {
	case "logger":
		fmt.Println("Logger job received")
		err := handlers.SendLog(job.Payload.Data)
		if err != nil {
			fmt.Println("Job is already reserved2 ERROR")
			return &queue.JobResponse{
				Error:        true,
				ErrorPayload: "Job is already reserved.",
			}, nil
		}
	case "project":
		fmt.Println("Project job received", input.GetId())
		err := handlers.SendProjectJob(job.Payload.Data)
		if err != nil {
			return &queue.JobResponse{
				Error:        true,
				ErrorPayload: err.Error(),
			}, nil
		}
	}

	return &queue.JobResponse{
		Error:        false,
		ErrorPayload: "",
	}, nil
}

func (app *Config) gRPCListen() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", gRPCPort))
	if err != nil {
		log.Fatalln("Failed to listen for gRPC:", err)
	}

	s := grpc.NewServer()

	queue.RegisterQueueWorkerServiceServer(s, &QueueWorkerServer{})
	log.Println("gRPC Server started on port", gRPCPort)

	if err = s.Serve(lis); err != nil {
		log.Fatalln("Failed to listen for gRPC:", err)
	}
}
