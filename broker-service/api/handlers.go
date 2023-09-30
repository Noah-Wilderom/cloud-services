package main

import (
	"context"
	"encoding/json"
	"github.com/Noah-Wilderom/cloud-services/broker/logs"
	"github.com/Noah-Wilderom/cloud-services/broker/queue"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
	"time"
)

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type JobRequestPayload struct {
	Id         string     `json:"id"`
	Payload    JobPayload `json:"payload"`
	ReservedAt *time.Time `json:"reserved_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

type JobPayload struct {
	Service string          `json:"service"`
	Data    json.RawMessage `json:"data"`
}

func (app *Config) LogItem(w http.ResponseWriter, r *http.Request) {
	var requestPayload LogPayload

	err := app.readJson(w, r, &requestPayload)
	if err != nil {
		_ = app.ErrorJson(w, err)
	}

	err = sendPayloadThroughGRPC("logger-service:5001", func(conn *grpc.ClientConn, ctx context.Context) {
		c := logs.NewLogServiceClient(conn)
		_, err = c.WriteLog(ctx, &logs.LogRequest{
			LogEntry: &logs.Log{
				Name: requestPayload.Name,
				Data: requestPayload.Data,
			},
		})
	})
	if err != nil {
		_ = app.ErrorJson(w, err)
		return
	}

	payload := JsonResponse{
		Error:   false,
		Message: "logged",
	}

	_ = app.WriteJson(w, http.StatusAccepted, payload)
}

func (app *Config) JobDispatch(w http.ResponseWriter, r *http.Request) {
	var requestPayload JobRequestPayload
	err := app.readJson(w, r, &requestPayload)
	if err != nil {
		_ = app.ErrorJson(w, err)
	}

	err = sendPayloadThroughGRPC("queue-listener:5003", func(conn *grpc.ClientConn, ctx context.Context) {
		c := queue.NewQueueListenerServiceClient(conn)
		_, err = c.InsertJob(ctx, &queue.Job{
			Payload: &queue.JobPayload{
				Service: requestPayload.Payload.Service,
				Data:    requestPayload.Payload.Data,
			},
		})
	})
	if err != nil {
		_ = app.ErrorJson(w, err)
		return
	}

	payload := JsonResponse{
		Error:   false,
		Message: "queue listener inserted",
	}

	_ = app.WriteJson(w, http.StatusAccepted, payload)
}

func sendPayloadThroughGRPC(target string, callback func(conn *grpc.ClientConn, ctx context.Context)) error {
	conn, err := grpc.Dial(
		target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return err
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	callback(conn, ctx)

	return nil
}
