// Package testnet provides a Docker test net for SkyCoin
package main

import (
	"html/template"
	"log"
	"os"
	"os/exec"
	"path"
)

// DockerParameters Dockerfile parameters
type DockerParameters struct {
	Execution    string
	BuildContext string
	GitCommit    string
}

// GetCurrentGitCommit returns the current git commit SHA
func GetCurrentGitCommit() string {
	var (
		cmdOut []byte
		err    error
	)
	cmdName := "git"
	cmdArgs := []string{"rev-parse", "--verify", "HEAD"}
	if cmdOut, err = exec.Command(cmdName, cmdArgs...).Output(); err != nil {
		log.Print("There was an error running git rev-parse command: ", err)
		os.Exit(1)
	}
	sha := string(cmdOut)
	firstSix := sha[:6]
	return firstSix
}

// CreateDockerfiles makes the Dockerfiles needed to build the images
// for the testnet
func CreateDockerfiles() {
	gitCommit := GetCurrentGitCommit()
	parameters := DockerParameters{
		Execution:    "run",
		BuildContext: "run",
		GitCommit:    gitCommit,
	}
	dockerfileTemplate := path.Join("templates", "Dockerfile")
	_, err := os.Stat(dockerfileTemplate)
	if err != nil {
		log.Print(err)
		return
	}
	buildTemplate, err := template.ParseFiles(dockerfileTemplate)
	err = buildTemplate.Execute(os.Stdout, parameters)
	if err != nil {
		log.Print(err)
		return
	}
}
func main() {
	CreateDockerfiles()
}
