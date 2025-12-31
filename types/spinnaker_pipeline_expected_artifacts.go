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

// PipelineExpectedArtifact represents spinnaker pipeline expected artifact config
type PipelineExpectedArtifact struct {
	DefaultArtifact    *PipelineArtifact `json:"defaultArtifact"`
	DisplayName        string            `json:"displayName"`
	ID                 string            `json:"id"`
	MatchArtifact      *PipelineArtifact `json:"matchArtifact"`
	UseDefaultArtifact bool              `json:"useDefaultArtifact"`
	UsePriorArtifact   bool              `json:"usePriorArtifact"`
}

// newDockerPipelineExpectedArtifact return new expected docker image artifact
func newDockerPipelineExpectedArtifact(organization, image, version string) *PipelineExpectedArtifact {
	return &PipelineExpectedArtifact{
		DefaultArtifact: &PipelineArtifact{
			ArtifactAccount: "docker-registry",
			Name:            dockerHubUrl + organization + "/" + image,
			Reference:       dockerHubUrl + organization + "/" + image + ":" + version,
			Type:            "docker/image",
			Version:         version,
		},
		DisplayName: dockerHubUrl + organization + "/" + image,
		ID:          organization + "/" + image,
		MatchArtifact: &PipelineArtifact{
			ArtifactAccount: "docker-registry",
			Name:            dockerHubUrl + organization + "/" + image,
			Type:            "docker/image",
		},
		UseDefaultArtifact: true,
		UsePriorArtifact:   true,
	}
}

// newManifestPipelineExpectedArtifact return new expected k8s manifest artifact
func newManifestPipelineExpectedArtifact(githubContentUrl, relativePath string) *PipelineExpectedArtifact {
	return &PipelineExpectedArtifact{
		DefaultArtifact: &PipelineArtifact{
			ArtifactAccount: "spinnaker-github-token",
			Name:            relativePath,
			Reference:       githubContentUrl + relativePath,
			Type:            "github/file",
			Version:         "master",
		},
		DisplayName: relativePath,
		ID:          relativePath,
		MatchArtifact: &PipelineArtifact{
			ArtifactAccount: "spinnaker-github-token",
			CustomKind:      true,
			Name:            relativePath,
			Type:            "github/file",
		},
		UseDefaultArtifact: true,
		UsePriorArtifact:   false,
	}
}
