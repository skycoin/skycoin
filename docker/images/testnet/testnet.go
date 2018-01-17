// Package testnet provides a Docker test net for SkyCoin
package main

import (
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"text/template"
)

// DockerParameters Dockerfile parameters
type DockerParameters struct {
	SkyCoinParameters string
	BuildContext      string
	GitCommit         string
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

// CreateDockerFile makes the Dockerfiles needed to build the images
// for the testnet
func CreateDockerFile(parameters DockerParameters) {
	dockerfileTemplate := path.Join("templates", "Dockerfile")
	_, err := os.Stat(dockerfileTemplate)
	if err != nil {
		log.Print(err)
		return
	}
	f, err := os.Create("Dockerfile")
	if err != nil {
		log.Print(err)
		return
	}
	buildTemplate, err := template.ParseFiles(dockerfileTemplate)
	err = buildTemplate.Execute(f, parameters)
	if err != nil {
		log.Print(err)
		return
	}
}

// GenerateDockerFile generates a Dockerfile with the passed parameters
func GenerateDockerFiles() {
	commonParameters := []string{
		"-launch-browser false",
	}
	guiParameters := []string{
		"--gui-dir=/usr/local/skycoin/static",
		"--web-interface-addr=0.0.0.0",
	}
	noGuiParameters := []string{
		"-web-interface false",
		"-web-interface false",
	}
	parameters := DockerParameters{
		SkyCoinParameters: strings.Join(append(commonParameters, noGuiParameters...), "\", \""),
		BuildContext:      "run",
		GitCommit:         GetCurrentGitCommit(),
	}
	noparameters := DockerParameters{
		SkyCoinParameters: strings.Join(append(commonParameters, guiParameters...), "\", \""),
		BuildContext:      "run",
		GitCommit:         GetCurrentGitCommit(),
	}
	CreateDockerFile(parameters)
	CreateDockerFile(noparameters)
}
func main() {
	GenerateDockerFiles()
}
