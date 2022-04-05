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

const (
	dockerHubUrl = "index.docker.io/"
)

type Configuration struct {
	Application                       string          `json:"application"`
	DockerImage                       string          `json:"image,omitempty"`
	Profiles                          *[]*Profile     `json:"profiles,omitempty"`
	EnvFrom                           []string        `json:"envFrom,omitempty"`
	DependsOn                         []DependsOn     `json:"dependsOn,omitempty"`
	Type                              string          `json:"type"`
	Owners                            string          `json:"owners,omitempty"`
	OwnerEmail                        string          `json:"ownerEmail,omitempty"`
	NodePool                          string          `json:"nodePool,omitempty"`
	Namespace                         string          `json:"namespace,omitempty"`
	SlackChannel                      string          `json:"slackChannel,omitempty"`
	JenkinsJobName                    string          `json:"jenkinsJobName,omitempty"`
	Ports                             []Port          `json:"ports,omitempty"`
	Strategy                          *DeployStrategy `json:"strategy,omitempty"`
	ChaosMonkey                       *ChaosMonkey    `json:"chaosMonkey,omitempty"`
	Version                           string          `json:"version,omitempty"`
	RestrictExecutionDuringTimeWindow bool            `json:"restrictExecutionDuringTimeWindow,omitempty"`
	SkipAutogeneration                bool            `json:"skipAutogeneration,omitempty"`
}

type Profile struct {
	ProfileName string         `json:"profileName"`
	Datacenters *[]*Datacenter `json:"datacenters"`
}

type Datacenter struct {
	TierName         string                `json:"tierName"`
	Replicas         int32                 `json:"replicas"`
	NodePool         string                `json:"nodePool,omitempty"`
	Env              *[]EnvVar             `json:"env"`
	EnvFrom          []string              `json:"envFrom,omitempty"`
	Command          []string              `json:"command,omitempty"`
	Resources        *ResourceRequirements `json:"resources"`
	ProgressDeadline int32                 `json:"progressDeadline,omitempty"`
	LivenessProbe    *Probe                `json:"livenessProbe,omitempty"`
	ReadinessProbe   *Probe                `json:"readinessProbe,omitempty"`
	StartupProbe     *Probe                `json:"startupProbe,omitempty"`
	PodPriority      string                `json:"podPriority,omitempty"`
	ChaosMonkey      *ChaosMonkey          `json:"chaosMonkey,omitempty"`
	Version          string                `json:"version,omitempty"`
}

type EnvVar struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type DeployStrategy struct {
	// Type of deployment. Can be "Recreate" or "RollingUpdate". Default is RollingUpdate.
	Type string `json:"type,omitempty"`

	// Rolling update config params. Present only if DeployStrategyType = RollingUpdate.
	RollingUpdate *RollingUpdateDeployment `json:"rollingUpdate,omitempty"`
}

type RollingUpdateDeployment struct {
	// The maximum number of pods that can be unavailable during the update.
	// Value can be ONLY a percentage of desired pods (ex: 10%).
	// Absolute number is calculated from percentage by rounding down.
	// This can not be 0 if MaxSurge is 0.
	// Defaults to 25%.
	// Example: when this is set to 30%, the old ReplicaSet can be scaled down to 70% of desired pods
	// immediately when the rolling update starts. Once new pods are ready, old ReplicaSet
	// can be scaled down further, followed by scaling up the new ReplicaSet, ensuring
	// that the total number of pods available at all times during the update is at
	// least 70% of desired pods.
	MaxUnavailable string `json:"maxUnavailable,omitempty"`

	// The maximum number of pods that can be scheduled above the desired number of
	// pods.
	// Value can be ONLY a percentage of desired pods (ex: 10%).
	// This can not be 0 if MaxUnavailable is 0.
	// Absolute number is calculated from percentage by rounding up.
	// Defaults to 25%.
	// Example: when this is set to 30%, the new ReplicaSet can be scaled up immediately when
	// the rolling update starts, such that the total number of old and new pods do not exceed
	// 130% of desired pods. Once old pods have been killed,
	// new ReplicaSet can be scaled up further, ensuring that total number of pods running
	// at any time during the update is at most 130% of desired pods.
	MaxSurge string `json:"maxSurge,omitempty"`
}

type ResourceRequirements struct {
	Limits   *ResourceList `json:"limits"`
	Requests *ResourceList `json:"requests"`
}

type ResourceList struct {
	CPU    string `json:"cpu"`
	Memory string `json:"memory"`
}

type Probe struct {
	Type             string `json:"type"`
	Path             string `json:"path"`
	Delay            int32  `json:"delay"`
	Port             int    `json:"port"`
	TimeoutSeconds   int32  `json:"timeoutSeconds"`
	PeriodSeconds    int32  `json:"periodSeconds"`
	SuccessThreshold int32  `json:"successThreshold"`
	FailureThreshold int32  `json:"failureThreshold"`
}

type ChaosMonkey struct {
	Enabled   bool   `json:"enabled"`
	MTBF      string `json:"mtbf"`
	KillMode  string `json:"killMode"`
	KillValue string `json:"killValue"`
}

type DependsOn struct {
	Name string `json:"name"`
	URL  string `json:"url,omitempty"`
}

type Port struct {
	Name          string `json:"name"`
	ContainerPort int32  `json:"containerPort"`
}

// dependencyContains checks if a string is present in a dependencies slice as name
func dependencyContains(d []DependsOn, str string) bool {
	for _, v := range d {
		if v.Name == str {
			return true
		}
	}

	return false
}
