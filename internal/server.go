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

	log.Printf("Server is listening on port %s...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
