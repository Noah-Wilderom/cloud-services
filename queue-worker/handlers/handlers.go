package handlers

import (
	"context"
	"encoding/json"
	"github.com/Noah-Wilderom/cloud-services/queue-worker/logs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func SendLog(data []byte) error {
	var requestPayload LogPayload

	err := json.Unmarshal(data, &requestPayload)
	if err != nil {
		return err
	}

	conn, err := grpc.Dial(
		"logger-service:5001",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return err
	}
	defer conn.Close()

	c := logs.NewLogServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err = c.WriteLog(ctx, &logs.LogRequest{
		LogEntry: &logs.Log{
			Name: requestPayload.Name,
			Data: requestPayload.Data,
		},
	})
	if err != nil {
		return err
	}

	return nil
}
