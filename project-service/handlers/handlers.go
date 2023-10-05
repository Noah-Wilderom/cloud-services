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

	if _, err = os.Stat(dir); err != nil {
		err = helpers.RemoveAllFilesInDirectory(dir)
	}

	err = os.MkdirAll(dir, 644)
	if err != nil {
		return err
	}

	_, err = conn.UpdateFilesPath(project.GetId(), dir)
	if err != nil {
		return err
	}
	fmt.Println(strings.ToUpper(project.GetStack()), strings.ToUpper(project.GetStack()) == "PHP")

	if strings.ToUpper(project.GetStack()) == "PHP" || true { // TODO: bug fix stack not applied in create project form

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
		confPath := fmt.Sprintf("/etc/nginx/sites-available/%s", fullDomain)
		err := helpers.ReplaceStubVariables("/templates/nginx/laravel/http.conf", confPath, vars)
		if err != nil {
			fmt.Println(err)
			return err
		}

		symLink := fmt.Sprintf("/etc/nginx/sites-enabled/%s", fullDomain)
		if _, err = os.Stat(symLink); err != nil {
			_ = os.Remove(symLink)
		}
		err = os.Symlink(confPath, symLink)
		if err != nil {
			fmt.Println(err)
		}

	}

	image := helpers.NavigateAndTakeScreenshot(fmt.Sprintf("http://%s", fullDomain))
	err = conn.UpdateScreenshot(project.GetId(), image)

	if len(project.GetGit().Repository) > 1 {
		fmt.Println("Git enabled")
	}

	_, err = conn.UpdateJobStatus(project.GetId(), "running")
	if err != nil {
		return err
	}

	return nil
}
