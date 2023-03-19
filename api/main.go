package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os/exec"
)

struct Job {
	JobID string `json:"job_id"`
	Status string `json:"status"`
	Command string `json:"command"`
	ExitCode int `json:"exit_code"`
	CreatedAt string `json:"created_at"`
	StartedAt string `json:"started_at"`
	FinishedAt string `json:"finished_at"`
	Duration int `json:"duration"`
	Nodes int `json:"nodes"`
	Memory int `json:"memory"`
	CPU int `json:"cpu"`
	Disk int `json:"disk"`
}

func main() {
	// Define the route handler for GET /jobs
	http.HandleFunc("/jobs", getJobs)

	// Start the server on port 8080
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// getJobs handles GET requests to /jobs
func getJobs(writer http.ResponseWriter, r *http.Request) {
	//TODO: Implement function to retrieve SLURM jobs
	output, err := exec.Command("squeue -u").Output()
	if err!= nil {
        writer.WriteHeader(http.StatusInternalServerError)
        writer.Write([]byte(err.Error()))
        return
    }

	// Extract a slice of jobs from the output of `squeue -u`
	var jobs []Job
	
	// Encode the jobs list as JSON and write to the response
	json.NewEncoder(writer).Encode(jobs)
}
