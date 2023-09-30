package handlers

import (
	"github.com/Noah-Wilderom/cloud-services/project-service/api"
	"github.com/Noah-Wilderom/cloud-services/project-service/projects"
	"log"
)

func ProvisionProject(project *projects.Project) error {
	log.Println("PROVISIONPROJECT")
	conn := api.NewApi()
	_, err := conn.UpdateJobStatus(project.Id, "starting")
	if err != nil {
		return err
	}

	return nil
}
