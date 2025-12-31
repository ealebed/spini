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
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NewDeployment return k8s deployment object
func NewDeployment(config *Configuration, tier *Datacenter, stage, organization string) *appsv1.Deployment {
	application := config.Application
	if stage != stageProduction {
		application = config.Application + "-" + stage
	}

	deployment := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        application,
			Namespace:   config.Namespace,
			Annotations: defaultAnnotations(organization, stage, config.Owners),
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(tier.Replicas),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": application,
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": application,
					},
				},
				Spec: apiv1.PodSpec{
					ServiceAccountName:            application,
					Containers:                    newContainer(config, tier, organization, stage),
					Affinity:                      newAffinity(application, config.NodePool),
					Tolerations:                   newToleration(config.NodePool),
					Volumes:                       newVolume(config.Application, config.DependsOn),
					TerminationGracePeriodSeconds: int64Ptr(20),
					PriorityClassName:             tier.PodPriority,
				},
			},
			Strategy: NewDeploymentStrategy(config.Strategy),
		},
	}

	if tier.ProgressDeadline != 0 {
		deployment.Spec.ProgressDeadlineSeconds = int32Ptr(tier.ProgressDeadline)
	}

	if tier.PodPriority == "high-priority" {
		deployment.Spec.Template.Spec.TerminationGracePeriodSeconds = int64Ptr(60)
	}

	if dependencyContains(config.DependsOn, "maxmind") {
		deployment.Spec.Template.Spec.InitContainers = []apiv1.Container{
			{
				Name:    "data-container",
				Image:   dockerHubUrl + organization + "/maxmind-geoip",
				Command: []string{"cp", "-a", "/usr/share/GeoIP/.", "/tmp"},
				VolumeMounts: []apiv1.VolumeMount{
					{
						Name:      "geoip-files",
						MountPath: "/tmp",
					},
				},
			},
		}
	}

	if tier.ChaosMonkey.Enabled {
		deployment.Labels = newMetadataLabels(application, tier.ChaosMonkey)
		deployment.Spec.Template.Labels = newTemplateLabels(application)
	}

	return deployment
}

func int32Ptr(i int32) *int32 { return &i }

func int64Ptr(i int64) *int64 { return &i }
