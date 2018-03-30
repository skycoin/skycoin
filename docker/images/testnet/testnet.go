// Package testnet provides a Docker test net for SkyCoin
package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/skycoin/skycoin/src/util/file"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
)

type dockerService struct {
	Image   string
	Ports   []string     `yaml:"ports,omitempty"`
	Volumes []string     `yaml:"volumes,omitempty"`
	Command string       `yaml:"command,omitempty"`
	Build   serviceBuild `yaml:"build,omitempty"`
}
type serviceBuild struct {
	Context    string `yaml:"context,omitempty"`
	Dockerfile string `yaml:"dockerfile,omitempty"`
	Template   string `yaml:"-"`
}
type dockerCompose struct {
	Version  string
	Services map[string]dockerService
}

// BuildImage builds the node docker image
func runTestNet(tempDir string, scale int) {
	cmdName := "docker-compose"
	cmdArgs := []string{"up", "-d", "--scale", "skycoin-peer=" + strconv.Itoa(scale)}
	cmd := exec.Command(cmdName, cmdArgs...)
	cmd.Dir = tempDir
	err := cmd.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error starting Testnet", err)
		cleanUp(tempDir)
		os.Exit(1)
	}
}

func stopTestNet(tempDir string) {
	cmdName := "docker-compose"
	cmdArgs := []string{"down"}
	cmd := exec.Command(cmdName, cmdArgs...)
	cmd.Dir = tempDir
	err := cmd.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error stopping Testnet", err)
		cleanUp(tempDir)
		os.Exit(1)
	}
	cleanUp(tempDir)
	fmt.Println("Done.")
}

func cleanUp(tempDir string) {
	err := os.RemoveAll(tempDir)
	if err != nil {
		log.Fatal(err)
	}
}

func processComposeFile(composeFile *dockerCompose, tempDir string) {
	for k, v := range composeFile.Services {
		if v.Image == "" {
			v.Image = k
		}
		templateFile := path.Join("templates", "Dockerfile-"+k)
		_, err := os.Stat(templateFile)
		if err == nil {
			v.Build.Template = templateFile
			v.Build.Dockerfile = path.Join(tempDir, "Dockerfile-"+k)
		}
		composeFile.Services[k] = v
	}
}

func createComposeFile(composeFile dockerCompose, tempDir string) {
	text, err := yaml.Marshal(composeFile)
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

func prepareTestEnv(composeFile dockerCompose) {
	// Copies Dockerfiles
	for _, service := range composeFile.Services {
		if service.Build.Template != "" {
			df, err := os.Open(service.Build.Template)
			_, err = file.CopyFile(service.Build.Dockerfile, df)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

}
func main() {
	_, callerFile, _, _ := runtime.Caller(0)
	projectPath, _ := filepath.Abs(filepath.Join(filepath.Dir(callerFile), "../../../"))
	log.Print("Source code base dir at ", projectPath)
	nodesPtr := flag.Int("nodes", 5, "Number of nodes to launch.")
	buildContextPtr := flag.String("buildcontext", projectPath, "Docker build context (source code root).")
	flag.Parse()
	buildContext, err := filepath.Abs(*buildContextPtr)
	tempDir, err := ioutil.TempDir("", "skycointest")
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Fatal(err)
	}
	composeFile := dockerCompose{
		Version: "3",
		Services: map[string]dockerService{
			"skycoin-master": dockerService{
				Build: serviceBuild{
					Context: buildContext,
				},
				Ports: []string{"16420:16420"},
			},
			"skycoin-peer": dockerService{
				Build: serviceBuild{
					Context: buildContext,
				},
			},
			"skycoin-explorer": dockerService{
				Build: serviceBuild{
					Context: ".",
				},
				Ports: []string{"8001:8001"},
			},
			"gliderlabs_logspout": dockerService{
				Image:   "gliderlabs/logspout",
				Volumes: []string{"/var/run/docker.sock:/var/run/docker.sock"},
				Command: "raw+tcp://oklog:7651?filter.name=*skycoin*",
			},
			"oklog": dockerService{
				Build: serviceBuild{
					Context: ".",
				},
				Ports: []string{"7651:7651"},
			},
		},
	}
	processComposeFile(&composeFile, tempDir)
	createComposeFile(composeFile, tempDir)
	prepareTestEnv(composeFile)
	runTestNet(tempDir, *nodesPtr)
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Press enter to stop TestNet: ")
	_, _ = reader.ReadString('\n')
	stopTestNet(tempDir)
}
