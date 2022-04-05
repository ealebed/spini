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

// newMetadataLabels return k8s metadata labels for deployment object with default chaosMonkey values
func newMetadataLabels(application string, chaosConfig *ChaosMonkey) map[string]string {
	return map[string]string{
		"kube-monkey/enabled":    "enabled",
		"kube-monkey/identifier": application,
		"kube-monkey/kill-mode":  chaosConfig.KillMode,
		"kube-monkey/kill-value": chaosConfig.KillValue,
		"kube-monkey/mtbf":       chaosConfig.MTBF,
	}
}

// newTemplateLabels return k8s template labels for deployment object with default chaosMonkey values
func newTemplateLabels(application string) map[string]string {
	return map[string]string{
		"app":                    application,
		"kube-monkey/enabled":    "enabled",
		"kube-monkey/identifier": application,
	}
}
