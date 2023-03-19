package main

import (
  "log"
  "net/http"
)

func main() {
	// Define the route handler for GET /jobs
	http.HandleFunc("/jobs", getJobs)

	// Start the server on port 8080
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// getJobs handles GET requests to /jobs
func getJobs(writer http.ResponseWriter, r *http.Request) {
	//TODO: Implement function to retrieve jobs
	// For now, just return an empty list
	jobs := []string{}

	// Encode the jobs list as JSON and write to the response
	json.NewEncoder(writer).Encode(jobs)
}
