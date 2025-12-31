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

package types

import (
	"strings"

	dha "github.com/ealebed/dha/pkg/dockerhub"
)

const stageProduction = "production"

// PipelineConfig represents full pipeline config - fields for the top level object of a spinnaker
// pipeline. Mostly used for constructing JSON
type Pipeline struct {
	AppConfig            map[string]interface{}      `json:"appConfig,omitempty"`
	Application          string                      `json:"application"`
	Disabled             bool                        `json:"disabled,omitempty"`
	ExpectedArtifacts    []*PipelineExpectedArtifact `json:"expectedArtifacts,omitempty"`
	ID                   string                      `json:"id,omitempty"`
	Index                int                         `json:"index,omitempty"`
	KeepWaitingPipelines bool                        `json:"keepWaitingPipelines"`
	LastModifiedBy       string                      `json:"lastModifiedBy"`
	LimitConcurrent      bool                        `json:"limitConcurrent"`
	Name                 string                      `json:"name"`
	Notifications        []Notification              `json:"notifications,omitempty"`
	ParameterConfig      *[]*Parameter               `json:"parameterConfig,omitempty"`
	Roles                []string                    `json:"roles,omitempty"`
	SpelEvaluator        string                      `json:"spelEvaluator,omitempty"`
	Stages               []*Stage                    `json:"stages,omitempty"`
	Triggers             []*Trigger                  `json:"triggers,omitempty"`
	UpdateTs             string                      `json:"updateTs,omitempty"`
}

// NewBuildPipeline return build pipeline with default values
func NewBuildPipeline(pipe *Configuration) *Pipeline {
	return &Pipeline{
		Application:          pipe.Application,
		Disabled:             false,
		KeepWaitingPipelines: false,
		LastModifiedBy:       pipe.OwnerEmail,
		LimitConcurrent:      true,
		Name:                 "build-image",
		SpelEvaluator:        "v4",
		Stages:               []*Stage{defaultJenkinsStage(pipe.JenkinsJobName, pipe.Application)},
		Triggers:             []*Trigger{},
	}
}

// NewPromotePipeline return promote to stage pipeline with default values
func NewPromotePipeline(pipe *Configuration, pipeValues map[string]interface{}) *Pipeline {
	var organization = pipeValues["organization"].(string)

	return &Pipeline{
		AppConfig:            map[string]interface{}{},
		Application:          pipe.Application,
		ExpectedArtifacts:    []*PipelineExpectedArtifact{},
		Disabled:             false,
		ID:                   pipeValues["id"].(string),
		Index:                0,
		KeepWaitingPipelines: false,
		LastModifiedBy:       pipe.OwnerEmail,
		LimitConcurrent:      true,
		Name:                 "promote-to-" + pipeValues["stage"].(string),
		Roles:                []string{"devops", pipe.Owners},
		SpelEvaluator:        "v4",
		Stages:               []*Stage{defaultPromoteStage(pipeValues["stage"].(string))},
		Triggers:             []*Trigger{newPipelineTrigger(organization, pipe.Application, pipe.Owners, pipeValues["parentPipelineId"].(string))},
	}
}

// NewDeployPipeline return deploy to DC pipeline with default values
func NewDeployPipeline(pipe *Configuration, pipeValues map[string]interface{}) *Pipeline {
	var organization = pipeValues["organization"].(string)
	var githubRepositoryName = pipeValues["githubRepositoryName"].(string)
	var githubContentUrl = "https://api.github.com/repos/" + organization + "/" + githubRepositoryName + "/contents/"

	var fullListStageRefIds = []string{}
	var requiredArtifactIds = []string{organization + "/" + pipe.DockerImage}

	var manifestPath string
	var expectedArtifacts = []*PipelineExpectedArtifact{}
	var stages = []*Stage{}
	var triggers = []*Trigger{}
	var expectedArtifactIds = []string{}

	if pipe.Version == "" {
		tags, _ := dha.NewClient(organization, "").ListTags(pipe.DockerImage)
		pipe.Version = tags[0].Name
	}

	if dependencyContains(pipe.DependsOn, "maxmind") {
		tags, _ := dha.NewClient(organization, "").ListTags("maxmind-geoip")
		maxmindDefaultTag := tags[0].Name

		expectedArtifacts = append(expectedArtifacts, newDockerPipelineExpectedArtifact(
			organization,
			"maxmind-geoip",
			maxmindDefaultTag))
		requiredArtifactIds = append(requiredArtifactIds, organization+"/maxmind-geoip")
		triggers = append(triggers, newDockerTrigger(
			organization,
			"maxmind-geoip",
			pipe.Owners,
			true))
	}

	if pipe.Namespace != "default" {
		manifestPath = "datacenters/" + pipeValues["cluster"].(string) + "/" + pipe.Namespace + "/_namespace.yaml"
		expectedArtifacts = append(expectedArtifacts, newManifestPipelineExpectedArtifact(githubContentUrl, manifestPath))
		stages = append(stages, defaultDeployManifestStage(
			pipeValues["cluster"].(string),
			pipe.Application,
			pipe.Namespace,
			manifestPath,
			[]string{},
			[]string{}))
		expectedArtifactIds = append(expectedArtifactIds, manifestPath)
		fullListStageRefIds = append(fullListStageRefIds, "Deploy "+manifestPath)
	}

	for _, envFile := range pipe.EnvFrom {
		manifestPath = "datacenters/_commons/" + envFile + ".yaml"

		expectedArtifacts = append(expectedArtifacts, newManifestPipelineExpectedArtifact(githubContentUrl, manifestPath))
		stages = append(stages, defaultDeployManifestStage(
			pipeValues["cluster"].(string),
			pipe.Application,
			pipe.Namespace,
			manifestPath,
			[]string{},
			[]string{}))
		expectedArtifactIds = append(expectedArtifactIds, manifestPath)
		fullListStageRefIds = append(fullListStageRefIds, "Deploy "+manifestPath)
	}

	for _, dependency := range pipe.DependsOn {
		if !strings.HasSuffix(dependency.Name, "-config") {
			continue
		}
		manifestPath = "datacenters/_commons/" + dependency.Name + ".yaml"

		expectedArtifacts = append(expectedArtifacts, newManifestPipelineExpectedArtifact(githubContentUrl, manifestPath))
		stages = append(stages, defaultDeployManifestStage(
			pipeValues["cluster"].(string),
			pipe.Application,
			pipe.Namespace,
			manifestPath,
			[]string{},
			[]string{}))
		expectedArtifactIds = append(expectedArtifactIds, manifestPath)
		fullListStageRefIds = append(fullListStageRefIds, "Deploy "+manifestPath)
	}

	if pipeValues["stage"].(string) == stageProduction {
		manifestPath = "datacenters/" + pipeValues["cluster"].(string) + "/" + pipe.Namespace + "/" + pipe.Application + ".yaml"
	} else {
		manifestPath = "datacenters/" + pipeValues["cluster"].(string) + "/" + pipe.Namespace + "/" + pipe.Application + "-" + pipeValues["stage"].(string) + ".yaml"
	}

	expectedArtifacts = append(expectedArtifacts,
		newDockerPipelineExpectedArtifact(organization, pipe.DockerImage, pipe.Version),
		newManifestPipelineExpectedArtifact(githubContentUrl, manifestPath))
	expectedArtifactIds = append(expectedArtifactIds,
		manifestPath)
	stages = append(stages, defaultDeployManifestStage(
		pipeValues["cluster"].(string),
		pipe.Application,
		pipe.Namespace,
		manifestPath,
		fullListStageRefIds,
		requiredArtifactIds))
	triggers = append(triggers, newDockerTrigger(
		organization,
		pipe.DockerImage,
		pipe.Owners,
		pipeValues["dockerTriggerEnabled"].(bool)))
	triggers = append(triggers, newGitTrigger(
		organization,
		pipeValues["githubRepositoryName"].(string),
		pipe.Owners, expectedArtifactIds))

	if pipeValues["GeneratePromotePipeline"].(bool) {
		triggers = append(triggers, newPipelineTrigger(
			organization,
			pipe.Application,
			pipe.Owners,
			pipeValues["parentPipelineId"].(string)))
	}

	var notification = Notification{
		Address: pipe.SlackChannel,
		Level:   "pipeline",
		Message: map[string]NotificationMessage{
			"pipeline.complete": {
				Text: "*Deploy* ${ trigger['artifacts'].?[type == 'docker/image'].![reference] }\n*User:* ${ trigger['user'] }",
			},
			"pipeline.failed": {
				Text: "*Deploy* ${ trigger['artifacts'].?[type == 'docker/image'].![reference] }\n*User:* ${ trigger['user'] }",
			},
			"pipeline.starting": {
				Text: "*Deploy* ${ trigger['artifacts'].?[type == 'docker/image'].![reference] }\n*User:* ${ trigger['user'] }",
			},
		},
		Type: "slack",
		When: []string{"pipeline.starting", "pipeline.failed", "pipeline.complete"},
	}

	return &Pipeline{
		Application:          pipe.Application,
		Disabled:             false,
		ExpectedArtifacts:    expectedArtifacts,
		ID:                   pipeValues["id"].(string),
		Index:                0,
		KeepWaitingPipelines: false,
		LastModifiedBy:       pipe.OwnerEmail,
		LimitConcurrent:      true,
		Name:                 "deploy-" + pipeValues["cluster"].(string) + "-dc(" + pipeValues["stage"].(string) + ")",
		Roles:                []string{"devops", pipe.Owners},
		ParameterConfig:      &[]*Parameter{},
		Notifications:        []Notification{notification},
		Stages:               stages,
		Triggers:             triggers,
	}
}
