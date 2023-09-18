package main

import (
	"context"
	"github.com/Noah-Wilderom/cloud-services/broker/logs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
	"time"
)

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) LogItem(w http.ResponseWriter, r *http.Request) {
	var requestPayload LogPayload

	err := app.readJson(w, r, &requestPayload)
	if err != nil {
		_ = app.ErrorJson(w, err)
	}

	conn, err := grpc.Dial(
		"logger-service:5001",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		_ = app.ErrorJson(w, err)
		return
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
		_ = app.ErrorJson(w, err)
		return
	}

	payload := JsonResponse{
		Error:   false,
		Message: "logged",
	}

	_ = app.WriteJson(w, http.StatusAccepted, payload)
}
