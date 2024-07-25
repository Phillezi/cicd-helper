package internal

import (
	"bytes"
	"cicd-helper/model/harbor"
	"encoding/json"
	"errors"
	"fmt"
	"go-deploy/dto/v2/body"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type ForwardRequestBody struct {
	Command string `json:"command"`
}

func getAffectedDepls(apiURL string, token string, image string) ([]body.DeploymentRead, error) {
	var deployments []body.DeploymentRead
	url := fmt.Sprintf("%s/deploy/v2/deployments?all=true&shared=true", apiURL)

	log.Infoln("Sending GET request to:", url)

	req, err := http.NewRequest("GET", url, http.NoBody)
	if err != nil {
		return nil, errors.New("error creating forward request")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Key", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Errorln("Error GET-ting deployments", err)
		return nil, err
	}
	defer resp.Body.Close()
	rbody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if len(rbody) == 0 {
		return nil, errors.New("empty response")
	}
	err = json.Unmarshal(rbody, &deployments)
	if err != nil {
		return nil, err
	}
	if deployments == nil {
		return nil, errors.New("no deployments in response")
	}

	var affected []body.DeploymentRead
	for _, deployment := range deployments {
		if *deployment.Image == image {
			affected = append(affected, deployment)
		}
	}
	return affected, nil
}

func HarborAutoRestart(apiURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Infoln("Request recieved for /harbor/restart")

		token := r.Header.Get("Authorization")
		if token == "" {
			log.Errorln("Authorization header is required", http.StatusUnauthorized)
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		log.Infoln("Request is authorized")

		defer r.Body.Close()
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println("Error reading request body:", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		var requestBody harbor.RequestBody
		err = json.Unmarshal(body, &requestBody)
		if err != nil {
			log.Println("Error unmarshaling JSON:", err)
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if len(requestBody.EventData.Resources) <= 0 {
			log.Errorln("No resources found in the event data.")
			http.Error(w, "No resources found in the event data", http.StatusBadRequest)
			return
		}
		resourceURL := requestBody.EventData.Resources[0].ResourceURL

		deployments, err := getAffectedDepls(apiURL, token, resourceURL)
		if err != nil {
			log.Errorln("Could not get affected deployments: ", err)
			http.Error(w, "Could not get affected deployments", http.StatusInternalServerError)
			return
		}

		for _, deployment := range deployments {
			// TODO: find better way
			if deployment.PingResult == nil {
				log.Infoln("Deployment with ID: ", deployment.ID, " has status: ", deployment.Status, " (skipping restart)")
				continue
			}
			code, message, err := sendRestartRequest(apiURL, token, deployment.ID)
			if err != nil {
				msg := "Could not restart deployment with ID: " + deployment.ID
				log.Errorln(msg, " error: ", err)
				http.Error(w, msg, http.StatusInternalServerError)
				return
			}
			if code != 204 {
				msg := "Could not restart deployment with ID: " + deployment.ID + " response: " + message
				log.Errorln(msg, " error: ", err)
				http.Error(w, msg, http.StatusInternalServerError)
				return
			}
		}
		w.WriteHeader(204)
	}
}

func sendRestartRequest(apiURL, token, deploymentID string) (int, string, error) {
	url := fmt.Sprintf("%s/deploy/v2/deployments/%s/command", apiURL, deploymentID)

	forwardBody := ForwardRequestBody{
		Command: "restart",
	}

	jsonData, err := json.Marshal(forwardBody)
	if err != nil {
		return http.StatusInternalServerError, "", fmt.Errorf("error marshalling JSON: %v", err)
	}

	log.Println("Sending restart request to:", url)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return http.StatusInternalServerError, "", fmt.Errorf("error creating forward request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Key", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return http.StatusInternalServerError, "", fmt.Errorf("error forwarding request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, "", fmt.Errorf("error reading response body: %v", err)
	}

	return resp.StatusCode, string(body), nil
}

func ForwardRequest(apiURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Infoln("Request recieved for /forward")

		deploymentID := r.URL.Query().Get("deploymentid")
		if deploymentID == "" {
			log.Errorln("deploymentid query parameter is required", http.StatusBadRequest)
			http.Error(w, "deploymentid query parameter is required", http.StatusBadRequest)
			return
		}

		log.Infoln("Request regards deploymentid:", deploymentID)

		token := r.Header.Get("Authorization")
		if token == "" {
			log.Errorln("Authorization header is required", http.StatusUnauthorized)
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		log.Infoln("Request is authorized")

		url := fmt.Sprintf("%s/deploy/v2/deployments/%s/command", apiURL, deploymentID)

		forwardBody := ForwardRequestBody{
			Command: "restart",
		}

		jsonData, err := json.Marshal(forwardBody)
		if err != nil {
			log.Errorln("Error marshalling JSON", http.StatusInternalServerError)
			http.Error(w, "Error marshalling JSON", http.StatusInternalServerError)
			return
		}

		log.Infoln("Sending restart request to:", url)

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		if err != nil {
			log.Errorln("Error creating forward request", http.StatusInternalServerError)
			http.Error(w, "Error creating forward request", http.StatusInternalServerError)
			return
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Api-Key", token)

		tr := &http.Transport{
			//TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}

		client := &http.Client{Transport: tr}
		resp, err := client.Do(req)
		if err != nil {
			log.Errorln("Error forwarding request", err)
			http.Error(w, "Error forwarding request", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	}
}
