package internal

import (
	"log"
	"net/http"
	"os"
)

func Run() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	apiURL := os.Getenv("API_URL")
	if apiURL == "" {
		apiURL = "https://api.cloud.cbh.kth.se"
	}

	http.HandleFunc("/forward", ForwardRequest(apiURL))
	http.HandleFunc("/harbor/restart", HarborAutoRestart(apiURL))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		githubURL := "https://github.com/Phillezi/cicd-helper"
		http.Redirect(w, r, githubURL, http.StatusSeeOther)
	})

	log.Printf("Server is listening on port %s...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
