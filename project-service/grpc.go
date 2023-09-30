package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/Noah-Wilderom/cloud-services/project-service/handlers"
	"github.com/Noah-Wilderom/cloud-services/project-service/projects"
	"google.golang.org/grpc"
	"log"
	"net"
)

type Config struct{}

type ProjectServiceServer struct {
	projects.UnimplementedProjectServiceServer
}

func (app *Config) gRPCListen() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", gRPCPort))
	if err != nil {
		log.Fatalln("Failed to listen for gRPC:", err)
	}

	s := grpc.NewServer()

	projects.RegisterProjectServiceServer(s, &ProjectServiceServer{})
	log.Println("gRPC Server started on port", gRPCPort)

	if err = s.Serve(lis); err != nil {
		log.Fatalln("Failed to listen for gRPC:", err)
	}
}

func (q *ProjectServiceServer) HandleJob(ctx context.Context, request *projects.ProjectRequest) (*projects.ProjectResponse, error) {
	switch request.Action {
	case "provision":
		err := handlers.ProvisionProject(request.Project)
		if err != nil {
			log.Println("error on provisioning project", err)
			return &projects.ProjectResponse{
				Error:        true,
				ErrorPayload: err.Error(),
			}, nil
		}
	default:
		return &projects.ProjectResponse{
			Error:        true,
			ErrorPayload: errors.New("no valid action found").Error(),
		}, nil
	}

	return &projects.ProjectResponse{
		Error:        false,
		ErrorPayload: "",
	}, nil

}
