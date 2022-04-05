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

// Option contains the value of the option in a given pipeline parameter
type Option struct {
	Value string `json:"value,omitempty"`
}

// Parameter is a parameter declaration for a spinnaker pipeline config
type Parameter struct {
	ID          string   `json:"id"`
	Default     string   `json:"default"`
	Description string   `json:"description"`
	HasOptions  bool     `json:"hasOptions,omitempty"`
	Label       string   `json:"label"`
	Name        string   `json:"name"`
	Options     []Option `json:"options,omitempty"`
	Pinned      bool     `json:"pinned"`
	Required    bool     `json:"required"`
}
