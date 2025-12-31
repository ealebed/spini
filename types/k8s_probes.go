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
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	probeTypeFile   = "file"
	probePathHealth = "/health"
)

// newLivenessProbe return k8s liveness probe object
func newLivenessProbe(tier *Datacenter) *apiv1.Probe { //nolint:dupl // similar structure to newStartupProbe is acceptable
	var probeHandler apiv1.ProbeHandler

	if tier.LivenessProbe.Type == probeTypeFile {
		probeHandler = apiv1.ProbeHandler{
			Exec: &apiv1.ExecAction{
				Command: []string{"cat", "/tmp/live"},
			},
		}
	} else {
		if tier.LivenessProbe.Path == "" {
			tier.LivenessProbe.Path = probePathHealth
		}
		probeHandler = apiv1.ProbeHandler{
			HTTPGet: &apiv1.HTTPGetAction{
				Port: intstr.FromInt(getIntOrDefault(tier.LivenessProbe.Port, 8080)),
				Path: tier.LivenessProbe.Path,
			},
		}
	}

	return &apiv1.Probe{
		ProbeHandler:        probeHandler,
		FailureThreshold:    tier.LivenessProbe.FailureThreshold,
		InitialDelaySeconds: tier.LivenessProbe.Delay,
		PeriodSeconds:       tier.LivenessProbe.PeriodSeconds,
		SuccessThreshold:    tier.LivenessProbe.SuccessThreshold,
		TimeoutSeconds:      tier.LivenessProbe.TimeoutSeconds,
	}
}

// newReadinessProbe return k8s readiness probe object
func newReadinessProbe(tier *Datacenter) *apiv1.Probe {
	if tier.ReadinessProbe == nil {
		tier.ReadinessProbe = tier.LivenessProbe
	}

	var probeHandler apiv1.ProbeHandler

	if tier.LivenessProbe.Type == probeTypeFile {
		probeHandler = apiv1.ProbeHandler{
			Exec: &apiv1.ExecAction{
				Command: []string{"cat", "/tmp/ready"},
			},
		}
	} else {
		if tier.ReadinessProbe.Path == "" {
			tier.ReadinessProbe.Path = probePathHealth
		}
		probeHandler = apiv1.ProbeHandler{
			HTTPGet: &apiv1.HTTPGetAction{
				Port: intstr.FromInt(getIntOrDefault(tier.ReadinessProbe.Port, 8080)),
				Path: tier.ReadinessProbe.Path,
			},
		}
	}

	return &apiv1.Probe{
		ProbeHandler:        probeHandler,
		FailureThreshold:    tier.ReadinessProbe.FailureThreshold,
		InitialDelaySeconds: tier.ReadinessProbe.Delay,
		PeriodSeconds:       tier.ReadinessProbe.PeriodSeconds,
		SuccessThreshold:    tier.ReadinessProbe.SuccessThreshold,
		TimeoutSeconds:      tier.ReadinessProbe.TimeoutSeconds,
	}
}

// newStartupProbe return k8s startup probe object
func newStartupProbe(tier *Datacenter) *apiv1.Probe { //nolint:dupl // similar structure to newLivenessProbe is acceptable
	var probeHandler apiv1.ProbeHandler

	if tier.StartupProbe.Type == probeTypeFile {
		probeHandler = apiv1.ProbeHandler{
			Exec: &apiv1.ExecAction{
				Command: []string{"cat", "/tmp/started"},
			},
		}
	} else {
		if tier.StartupProbe.Path == "" {
			tier.StartupProbe.Path = probePathHealth
		}
		probeHandler = apiv1.ProbeHandler{
			HTTPGet: &apiv1.HTTPGetAction{
				Port: intstr.FromInt(getIntOrDefault(tier.StartupProbe.Port, 8080)),
				Path: tier.StartupProbe.Path,
			},
		}
	}

	return &apiv1.Probe{
		ProbeHandler:        probeHandler,
		FailureThreshold:    tier.StartupProbe.FailureThreshold,
		InitialDelaySeconds: tier.StartupProbe.Delay,
		PeriodSeconds:       tier.StartupProbe.PeriodSeconds,
		SuccessThreshold:    tier.StartupProbe.SuccessThreshold,
		TimeoutSeconds:      tier.StartupProbe.TimeoutSeconds,
	}
}

func getIntOrDefault(value, defaultValue int) int {
	if value > 0 {
		return value
	} else {
		return defaultValue
	}
}
