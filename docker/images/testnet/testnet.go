// Package testnet provides a Docker test net for SkyCoin
package main

import (
	//"bufio"
	//"fmt"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"
	"text/template"
)

type dockerService struct {
	SkyCoinParameters []string
	Image             string
	Tag               string
	NodesNum          int
}

type serviceNetwork struct {
	IPv4Address string `yaml:"ipv4_address"`
}
type service struct {
	Image    string
	Networks map[string]serviceNetwork
}
type networkIpamConfig struct {
	Subnet string
}
type networkIpam struct {
	Driver string
	Config networkIpamConfig
}
type network struct {
	Driver string
	Ipam   networkIpam
}
type dockerCompose struct {
	Version  int
	Services map[string]service
	Networks map[string]network
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
func (d *dockerService) CreateDockerFile() {
	dockerfileTemplate := path.Join("templates", "Dockerfile")
	_, err := os.Stat(dockerfileTemplate)
	if err != nil {
		log.Print(err)
		return
	}
	buildTemplate, err := template.ParseFiles(dockerfileTemplate)
	f, err := os.Create("Dockerfile-" + d.Image)
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

/*
// BuildImage builds the node docker image
func (d *DockerService) BuildImage() {
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
*/

// Configure generates a Dockerfile with the passed parameters
func Configure(nodesNum int) {
	ipHostNum := 1
	networkAddr := "172.16.200."
	commonParameters := []string{
		"--launch-browser=false",
		"--gui-dir=/usr/local/skycoin/static",
	}
	currentCommit := GetCurrentGitCommit()
	networkName := "skycoin-" + currentCommit
	nodes := [2]dockerService{
		dockerService{
			Image: "skycoin-gui",
			SkyCoinParameters: []string{
				"--web-interface-addr=0.0.0.0",
			},
			Tag:      currentCommit,
			NodesNum: 1,
		},
		dockerService{
			Image: "skycoin-nogui",
			SkyCoinParameters: []string{
				"--web-interface=false",
			},
			Tag:      currentCommit,
			NodesNum: nodesNum,
		},
	}
	compose := dockerCompose{
		Version:  3,
		Services: make(map[string]service),
		Networks: map[string]network{
			string(networkName): network{
				Driver: "bridge",
				Ipam: networkIpam{
					Driver: "default",
					Config: networkIpamConfig{
						Subnet: networkAddr + "0/24",
					},
				},
			},
		},
	}
	for _, d := range nodes {
		for i := 1; i <= d.NodesNum; i++ {
			num := strconv.Itoa(ipHostNum)
			compose.Services["skycoin-"+num] = service{
				Image: d.Image + ":" + d.Tag,
				Networks: map[string]serviceNetwork{
					string(networkName): serviceNetwork{
						IPv4Address: networkAddr + num,
					}},
			}
			ipHostNum++
		}
		d.SkyCoinParameters = append(d.SkyCoinParameters, commonParameters...)
		d.CreateDockerFile()
		//d.BuildImage()
	}
	text, err := yaml.Marshal(compose)
	if err != nil {
		log.Fatal(err)
	}
	_, err = os.Stdout.Write(text)
	if err != nil {
		log.Fatal(err)
		return
	}
}

// GenerateDockerCompose generates the compose YAML file.
func GenerateDockerCompose(nodesNum int, commit string) {
}

func main() {
	//	ConfigureNodes()
	GenerateDockerCompose(10, "0f4e23")
	Configure(15)
}
