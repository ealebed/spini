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

// Moniker represents moniker config - component of the V2 Spinnaker Manifest Stage
// that allows users to label assets created by the Spinnaker v2 provider
type Moniker struct {
	App      string `json:"app"`
	Cluster  string `json:"cluster,omitempty"`
	Detail   string `json:"detail,omitempty"`
	Stack    string `json:"stack,omitempty"`
	Sequence string `json:"sequence,omitempty"`
}
