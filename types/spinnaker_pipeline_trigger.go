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

// Triggers represents full trigers config
type Trigger struct {
	Account             string   `json:"account,omitempty"`
	Application         string   `json:"application,omitempty"`
	Branch              string   `json:"branch,omitempty"`
	Enabled             bool     `json:"enabled"`
	ExpectedArtifactIds []string `json:"expectedArtifactIds"`
	Organization        string   `json:"organization,omitempty"`
	Pipeline            string   `json:"pipeline,omitempty"`
	Project             string   `json:"project,omitempty"`
	Registry            string   `json:"registry,omitempty"`
	Repository          string   `json:"repository,omitempty"`
	RunAsUser           string   `json:"runAsUser"`
	Secret              string   `json:"secret,omitempty"`
	Slug                string   `json:"slug,omitempty"`
	Source              string   `json:"source,omitempty"`
	Status              []string `json:"status,omitempty"`
	Tag                 string   `json:"tag,omitempty"`
	Type                string   `json:"type"`
}

// newDockerTrigger return Trigger object with default values for docker hub trigger type
func newDockerTrigger(organization, dockerImage, owner string, enabled bool) *Trigger {
	return &Trigger{
		Account:             organization,
		Enabled:             enabled,
		ExpectedArtifactIds: []string{organization + "/" + dockerImage},
		Organization:        organization,
		Registry:            "index.docker.io",
		Repository:          organization + "/" + dockerImage,
		RunAsUser:           owner + "-service-account@" + organization + ".com",
		Tag:                 "^\\d{2}\\.\\d{2}\\.\\d{2}\\-\\d{2}\\.\\d{2}$",
		Type:                "docker",
	}
}

// newGitTrigger return Trigger object with default values for github trigger type
func newGitTrigger(organization, repositoryName, owner string, expectedArtifacts []string) *Trigger {
	return &Trigger{
		Branch:              "master",
		Enabled:             true,
		ExpectedArtifactIds: expectedArtifacts,
		Project:             organization,
		RunAsUser:           owner + "-service-account@" + organization + ".com",
		Slug:                repositoryName,
		Source:              "github",
		Type:                "git",
	}
}

// newPipelineTrigger return Trigger object with default values for spinnaker pipeline trigger type
func newPipelineTrigger(organization, application, owner, parentPipelineId string) *Trigger {
	return &Trigger{
		Application:         application,
		Enabled:             true,
		ExpectedArtifactIds: []string{},
		Pipeline:            parentPipelineId,
		RunAsUser:           owner + "-service-account@" + organization + ".com",
		Status:              []string{"successful"},
		Type:                "pipeline",
	}
}
