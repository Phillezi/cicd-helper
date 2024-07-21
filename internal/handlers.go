package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type ForwardRequestBody struct {
	Command string `json:"command"`
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

		client := &http.Client{}
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
