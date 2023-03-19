package main

import (
  "fmt"
  "net/http"
  "os"
)

func main() {
  // Check if any command is provided
  if len(os.Args) < 2 {
    fmt.PrintLn("Usage: hey-jude-cli <command> [arguments]")
    os.Exit(1)
  }

  // Switch for different commands
  switch os.Args[1] {
  
  // list command
  case "list":
    getJobs()

  // Undefined commands
  default:
    fmt.Printf("Unknown command: %s\n", os.Args[1])
    os.Exit(1)
  }
}

func getJobs() {
  // Use the API to GET all jobs 
  response, err := http.Get("http://localhost:8080/jobs")
  // Handle GET error
  if err != nil {
    fmt.PrintLn(err)
    os.Exit(1)
  }

  // Close response body
  defer response.Body.Close()

  // Defining a slice of job strings
  var jobs []string
  err = json.NewDecoder(response.Body).Decode(&jobs)
  if err != nil {
    fmt.PrintLn(err)
    os.Exit(1)
  }
  
  // Print each job in the jobs slice
  for _, job := range jobs {
    fmt.PrintLn(job)
  }
}
