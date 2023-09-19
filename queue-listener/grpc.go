package main

import (
	"context"
	"fmt"
	"github.com/Noah-Wilderom/cloud-services/queue-listener/data"
	"github.com/Noah-Wilderom/cloud-services/queue-listener/queue"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"net"
)

type QueueListenerServer struct {
	queue.UnimplementedQueueListenerServiceServer
}

func (q *QueueListenerServer) SetJobReserved(ctx context.Context, cJob *queue.Job) (*queue.Job, error) {
	job, err := data.SetReservedById(cJob.GetId())
	if err != nil {
		log.Println("error refreshing database record via grpc", err)
	}
	var reservedAt timestamppb.Timestamp
	if job.ReservedAt == nil {
		reservedAt = timestamppb.Timestamp{}
	} else {
		reservedAt = *timestamppb.New(*job.ReservedAt)
	}

	return &queue.Job{
		Id: job.Id,
		Payload: &queue.JobPayload{
			Service: job.Payload.Service,
			Data:    job.Payload.Data,
		},
		ReservedAt: &reservedAt,
		UpdatedAt:  timestamppb.New(job.UpdatedAt),
		CreatedAt:  timestamppb.New(job.CreatedAt),
	}, nil

}

func (q *QueueListenerServer) RefreshJob(ctx context.Context, cJob *queue.Job) (*queue.Job, error) {
	job, err := data.RefreshById(cJob.GetId())
	if err != nil {
		log.Println("error refreshing database record via grpc", err)
	}

	var reservedAt *timestamppb.Timestamp
	if job.ReservedAt.IsZero() {
		reservedAt = &timestamppb.Timestamp{}
	} else {
		reservedAt = timestamppb.New(*job.ReservedAt)
	}

	return &queue.Job{
		Id: job.Id,
		Payload: &queue.JobPayload{
			Service: job.Payload.Service,
			Data:    job.Payload.Data,
		},
		ReservedAt: reservedAt,
		UpdatedAt:  timestamppb.New(job.UpdatedAt),
		CreatedAt:  timestamppb.New(job.CreatedAt),
	}, nil

}

func (app *Config) gRPCListen() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", gRPCPort))
	if err != nil {
		log.Fatalln("Failed to listen for gRPC:", err)
	}

	s := grpc.NewServer()

	queue.RegisterQueueListenerServiceServer(s, &QueueListenerServer{})
	log.Println("gRPC Server started on port", gRPCPort)

	if err = s.Serve(lis); err != nil {
		log.Fatalln("Failed to listen for gRPC:", err)
	}
}
