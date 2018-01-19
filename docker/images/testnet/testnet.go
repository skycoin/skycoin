// Package testnet provides a Docker test net for SkyCoin
package main

import (
	"log"
	"os"
	"os/exec"
	"path"
//	"strings"
	"text/template"
)

// DockerFile object to generate dockerfiles.
type DockerFile struct {
	SkyCoinParameters []string
	NodeType          string
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
func (d *DockerFile) CreateDockerFile() {
	dockerfileTemplate := path.Join("templates", "Dockerfile")
	_, err := os.Stat(dockerfileTemplate)
	if err != nil {
		log.Print(err)
		return
	}
	buildTemplate, err := template.ParseFiles(dockerfileTemplate)
	f, err := os.Create("Dockerfile-" + d.NodeType)
	if err != nil {
		log.Print(err)
		return
	}
	err = buildTemplate.Execute(f, d)
	if err != nil {
		log.Print(err)
		return
	}
	f.Close()
}

// ConfigureNodes generates a Dockerfile with the passed parameters
func ConfigureNodes() {
	var nodes []DockerFile
	currentCommit := GetCurrentGitCommit()
	nodes = append(nodes, DockerFile{
		NodeType:          "gui",
		SkyCoinParameters: []string{
			"--gui-dir=/usr/local/skycoin/static",
			"--web-interface-addr=0.0.0.0",
		},
		GitCommit:         currentCommit,
	})
	nodes = append(nodes, DockerFile{
		NodeType:          "nogui",
		SkyCoinParameters: []string{
			"-web-interface false",
		},
		GitCommit:         currentCommit,
	})
	for _, d := range nodes {
		d.CreateDockerFile()
	}
}
func main() {
	ConfigureNodes()
}
