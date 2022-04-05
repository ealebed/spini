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

// newVolumeMount return k8s volume mount objects list
func newVolumeMount(application string, dependencies []DependsOn) []apiv1.VolumeMount {
	listVolumeMounts := []apiv1.VolumeMount{}

	for _, dependency := range dependencies {
		if strings.HasSuffix(dependency.Name, "-config") {
			dependency := strings.Replace(dependency.Name, "-config", "", 1)
			listVolumeMounts = append(listVolumeMounts, apiv1.VolumeMount{
				Name:      dependency,
				MountPath: "/app/conf/" + dependency + ".conf",
				SubPath:   dependency + ".conf",
			})
		}

		if dependency.Name == "GoogleCloudStorage" {
			listVolumeMounts = append(listVolumeMounts, apiv1.VolumeMount{
				Name:      "google-cloud-" + application,
				ReadOnly:  true,
				MountPath: "/google-cloud-" + application,
			})
		}

		if dependency.Name == "maxmind" {
			listVolumeMounts = append(listVolumeMounts, apiv1.VolumeMount{
				Name:      "geoip-files",
				MountPath: "/usr/share/GeoIP",
			})
		}
	}

	sort.Slice(listVolumeMounts, func(i, j int) bool { return listVolumeMounts[i].Name < listVolumeMounts[j].Name })

	return listVolumeMounts
}
