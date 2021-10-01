package main

import (
	"bytes"
	"encoding/json"
	actions "github.com/sethvargo/go-githubactions"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
)

var defaultRepositoryOwner string = "hashicorp"
var defaultMetadataFileName string = "metadata.json"

type args struct {
	name  string
	value string
}

const (
	repository       = iota
	repositoryOwner  = iota
	filePath         = iota
	metadataFileName = iota
	version          = iota
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
			name: "filePath",
		},
		args{
			name: "metadataFileName",
		},
		args{
			name: "version",
		},
	}
	for i := range in {
		getInputsValue(&in[i].value, in[i].name)
	}
	generatedFile := createMetadataJson(in)

	if checkFileIsExist(generatedFile) {
		actions.SetOutput("filepath", generatedFile)
		actions.SetEnv("filepath", generatedFile)
		actions.Infof("Successfully created %v file\n", generatedFile)
	} else {
		actions.Fatalf("File %v does not exist", generatedFile)
	}
}

func checkFileIsExist(filepath string) bool {
	fileInfo, err := os.Stat(filepath)

	if os.IsNotExist(err) {
		return false
	}
	// Return false if the fileInfo says the file path is a directory.
	return !fileInfo.IsDir()
}

func createMetadataJson(in []args) string {
	file := in[metadataFileName].value
	if file == "" {
		file = defaultMetadataFileName
	}
	filePath := path.Join(in[filePath].value, file)

	sha := getSha()
	actions.Infof("Working sha %v\n", sha)
	repository := in[repository].value

	org := in[repositoryOwner].value
	if org == "" {
		org = defaultRepositoryOwner
	}

	runId := os.Getenv("GITHUB_RUN_ID")
	if runId == "" {
		actions.Fatalf("GITHUB_RUN_ID is empty")
	}

	version := in[version].value
	if version == "" {
		actions.Fatalf("The version or version command is not provided")
	} else if strings.Contains(version, " ") {
		version = getVersion(version)
	}
	actions.Infof("Working version %v\n", version)

	actions.Infof("Creating metadata file in %v\n", filePath)
	m := &Metadata{
		Product:         repository,
		Org:             org,
		Revision:        sha,
		BuildWorkflowId: runId,
		Version:         version}
	output, err := json.MarshalIndent(m, "", "\t\t")
	if err != nil {
		actions.Fatalf("JSON marshal failure. Error:%v\n", output, err)

	} else {
		err = ioutil.WriteFile(filePath, output, 0644)
		if err != nil {
			actions.Fatalf("Failed writing data into %v file. Error: %v\n", in[metadataFileName].value, err)
		}
	}
	return filePath
}

func getSha() string {
	// git rev-parse HEAD
	sha := execCommand("git", "rev-parse", "HEAD")
	if len(sha) == 0 {
		actions.Fatalf("Failed to determine git sha for this commit")
	} else {
		sha = sha[0 : len(sha)-1]
	}
	return sha
}

func getVersion(command string) string {
	version := execCommand(strings.Fields(command)...)
	if version == "" {
		actions.Fatalf("Failed to setup version using %v command", command)
	}
	return strings.TrimSuffix(version, "\n")
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
