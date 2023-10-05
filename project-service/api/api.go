package api

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

const (
	//apiUrl = "http://localhost/api"
	//apiUrl = "http://cloudservices-site-master.test.noahdev.nl/api"
	apiUrl = "https://cloud.noahdev.nl/api"
)

type Api struct{}

func NewApi() *Api {
	return &Api{}
}

func (api *Api) sendRequest(method string, endpoint string, payload []byte) (*http.Response, string, error) {
	log.Println("SENDING TO API")

	// Create a new request with the appropriate method, URL, and payload
	req, err := http.NewRequest(method, fmt.Sprint(apiUrl, endpoint), bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, "", err
	}

	// Set the Content-Type header for JSON
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Make the POST request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return nil, "", err
	}
	defer resp.Body.Close()

	// Read the response body
	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(resp.Body)
	respBody := buf.String()
	log.Println("STATUSCODE:", resp.StatusCode)

	return resp, respBody, nil
}

func (api *Api) UpdateJobStatus(id string, status string) (string, error) {
	type UpdateJobStatusPayload struct {
		ProjectId string `json:"project_id"`
		Status    string `json:"status"`
	}
	payload := UpdateJobStatusPayload{
		ProjectId: id,
		Status:    status,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	resp, respBody, err := api.sendRequest("POST", "/projects/update", jsonPayload)
	if resp.StatusCode != http.StatusAccepted {
		return "", errors.New(respBody)
	}

	return respBody, nil
}

func (api *Api) UpdateFilesPath(id string, filepath string) (string, error) {
	type UpdateJobFilesPathPayload struct {
		ProjectId string `json:"project_id"`
		FilesPath string `json:"files_path"`
	}
	payload := UpdateJobFilesPathPayload{
		ProjectId: id,
		FilesPath: filepath,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	resp, respBody, err := api.sendRequest("POST", "/projects/update", jsonPayload)
	if resp.StatusCode != http.StatusAccepted {
		return "", errors.New(respBody)
	}

	return respBody, nil
}

func (api *Api) UpdateScreenshot(id string, image []byte) error {
	type UpdateScreenshotPayload struct {
		ProjectId    string `json:"project_id"`
		PreviewImage string `json:"preview_image"`
	}

	encodedScreenshot := base64.StdEncoding.EncodeToString(image)

	payload := UpdateScreenshotPayload{
		ProjectId:    id,
		PreviewImage: encodedScreenshot,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, respBody, err := api.sendRequest("POST", "projects/update", jsonPayload)
	if resp.StatusCode != http.StatusAccepted {
		return errors.New(respBody)
	}

	return nil
}
