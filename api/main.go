package main

import (
    "encoding/json"
    "log"
    "net/http"
    "strings"
    "time"
    "strconv"
    "os/exec"
)

type Job struct {
    JobID      string `json:"job_id"`
    Status     string `json:"status"`
    Command    string `json:"command"`
    ExitCode   int    `json:"exit_code"`
    CreatedAt  string `json:"created_at"`
    StartedAt  string `json:"started_at"`
    FinishedAt string `json:"finished_at"`
    Duration   int    `json:"duration"`
    Nodes      int    `json:"nodes"`
    Memory     int    `json:"memory"`
    CPU        int    `json:"cpu"`
    Disk       int    `json:"disk"`
}

func main() {
    // Define the route handler for GET /jobs
    http.HandleFunc("/jobs", getJobs)

    // Start the server on port 8080
    log.Fatal(http.ListenAndServe(":8080", nil))
}

// getJobs handles GET requests to /jobs
func getJobs(writer http.ResponseWriter, r *http.Request) {
    // Execute `squeue -u` command and capture output
    output, err := exec.Command("squeue", "-u", "<USERNAME>").Output()
    if err != nil {
        writer.WriteHeader(http.StatusInternalServerError)
        writer.Write([]byte(err.Error()))
        return
    }

    // Extract jobs from output and create Job objects
    var jobs []Job
    lines := strings.Split(string(output), "\n")[1:]
    for _, line := range lines {
        if line == "" {
            continue
        }
        fields := strings.Fields(line)
        jobID := fields[0]
        status := fields[4]
        command := strings.Join(fields[7:], " ")
        exitCode, err := strconv.Atoi(fields[5])
        if err != nil {
            exitCode = -1
        }
        createdAt, err := time.Parse("2006-01-02T15:04:05", fields[1]+"T"+fields[2])
        if err != nil {
            createdAt = time.Time{}
        }
        nodes, err := strconv.Atoi(fields[6])
        if err != nil {
            nodes = -1
        }
        jobs = append(jobs, Job{
            JobID:     jobID,
            Status:    status,
            Command:   command,
            ExitCode:  exitCode,
            CreatedAt: createdAt.String(),
            Nodes:     nodes,
        })
    }

    // Encode the jobs list as JSON and write to the response
    writer.Header().Set("Content-Type", "application/json")
    err = json.NewEncoder(writer).Encode(jobs)
    if err != nil {
        writer.WriteHeader(http.StatusInternalServerError)
        writer.Write([]byte(err.Error()))
        return
    }
}
