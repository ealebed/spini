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

// newContainer return set of k8s container objects
func newContainer(config *Configuration, tier *Datacenter, organization, stage string) []apiv1.Container {
	application := config.Application
	if stage != stageProduction {
		application = config.Application + "-" + stage
	}

	listContainers := []apiv1.Container{}
	containerPorts := []apiv1.ContainerPort{}
	containerEnvs := []apiv1.EnvVar{
		{
			Name:  "API_NAME",
			Value: application,
		},
		{
			Name:  "DC_NAME",
			Value: tier.TierName,
		},
	}
	containerEnvsFrom := []apiv1.EnvFromSource{}

	for _, port := range config.Ports {
		containerPort := apiv1.ContainerPort{
			Name:          port.Name,
			ContainerPort: port.ContainerPort,
		}
		containerPorts = append(containerPorts, containerPort)
	}

	if tier.Env != nil {
		for _, env := range *tier.Env {
			containerEnv := apiv1.EnvVar{
				Name:  env.Name,
				Value: env.Value,
			}
			containerEnvs = append(containerEnvs, containerEnv)
		}
	}
	sort.Slice(containerEnvs, func(i, j int) bool { return containerEnvs[i].Name < containerEnvs[j].Name })

	for _, envFromFile := range append(config.EnvFrom, tier.EnvFrom...) {
		envFrom := strings.Replace(envFromFile, "-configmap", "", 1)
		containerEnvFrom := apiv1.EnvFromSource{
			ConfigMapRef: &apiv1.ConfigMapEnvSource{
				LocalObjectReference: apiv1.LocalObjectReference{
					Name: envFrom,
				},
			},
		}
		containerEnvsFrom = append(containerEnvsFrom, containerEnvFrom)
	}

	listContainers = append(listContainers, apiv1.Container{
		Name:          application,
		Image:         dockerHubUrl + organization + "/" + config.DockerImage,
		Ports:         containerPorts,
		Env:           containerEnvs,
		Resources:     newResourceRequirements(tier),
		VolumeMounts:  newVolumeMount(config.Application, config.DependsOn),
		LivenessProbe: newLivenessProbe(tier),
	})

	if config.Type == "service" {
		listContainers[len(listContainers)-1].ReadinessProbe = newReadinessProbe(tier)
	}

	if tier.StartupProbe != nil {
		listContainers[len(listContainers)-1].StartupProbe = newStartupProbe(tier)
	}

	if containerEnvsFrom != nil {
		listContainers[len(listContainers)-1].EnvFrom = containerEnvsFrom
	}

	if tier.Command != nil {
		listContainers[len(listContainers)-1].Command = tier.Command
	}

	return listContainers
}
