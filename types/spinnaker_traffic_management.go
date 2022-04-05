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

// PipelineTrafficManagement represents pipeline traffic management config
type PipelineTrafficManagement struct {
	Enabled bool                              `json:"enabled"`
	Options *PipelineTrafficManagementOptions `json:"options,omitempty"`
}

// defaultPipelineTrafficManagement returns PipelineTrafficManagement object with default values
func defaultPipelineTrafficManagement() *PipelineTrafficManagement {
	return &PipelineTrafficManagement{
		Enabled: false,
		Options: defaultPipelineTrafficManagementOptions(),
	}
}
