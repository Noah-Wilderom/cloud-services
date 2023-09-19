package data

import (
	"context"
	"github.com/Noah-Wilderom/cloud-services/queue-worker/queue"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type Job struct {
	Id         string     `json:"id"`
	Payload    JobPayload `json:"payload"`
	ReservedAt *time.Time `json:"reserved_at"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

type JobPayload struct {
	Service string `json:"service"`
	Data    []byte `json:"data"`
}

func (j *Job) SetReserved() error {
	conn, err := grpc.Dial(
		"queue-listener:5003",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return err
	}
	defer conn.Close()

	c := queue.NewQueueListenerServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var reservedAtProto *timestamppb.Timestamp
	if j.ReservedAt == nil {
		reservedAtProto = nil
	} else {
		reservedAtProto = timestamppb.New(*j.ReservedAt)
	}

	createdAtProto := timestamppb.New(j.CreatedAt)
	updatedAtProto := timestamppb.New(j.UpdatedAt)
	cJob, err := c.SetJobReserved(ctx, &queue.Job{
		Id: j.Id,
		Payload: &queue.JobPayload{
			Service: j.Payload.Service,
			Data:    j.Payload.Data,
		},
		ReservedAt: reservedAtProto,
		UpdatedAt:  updatedAtProto,
		CreatedAt:  createdAtProto,
	})

	timeNow := time.Now()
	j.Id = cJob.GetId()
	j.Payload = JobPayload{
		Service: cJob.GetPayload().GetService(),
		Data:    cJob.GetPayload().GetData(),
	}
	j.ReservedAt = &timeNow
	j.UpdatedAt = cJob.GetUpdatedAt().AsTime()
	j.CreatedAt = cJob.GetCreatedAt().AsTime()

	return nil
}
