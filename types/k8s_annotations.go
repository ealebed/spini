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

// defaultAnnotations return k8s annotations object with default values for deployment
func defaultAnnotations(organization, stage, owner string) map[string]string {
	return map[string]string{
		"moniker.spinnaker.io/stack":                 stage,
		"service." + organization + ".dev/generated": "spini/v1",
		"service." + organization + ".dev/owners":    owner,
	}
}
