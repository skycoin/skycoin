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
	Image string
	Ports []string `yaml:"ports,omitempty"`
	Build serviceBuild
}
type serviceBuild struct {
	Context    string `yaml:"context,omitempty"`
	Dockerfile string
}
type dockerCompose struct {
	Version  string
	Services map[string]dockerService
}

// BuildImage builds the node docker image
func runTest(tempDir string, scale int) {
	cmdName := "docker-compose"
	composeFile := path.Join(tempDir, "docker-compose.yml")
	cmdArgs := []string{"-f", composeFile, "up", "--scale", "skycoin-slave=" + strconv.Itoa(scale)}
	cmd := exec.Command(cmdName, cmdArgs...)
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating StdoutPipe for Cmd", err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			fmt.Printf("%s\n", scanner.Text())
		}
	}()

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
	fmt.Println("Cleaning up...")
	cmdArgs = []string{"-f", composeFile, "down"}
	err = os.RemoveAll(tempDir)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Done.")
}

func processComposeFile(composeFile *dockerCompose, tempDir string) {
	for k, v := range composeFile.Services {
		v.Image = k
		v.Build.Dockerfile = tempDir + "/Dockerfile-" + k
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
	for k, service := range composeFile.Services {
		dockerfileExplorerSrc := path.Join("templates", "Dockerfile-"+k)
		dockerfileExplorerDst := service.Build.Dockerfile
		df, err := os.Open(dockerfileExplorerSrc)
		_, err = file.CopyFile(dockerfileExplorerDst, df)
		if err != nil {
			log.Fatal(err)
		}
	}

}
func main() {
	_, callerFile, _, _ := runtime.Caller(0)
	projectPath, _ := filepath.Abs(filepath.Join(filepath.Dir(callerFile), "../../../"))
	log.Print("Source code base dir at ", projectPath)
	//nodesPtr := flag.Int("-nodes", 5, "Number of nodes to launch.")
	//buildContextPtr := flag.String("-buildcontext", projectPath, "Docker build context (source code root).")
	flag.Parse()
	//buildContext, err := filepath.Abs(*buildContextPtr)
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
					Context: projectPath,
				},
				Ports: []string{"16420:16420"},
			},
			"skycoin-peer": dockerService{
				Build: serviceBuild{
					Context: projectPath,
				},
			},
			"skycoin-explorer": dockerService{
				Build: serviceBuild{
					Context: ".",
				},
				Ports: []string{"8001:8001"},
			},
		},
	}
	processComposeFile(&composeFile, tempDir)
	createComposeFile(composeFile, tempDir)
	prepareTestEnv(composeFile)
	//testNet := NewSkyCoinTestNetwork(*nodesPtr, buildContext, tempDir)
	//runTest(tempDir)
}
