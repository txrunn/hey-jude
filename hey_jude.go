package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func uploadFile(client *ssh.Client, localPath, remotePath string) error {
	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		return err
	}
	defer sftpClient.Close()

	srcFile, err := os.Open(localPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := sftpClient.Create(remotePath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

func readInput(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func createSSHClient(username, password, host string) (*ssh.Client, error) {
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", host+":22", config)
	return client, err
}

func executeSSHCommand(client *ssh.Client, command string) (string, error) {
	session, err := client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	outputBytes, err := session.CombinedOutput(command)
	if err != nil {
		return "", err
	}

	return string(outputBytes), nil
}

func main() {
	// Read user input for SSH credentials and other parameters
	username := readInput("SSH username: ")
	password := readInput("SSH password: ")
	host := readInput("HPC hostname: ")

	// Create an SSH client
	client, err := createSSHClient(username, password, host)
	if err != nil {
		fmt.Printf("Failed to create SSH client: %s\n", err)
		return
	}
	defer client.Close()

	// Read or generate the R script and SLURM script based on user input
	rScript := "generated_r_script.R"
	slurmScript := "generated_slurm_script.slurm"

	// Upload the R and SLURM scripts to the HPC system
	err = uploadFile(client, rScript, rScript)
	if err != nil {
		fmt.Printf("Failed to upload R script: %s\n", err)
		return
	}

	err = uploadFile(client, slurmScript, slurmScript)
	if err != nil {
		fmt.Printf("Failed to upload SLURM script: %s\n", err)
		return
	}

	// Submit the job using `sbatch`
	command := fmt.Sprintf("sbatch %s", slurmScript)
	output, err := executeSSHCommand(client, command)
	if err != nil {
		fmt.Printf("Failed to submit job: %s\n", err)
		return
	}
	fmt.Printf("Job submitted successfully:\n%s\n", output)
}
