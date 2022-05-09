/*
Copyright © 2022 Yevhen Lebid ealebed@gmail.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/ealebed/spini/types"
	git "github.com/ealebed/spini/utils/github"
	"github.com/google/go-github/v44/github"
	"github.com/google/uuid"
	"github.com/instrumenta/kubeval/kubeval"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	kjson "k8s.io/apimachinery/pkg/runtime/serializer/json"
	"sigs.k8s.io/kustomize/kyaml/kio"
	"sigs.k8s.io/kustomize/kyaml/kio/filters"
)

var (
	removeLines = [][]byte{
		// The Kubernetes standard header field `currentTimestamp` serializes weirdly, so filter it out.
		[]byte("  creationTimestamp: null"),
		[]byte("        resources: {}"),
		[]byte("      resources: {}"),
		[]byte("    resources: {}"),
		[]byte("status: {}"),
		[]byte("spec: {}"),
		[]byte("status:"),
		[]byte("  loadBalancer: {}"),
		[]byte("targetPort: 0"),
	}
)

// ReadJSONLocalToStruct returns struct with configuration values from selected repository
func ReadJSONLocalToStruct() []*types.Configuration {
	configuration := make([]*types.Configuration, 0)

	var filename = "configuration.json"

	file, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Print(err)
	}

	if err := json.Unmarshal([]byte(file), &configuration); err != nil {
		fmt.Printf("Failed to unmarshal JSON file 'configuration.json' due to %s", err)
		os.Exit(1)
	}

	return configuration
}

// GetFileContent loads the local content of a file and return the target name of the file in the target repository and its contents.
func GetFileContent(fileArg string) (targetName string, b []byte, err error) {
	var localFile string
	files := strings.Split(fileArg, ":")
	switch {
	case len(files) < 1:
		return "", nil, errors.New("empty `-files` parameter")
	case len(files) == 1:
		localFile = files[0]
		targetName = files[0]
	default:
		localFile = files[0]
		targetName = files[1]
	}

	b, err = ioutil.ReadFile(localFile)
	return targetName, b, err
}

// WriteFileOnDisk write generated file on disk
func WriteFileOnDisk(out []byte, outputFilePath string) error {
	if err := ioutil.WriteFile(outputFilePath, out, 0644); err != nil {
		fmt.Printf("Failed to safe the generated file. Error: %v", err)
		os.Exit(1)
	}

	return nil
}

// sliceContains checks if a string is present in a slice
func sliceContains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

// fillPipelineConfig fills pipeline configuration for generating spinnaker pipelines
func fillPipelineConfig(stage string, pList []string, pipelineIDs map[string]string) map[string]interface{} {
	var pipeValues = make(map[string]interface{})

	pipeValues["stage"] = stage

	if sliceContains(pList, "beta") || sliceContains(pList, "nightly") {
		if stage == "production" || (strings.Contains(stage, "dev")) {
			pipeValues["dockerTriggerEnabled"] = false
		} else if sliceContains(pList, "beta") && stage == "nightly" {
			pipeValues["dockerTriggerEnabled"] = false
		} else if !sliceContains(pList, "beta") && stage == "nightly" {
			pipeValues["dockerTriggerEnabled"] = true
		} else if !sliceContains(pList, "nightly") && stage == "beta" {
			pipeValues["dockerTriggerEnabled"] = true
		} else {
			pipeValues["dockerTriggerEnabled"] = true
		}
	} else {
		pipeValues["dockerTriggerEnabled"] = true
	}

	if sliceContains(pList, "beta") && (stage == "production" || stage == "nightly") {
		// by default set in all promote pipelines trigger parentPipelineID from 'beta' pipeline
		pipeValues["id"] = pipelineIDs["promote-to-"+stage]
		pipeValues["parentPipelineId"] = pipelineIDs["beta-gke1"]

		// in case we have 'nightly' stage, set in promote-to-production pipeline trigger parentPipelineID from 'nightly' pipeline
		if stage == "production" && pipelineIDs["nightly-gke1"] != "" {
			pipeValues["parentPipelineId"] = pipelineIDs["nightly-gke1"]
		}

		pipeValues["GeneratePromotePipeline"] = true
		pipeValues["pipelineTriggerEnabled"] = true
	} else if sliceContains(pList, "nightly") && stage == "production" {
		// in case we have 'nightly' stage, set in promote-to-production pipeline trigger parentPipelineID from 'nightly' pipeline
		pipeValues["id"] = pipelineIDs["promote-to-"+stage]
		pipeValues["parentPipelineId"] = pipelineIDs["nightly-gke1"]

		pipeValues["GeneratePromotePipeline"] = true
		pipeValues["pipelineTriggerEnabled"] = true
	} else {
		pipeValues["GeneratePromotePipeline"] = false
		pipeValues["pipelineTriggerEnabled"] = false
	}

	return pipeValues
}

// formatManifest validate and format generated kubernetes manifest
func formatManifest(in []byte) (*bytes.Buffer, error) {
	kubevalConfig := &kubeval.Config{
		DefaultNamespace:     "default",
		IgnoreMissingSchemas: true,
		KubernetesVersion:    "master",
		Strict:               true,
	}

	// validate generated kubernetes manifest with kubeval
	_, err := kubeval.Validate(in, kubevalConfig)
	if err != nil {
		return nil, err
	}

	// format the field ordering in generated YAML
	out := &bytes.Buffer{}
	p := kio.Pipeline{
		Inputs:  []kio.Reader{&kio.ByteReader{Reader: bytes.NewReader(in)}},
		Filters: []kio.Filter{filters.FormatFilter{UseSchema: true}},
		Outputs: []kio.Writer{kio.ByteWriter{Writer: out, KeepReaderAnnotations: false}},
	}
	if err := p.Execute(); err != nil {
		return nil, err
	}

	return out, nil
}

// GeneratePipelines returns list generated spinnaker pipeline objects
func GeneratePipelines(app *types.Configuration, organization, githubRepositoryName string) []*types.Pipeline {
	var pipelineNamesList []string

	var generatedPipelineList []*types.Pipeline
	var pipelineIDs = map[string]string{}

	// generate build pipeline
	buildPipeline := types.NewBuildPipeline(app)
	generatedPipelineList = append(generatedPipelineList, buildPipeline)

	// here we collect all pipeline names to list and create pipelineIDs, including IDs for promote to stage pipeline
	for _, profile := range *app.Profiles {
		pipelineNamesList = append(pipelineNamesList, profile.ProfileName)
		if profile.ProfileName != "beta" {
			pipelineIDs["promote-to-"+profile.ProfileName] = uuid.New().String()
		}
		for _, tier := range *profile.Datacenters {
			pipelineIDs[profile.ProfileName+"-"+tier.TierName] = uuid.New().String()
		}
	}

	for _, profile := range *app.Profiles {
		pipeValues := fillPipelineConfig(profile.ProfileName, pipelineNamesList, pipelineIDs)
		pipeValues["organization"] = organization
		pipeValues["githubRepositoryName"] = githubRepositoryName

		// generate promote-to-stage pipelines
		if pipeValues["GeneratePromotePipeline"].(bool) {
			promotePipeline := types.NewPromotePipeline(app, pipeValues)
			generatedPipelineList = append(generatedPipelineList, promotePipeline)
		}

		// generate deploy pipelines
		for _, tier := range *profile.Datacenters {
			pipeValues["cluster"] = tier.TierName
			pipeValues["id"] = pipelineIDs[profile.ProfileName+"-"+tier.TierName]
			pipeValues["parentPipelineId"] = pipelineIDs["promote-to-"+profile.ProfileName]

			if tier.EnvFrom != nil {
				app.EnvFrom = append(app.EnvFrom, tier.EnvFrom...)
			}

			deployPipeline := types.NewDeployPipeline(app, pipeValues)
			generatedPipelineList = append(generatedPipelineList, deployPipeline)
		}
	}

	return generatedPipelineList
}

// GenerateManifests returns generated kubernetes manifest objects
func GenerateManifests(app *types.Configuration, tier *types.Datacenter, stage, organization string) string {
	var buf bytes.Buffer
	list := &metav1.List{
		TypeMeta: metav1.TypeMeta{
			Kind:       "List",
			APIVersion: "v1",
		},
	}

	sa := types.NewServiceAccount(app.Application, stage, app.Namespace, "dockerhubkey")
	list.Items = append(list.Items, runtime.RawExtension{Object: sa})

	if app.Type == "service" {
		s := types.NewService(app.Application, stage, app.Namespace, app.Ports)
		list.Items = append(list.Items, runtime.RawExtension{Object: s})
	}

	if tier.ChaosMonkey == nil {
		tier.ChaosMonkey = app.ChaosMonkey
	}

	d := types.NewDeployment(app, tier, stage, organization)
	list.Items = append(list.Items, runtime.RawExtension{Object: d})

	options := kjson.SerializerOptions{
		Yaml:   true,
		Pretty: true,
		Strict: false,
	}

	e := kjson.NewSerializerWithOptions(kjson.DefaultMetaFactory, nil, nil, options)
	err := e.Encode(list, &buf)
	if err != nil {
		panic(err)
	}

	cleaned := buf.Bytes()

	// We currently use bytes.Replace for ease/speed but we have to repeat extra lines due to varying whitespace.
	// May consider doing regexp but it will probably result in it being slower/more
	for _, removeLine := range removeLines {
		cleaned = bytes.ReplaceAll(cleaned, removeLine, []byte(""))
	}

	// The Kubernetes standard header field `currentTimestamp` serializes weirdly, so filter it out.
	// filtered := bytes.ReplaceAll(buf.Bytes(), []byte("creationTimestamp: null"), []byte(""))

	// cleaned := bytes.ReplaceAll(filtered, []byte("status:\n  loadBalancer: {}\n"), []byte(""))

	out, err := formatManifest(cleaned)
	if err != nil {
		fmt.Printf("can't format manifest: %s", err)
	}

	directory := "datacenters/" + tier.TierName + "/" + app.Namespace + "/"

	_, cerr := os.Stat(directory)
	if os.IsNotExist(cerr) {
		errDir := os.MkdirAll(directory, 0755)
		if errDir != nil {
			panic(err)
		}

	}

	filePath := directory + app.Application + ".yaml"
	if stage != "production" {
		filePath = directory + app.Application + "-" + stage + ".yaml"
	}

	WriteFileOnDisk(out.Bytes(), filePath)

	return filePath
}

// CreatePullRequest create pull request with generated manifests to github repository
func CreatePullRequest(sourceFiles string, PROptions *types.PullRequestOptions) (err error) {
	gc := git.NewClient()

	// Create a tree with what to commit.
	entries := []*github.TreeEntry{}

	// Load each file into the tree.
	for _, fileArg := range strings.Split(sourceFiles, ",") {
		file, content, err := GetFileContent(fileArg)
		if err != nil {
			fmt.Printf("Error while getting file content: %s", err)
			return err
		}
		entries = append(entries, &github.TreeEntry{
			Path:    github.String(file),
			Type:    github.String("blob"),
			Content: github.String(string(content)),
			Mode:    github.String("100644")})
	}

	PROptions.Entries = entries

	if err := gc.NewPullRequest(PROptions); err != nil {
		fmt.Printf("Error while creating the pull request: %s", err)
		return err
	}

	return nil
}

// LoadConfiguration returns application config from local or remote configuration.json file
func LoadConfiguration(local bool, organization, repositoryName, branch string) []*types.Configuration {
	var configResponse = make([]*types.Configuration, 0)
	gc := git.NewClient()

	if local {
		configResponse = ReadJSONLocalToStruct()
	} else {
		if repositoryName != "" {
			configResponse = gc.ReadJSONFromRepoToStruct(
				organization,
				repositoryName,
				branch)
		}
	}

	return configResponse
}
