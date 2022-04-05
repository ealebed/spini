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
	"sort"
	"strings"

	apiv1 "k8s.io/api/core/v1"
)

// newVolume return k8s volume objects list
func newVolume(application string, dependencies []DependsOn) []apiv1.Volume {
	listVolumes := []apiv1.Volume{}

	for _, dependency := range dependencies {
		if strings.HasSuffix(dependency.Name, "-config") {
			dependency := strings.Replace(dependency.Name, "-config", "", 1)
			listVolumes = append(listVolumes, apiv1.Volume{
				Name: dependency,
				VolumeSource: apiv1.VolumeSource{
					ConfigMap: &apiv1.ConfigMapVolumeSource{
						LocalObjectReference: apiv1.LocalObjectReference{
							Name: dependency,
						},
					},
				},
			})
		}

		if dependency.Name == "GoogleCloudStorage" {
			listVolumes = append(listVolumes, apiv1.Volume{
				Name: "google-cloud-" + application,
				VolumeSource: apiv1.VolumeSource{
					Secret: &apiv1.SecretVolumeSource{
						SecretName: "google-cloud-" + application,
					},
				},
			})
		}

		if dependency.Name == "maxmind" {
			listVolumes = append(listVolumes, apiv1.Volume{
				Name: "geoip-files",
				VolumeSource: apiv1.VolumeSource{
					EmptyDir: &apiv1.EmptyDirVolumeSource{},
				},
			})
		}
	}

	sort.Slice(listVolumes, func(i, j int) bool { return listVolumes[i].Name < listVolumes[j].Name })

	return listVolumes
}
