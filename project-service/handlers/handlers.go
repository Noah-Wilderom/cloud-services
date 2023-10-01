package handlers

import (
	"fmt"
	"github.com/Noah-Wilderom/cloud-services/project-service/api"
	"github.com/Noah-Wilderom/cloud-services/project-service/projects"
	"log"
	"os"
)

func ProvisionProject(project *projects.Project) error {
	log.Println("PROVISIONPROJECT", project.GetId())
	conn := api.NewApi()
	_, err := conn.UpdateJobStatus(project.GetId(), "starting")
	if err != nil {
		return err
	}

	subdomain := project.GetSubdomain()
	if subdomain != "" {
		subdomain += "."
	}
	dir := fmt.Sprintf("/var/www/%s%s", subdomain, project.GetDomain().Domain)

	err = os.Mkdir(dir, 644)
	if err != nil {
		return err
	}
	_, err = conn.UpdateFilesPath(project.GetId(), dir)
	if err != nil {
		return err
	}

	return nil
}
