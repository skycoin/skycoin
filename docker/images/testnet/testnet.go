// Package testnet provides a Docker test net for SkyCoin
package main

import (
	"bufio"
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"
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

// BuildImage builds the node docker image
func (d *DockerFile) BuildImage() {
	cmdName := "docker"
	imageTag := "skycoin-" + d.NodeType + ":" + d.GitCommit
	dockerFile := "Dockerfile-" + d.NodeType
	cmdArgs := []string{
		"build",
		"-f", dockerFile,
		"-t", imageTag,
		"../../../",
	}
	cmd := exec.Command(cmdName, cmdArgs...)
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating StdoutPipe for Cmd", err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			fmt.Printf("docker build out | %s\n", scanner.Text())
		}
	}()

	fmt.Println("Building " + imageTag)

	err = cmd.Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error starting Cmd", err)
		os.Exit(1)
	}

	err = cmd.Wait()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error waiting for Cmd", err)
		os.Exit(1)
	}
	os.Remove(dockerFile)
}

// ConfigureNodes generates a Dockerfile with the passed parameters
func ConfigureNodes() {
	commonParameters := []string{
		"--launch-browser=false",
		"--gui-dir=/usr/local/skycoin/static",
	}
	currentCommit := GetCurrentGitCommit()
	nodes := [2]DockerFile{
		DockerFile{
			NodeType: "gui",
			SkyCoinParameters: []string{
				"--web-interface-addr=0.0.0.0",
			},
			GitCommit: currentCommit,
		},
		DockerFile{
			NodeType: "nogui",
			SkyCoinParameters: []string{
				"--web-interface=false",
			},
			GitCommit: currentCommit,
		},
	}
	for _, d := range nodes {
		d.SkyCoinParameters = append(d.SkyCoinParameters, commonParameters...)
		d.CreateDockerFile()
		d.BuildImage()
	}
}

// GenerateDockerCompose generates the compose YAML file.
func GenerateDockerCompose(nodesNum int, commit string) {
	network := "skycoin-" + commit
	compose := []yaml.MapItem{
		yaml.MapItem{"version", 3},
		yaml.MapItem{"services",
			yaml.MapSlice{},
		},
		yaml.MapItem{"networks",
			[]yaml.MapItem{
				yaml.MapItem{network,
					[]yaml.MapItem{
						yaml.MapItem{"driver", "bridge"},
						yaml.MapItem{"ipam",
							[]yaml.MapItem{
								yaml.MapItem{"driver", "default"},
								yaml.MapItem{"config",
									[]yaml.MapItem{
										yaml.MapItem{"subnet", "172.16.200.0/24"},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	for i := 1; i <= nodesNum; i++ {
		num := strconv.Itoa(i)
		service := yaml.MapItem{
			"skycoin-" + num,
			[]yaml.MapItem{
				yaml.MapItem{"image", "skycoin-nogui:" + commit},
				yaml.MapItem{"networks",
					[]yaml.MapItem{
						yaml.MapItem{network,
							[]yaml.MapItem{
								yaml.MapItem{"ipv4_address", "172.16.200." + num},
							},
						},
					},
				},
			},
		}
		compose[1].Value = append(compose[1].Value.(yaml.MapSlice), service)
	}
	text, err := yaml.Marshal(compose)
	if err != nil {
		log.Fatal(err)
		return
	}
	_, err = os.Stdout.Write(text)
}

func main() {
	//	ConfigureNodes()
	GenerateDockerCompose(10, "0f4e23")
}
