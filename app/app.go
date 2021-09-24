package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	actions "github.com/sethvargo/go-githubactions"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
)

var defaultRepositoryOwner string = "hashicorp"

type args struct {
	name  string
	value string
}

const (
	repository      = iota
	repositoryOwner = iota
	versionCommand  = iota
	filePath = iota
	metadataJsonFile = iota
)

type Metadata struct {
	Product         string `json:"repository"`
	Org             string `json:"repositoryOwner"`
	Revision        string `json:"sha"`
	Version         string `json:"version"`
	BuildWorkflowId string `json:"buildWorkflowId"`
}

func main() {

	in := []args{
		args{
			name: "repository",
		},
		args{
			name: "repositoryOwner",
		},
		args{
			name: "versionCommand",
		},
		args{
			name: "filePath",
		},
		args{
			name: "metadataJsonFile",
		},
	}
	for i := range in {
		getInputsValue(&in[i].value, in[i].name)
	}
	actions.Infof("ARGS %v\n", in)

	parsedOutputs := creteMetadataJson(in)

	actions.SetOutput("output", parsedOutputs)
}

func creteMetadataJson(in []args) {
	sha := getSha()
	actions.Infof("Working sha %s\n", sha)

	repository := in[repository].value

	org := in[repositoryOwner].value
	if org == "" {
		org = defaultRepositoryOwner
	}

	runId := os.Getenv("GITHUB_RUN_ID")
	if runId == "" {
		actions.Fatalf("GITHUB_RUN_ID is empty")
	}

	version := getVersion(in[versionCommand].value)

	actions.Infof("Creating metadata.json file")
	m := &Metadata{
		Product:         repository,
		Org:             org,
		Revision:        sha,
		BuildWorkflowId: runId,
		Version:         version}
	output, err := json.MarshalIndent(m, "", "\t\t")
	if err != nil {
		fmt.Println("Error", err)
		return
	}
	err = ioutil.WriteFile(path.Join(in[filePath].value), output))
	if err != nil {
		actions.Fatalf("Failed writing data into metadata.json file. Error: %v\n", err)
		return
	}
	actions.Infof("Successfully created metadata.json file")

}

func getSha() string {
	// git rev-parse HEAD
	sha := execCommand("git", "rev-parse", "HEAD")
	if len(sha) == 0 {
		actions.Fatalf("Unable to determine git sha for this commit")
	} else {
		sha = sha[0 : len(sha)-1]
	}

	return sha
}

func getVersion(command string) string {
	version := strings.TrimSuffix(execCommand(string(command[0]), command[1:]), "\n")
	actions.Infof("Running %v version\n", version)
	if version == "" {
		//sha = os.Getenv("GITHUB_SHA")
		actions.Fatalf("Failed to setup version using `make version` command")
	}
	return version
}

func execCommand(args ...string) string {
	name := args[0]
	stderr := new(bytes.Buffer)
	stdout := new(bytes.Buffer)

	cmd := exec.Command(name, args[1:]...)
	cmd.Stderr = stderr
	cmd.Stdout = stdout
	err := cmd.Run()
	actions.Infof("Running %v command: %v\nstdout: %v\nstderr: %v\n", name, cmd,
		strings.TrimSpace(string(stdout.Bytes())), strings.TrimSpace(string(stderr.Bytes())))

	if err != nil {
		actions.Fatalf("Failed to run %v command %v: %v", name, cmd, err)
	}

	return string(stdout.Bytes())
}

func getInputsValue(val *string, key string) {
	*val = actions.GetInput(key)
}
