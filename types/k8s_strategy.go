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
	"k8s.io/apimachinery/pkg/util/intstr"
)

// defaultDeploymentStrategy return k8s strategy object for deployment with default values
func defaultDeploymentStrategy() appsv1.DeploymentStrategy {
	return appsv1.DeploymentStrategy{
		Type: "RollingUpdate",
		RollingUpdate: &appsv1.RollingUpdateDeployment{
			MaxSurge: &intstr.IntOrString{
				Type:   intstr.String,
				StrVal: "25%",
			},
			MaxUnavailable: &intstr.IntOrString{
				Type:   intstr.String,
				StrVal: "25%",
			},
		},
	}
}

// NewDeploymentStrategy return k8s strategy for deployment
func NewDeploymentStrategy(s *DeployStrategy) appsv1.DeploymentStrategy {
	strategy := appsv1.DeploymentStrategy{}

	if s != nil {
		strategy.Type = appsv1.DeploymentStrategyType(s.Type)
		if s.Type == "RollingUpdate" {
			strategy.RollingUpdate = &appsv1.RollingUpdateDeployment{
				MaxSurge: &intstr.IntOrString{
					Type:   intstr.String,
					StrVal: s.RollingUpdate.MaxSurge,
				},
				MaxUnavailable: &intstr.IntOrString{
					Type:   intstr.String,
					StrVal: s.RollingUpdate.MaxUnavailable,
				},
			}
		}
	} else {
		strategy = defaultDeploymentStrategy()
	}

	return strategy
}
