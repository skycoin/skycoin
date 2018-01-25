// Package testnet provides a Docker test net for SkyCoin
package main

import (
	//"bufio"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"text/template"
)

// SkyCoinTestNetwork encapsulates the test data and functionality
type SkyCoinTestNetwork struct {
	Compose  dockerCompose
	Services []dockerService
	Peers    []string
}

type dockerService struct {
	SkyCoinParameters []string
	ImageName         string
	ImageTag          string
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
func (d *dockerService) CreateDockerFile(tempDir string) {
	dockerfileTemplate := path.Join("templates", "Dockerfile")
	_, err := os.Stat(dockerfileTemplate)
	if err != nil {
		log.Print(err)
		return
	}
	buildTemplate, err := template.ParseFiles(dockerfileTemplate)
	f, err := os.Create(path.Join(tempDir, "Dockerfile-"+d.ImageName))
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

// NewSkyCoinTestNetwork is the SkyCoinTestNetwork factory function
func NewSkyCoinTestNetwork(nodesNum int) SkyCoinTestNetwork {
	t := SkyCoinTestNetwork{}
	ipHostNum := 1
	networkAddr := "172.16.200."
	commonParameters := []string{
		"--launch-browser=false",
		"--gui-dir=/usr/local/skycoin/static",
	}
	currentCommit := GetCurrentGitCommit()
	networkName := "skycoin-" + currentCommit
	t.Services = []dockerService{
		dockerService{
			ImageName: "skycoin-gui",
			SkyCoinParameters: []string{
				"--web-interface-addr=0.0.0.0",
			},
			ImageTag: currentCommit,
			NodesNum: 1,
		},
		dockerService{
			ImageName: "skycoin-nogui",
			SkyCoinParameters: []string{
				"--web-interface=false",
			},
			ImageTag: currentCommit,
			NodesNum: nodesNum,
		},
	}
	t.Compose = dockerCompose{
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
	for _, s := range t.Services {
		for i := 1; i <= s.NodesNum; i++ {
			num := strconv.Itoa(ipHostNum)
			ipAddress := networkAddr + num
			serviceName := "skycoin-" + num
			t.Compose.Services[serviceName] = service{
				Image: s.ImageName + ":" + s.ImageTag,
				Networks: map[string]serviceNetwork{
					string(networkName): serviceNetwork{
						IPv4Address: ipAddress,
					}},
			}
			t.Peers = append(t.Peers, ipAddress+":6000")
			ipHostNum++
		}
		s.SkyCoinParameters = append(s.SkyCoinParameters, commonParameters...)
	}
	return t
}
func (t *SkyCoinTestNetwork) createComposeFile(tempDir string) {
	text, err := yaml.Marshal(t.Compose)
	if err != nil {
		log.Fatal(err)
	}
	f, err := os.Create(path.Join(tempDir, "docker-compose.yml"))
	if err != nil {
		log.Fatal(err)
	}
	_, err = f.Write(text)
	if err != nil {
		log.Fatal(err)
	}
	f.Close()
	if err != nil {
		log.Fatal(err)
	}

}
func (t *SkyCoinTestNetwork) prepareTestEnv() string {
	tempDir, err := ioutil.TempDir("", "skycointest")
	if err != nil {
		log.Fatal(err)
	}
	for _, s := range t.Services {
		s.CreateDockerFile(tempDir)
	}
	peersText := []byte(strings.Join(t.Peers, "\n"))
	for k, _ := range t.Compose.Services {
		err := os.Mkdir(path.Join(tempDir, k), os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
		f, err := os.Create(path.Join(tempDir, k, "peers.txt"))
		if err != nil {
			log.Fatal(err)
		}
		_, err = f.Write(peersText)
		if err != nil {
			log.Fatal(err)
		}
		f.Close()
		if err != nil {
			log.Fatal(err)
		}
	}
	return tempDir
}

func main() {
	testNet := NewSkyCoinTestNetwork(10)
	tempDir := testNet.prepareTestEnv()
	testNet.createComposeFile(tempDir)
	/*
		err := os.RemoveAll(tempDir)
		if err != nil {
			log.Fatal(err)
		}
	*/
	fmt.Println(tempDir)
}
