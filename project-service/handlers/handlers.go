package handlers

import (
	"fmt"
	"github.com/Noah-Wilderom/cloud-services/project-service/api"
	"github.com/Noah-Wilderom/cloud-services/project-service/helpers"
	"github.com/Noah-Wilderom/cloud-services/project-service/projects"
	"log"
	"os"
	"strings"
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
	fullDomain := fmt.Sprintf("%s%s", subdomain, project.GetDomain().Domain)
	dir := fmt.Sprintf("/var/www/%s", fullDomain)

	err = os.Mkdir(dir, 644)
	if err != nil {
		return err
	}
	_, err = conn.UpdateFilesPath(project.GetId(), dir)
	if err != nil {
		return err
	}
	fmt.Println(strings.ToUpper(project.GetStack()), strings.ToUpper(project.GetStack()) == "PHP")

	if strings.ToUpper(project.GetStack()) == "PHP" || true { // TODO: bug fix stack not applied in create project form
		//files, err := helpers.ReadTemplateFiles("nginx/laravel")
		//if err != nil {
		//	return err
		//}

		//for _, file := range files {
		//	fmt.Println(file.Name(), strings.HasSuffix(file.Name(), "http.conf"))
		//	if strings.HasSuffix(file.Name(), "http.conf") {
		domainServerName := fullDomain
		if len(strings.Split(fullDomain, ".")) < 3 {
			domainServerName = project.GetDomain().Domain
		}

		vars := map[string]string{
			"DOMAIN_SERVER_NAME": domainServerName,
			"DOMAIN":             fullDomain,
			"FILES_PATH":         dir,
			"PHP_VERSION":        "php8.2",
		} // TODO: PHP version moet configurable zijn
		err := helpers.ReplaceStubVariables("/templates/nginx/laravel/http.conf", fmt.Sprintf("/etc/nginx/sites-available/%s", fullDomain), vars)
		if err != nil {
			fmt.Println(err)
			return err
		}

		err = os.Symlink(fmt.Sprintf("/etc/nginx/sites-available/%s", fullDomain), fmt.Sprintf("/etc/nginx/sites-enabled/%s", fullDomain))
		if err != nil {
			fmt.Println(err)
			return err
		}
		//	}
		//}
	}

	_, err = conn.UpdateJobStatus(project.GetId(), "running")
	if err != nil {
		return err
	}

	return nil
}

//func setupNginxFiles() {
//
//}
