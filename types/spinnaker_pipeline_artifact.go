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

// PipelineArtifact represents default spinnaker pipeline artifact config
type PipelineArtifact struct {
	ArtifactAccount string            `json:"artifactAccount,omitempty"`
	CustomKind      bool              `json:"customKind,omitempty"`
	ID              string            `json:"id,omitempty"`
	Location        string            `json:"location,omitempty"`
	Metadata        map[string]string `json:"metadata,omitempty"`
	Name            string            `json:"name,omitempty"`
	Reference       string            `json:"reference,omitempty"`
	Type            string            `json:"type,omitempty"`
	Version         string            `json:"version,omitempty"`
}
