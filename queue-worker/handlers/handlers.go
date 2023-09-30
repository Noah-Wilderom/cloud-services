package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Noah-Wilderom/cloud-services/queue-worker/logs"
	"github.com/Noah-Wilderom/cloud-services/queue-worker/projects"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type ProjectPayload struct {
	Project Project `json:"project"`
	Action  string  `json:"action"`
}

type Project struct {
	Id           string  `json:"id"`
	Domain       Domain  `json:"domain"`
	Status       string  `json:"status"`
	Name         string  `json:"name"`
	Subdomain    string  `json:"subdomain"`
	Docker       bool    `json:"docker"`
	SFTP         *SFTP   `json:"sftp"`
	PreviewImage *string `json:"preview_image"`
	User         *string `json:"user"`
	SshKeyPath   *string `json:"ssh_key_path"`
	FilesPath    *string `json:"files_path"`
	Git          *Git    `json:"git"`
	Monitoring   bool    `json:"monitoring"`
}

type SFTP struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Path     string `json:"path"`
}

type Git struct {
	Repository string `json:"repository"`
	WebhookUrl string `json:"webhook_url"`
	SshKeyPath string `json:"ssh_key_path"`
}

type Domain struct {
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
		log.Println("Error on unmarshal")
		return err
	}

	if requestPayload.Project.SFTP == nil {
		requestPayload.Project.SFTP = &SFTP{}
	}

	if requestPayload.Project.Git == nil {
		requestPayload.Project.Git = &Git{}
	}

	err = sendPayloadThroughGRPC("project-service:5004", func(conn *grpc.ClientConn, ctx context.Context) {
		c := projects.NewProjectServiceClient(conn)

		resp, _ := c.HandleJob(ctx, &projects.ProjectRequest{
			Project: &projects.Project{
				Id: requestPayload.Project.Id,
				Domain: &projects.Domain{
					Id:     requestPayload.Project.Domain.Id,
					Domain: requestPayload.Project.Domain.Domain,
				},
				Status:    requestPayload.Project.Status,
				Name:      requestPayload.Project.Name,
				Subdomain: requestPayload.Project.Subdomain,
				Docker:    requestPayload.Project.Docker,
				SFTP: &projects.SFTP{
					User:     requestPayload.Project.SFTP.User,
					Password: requestPayload.Project.SFTP.Password,
					Path:     requestPayload.Project.SFTP.Path,
				},
				PreviewImage: *requestPayload.Project.PreviewImage,
				User:         *requestPayload.Project.User,
				SshKeyPath:   *requestPayload.Project.SshKeyPath,
				FilesPath:    *requestPayload.Project.FilesPath,
				Git: &projects.Git{
					Repository: requestPayload.Project.Git.Repository,
					WebhookUrl: requestPayload.Project.Git.WebhookUrl,
					SshKeyPath: requestPayload.Project.Git.SshKeyPath,
				},
				Monitoring: requestPayload.Project.Monitoring,
			},
			Action: requestPayload.Action,
		})

		if resp.Error {
			fmt.Println("error on handlejob", resp.ErrorPayload)
		}
	})
	if err != nil {
		log.Println("Error on sending to grpc")

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
