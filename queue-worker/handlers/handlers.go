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

type ProjectPayload struct {
	Id           string        `json:"id"`
	Domain       DomainPayload `json:"domain"`
	Status       string        `json:"status"`
	Name         string        `json:"name"`
	Subdomain    *string       `json:"subdomain"`
	Docker       bool          `json:"docker"`
	Stack        *string       `json:"stack"`
	SFTP         *string       `json:"sftp"`
	PreviewImage *string       `json:"preview_image"`
	User         *string       `json:"user"`
	SshKeyPath   *string       `json:"ssh_key_path"`
	FilesPath    *string       `json:"files_path"`
	Git          bool          `json:"git"`
	Monitoring   bool          `json:"monitoring"`
}

type DomainPayload struct {
	Id     string `json:"id"`
	Domain string `json:"domain"`
}

func SendLog(data []byte) error {
	var requestPayload LogPayload

	err := json.Unmarshal(data, &requestPayload)
	if err != nil {
		return err
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
		return err
	}

	return nil
}

func SendProjectJob(data []byte) error {
	var requestPayload ProjectPayload

	err := json.Unmarshal(data, &requestPayload)
	if err != nil {
		return err
	}

	err = sendPayloadThroughGRPC("project-service:5004", func(conn *grpc.ClientConn, ctx context.Context) {
		//
	})
	if err != nil {
		return err
	}

	return nil
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
