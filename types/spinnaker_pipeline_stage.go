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
)

// StageEnabled represents stage enabled config
type StageEnabled struct {
	Expression string `json:"expression"`
	Type       string `json:"type"`
}

// Stages represents full stages config
type Stage struct {
	Account                           string                     `json:"account,omitempty"`
	CloudProvider                     string                     `json:"cloudProvider,omitempty"`
	Comments                          string                     `json:"comments,omitempty"`
	ContinuePipeline                  bool                       `json:"continuePipeline,omitempty"`
	FailOnFailedExpressions           bool                       `json:"failOnFailedExpressions,omitempty"`
	FailPipeline                      bool                       `json:"failPipeline,omitempty"`
	Instructions                      string                     `json:"instructions,omitempty"`
	Job                               string                     `json:"job,omitempty"`
	JudgmentInputs                    []string                   `json:"judgmentInputs,omitempty"`
	Master                            string                     `json:"master,omitempty"`
	ManifestArtifactAccount           string                     `json:"manifestArtifactAccount,omitempty"`
	ManifestArtifactID                string                     `json:"manifestArtifactId,omitempty"`
	Manifests                         []appsv1.Deployment        `json:"manifests,omitempty"`
	Moniker                           *Moniker                   `json:"moniker,omitempty"`
	Name                              string                     `json:"name"`
	NamespaceOverride                 string                     `json:"namespaceOverride,omitempty"`
	Notifications                     []Notification             `json:"notifications,omitempty"`
	OverrideTimeout                   bool                       `json:"overrideTimeout,omitempty"`
	Parameters                        map[string]string          `json:"parameters,omitempty"`
	PropagateAuthenticationContext    bool                       `json:"propagateAuthenticationContext,omitempty"`
	RefID                             string                     `json:"refId"`
	RequiredArtifactIds               []string                   `json:"requiredArtifactIds,omitempty"`
	RequisiteStageRefIds              []string                   `json:"requisiteStageRefIds"`
	RestrictExecutionDuringTimeWindow bool                       `json:"restrictExecutionDuringTimeWindow,omitempty"`
	RestrictedExecutionWindow         *StageExecutionWindow      `json:"restrictedExecutionWindow,omitempty"`
	SkipExpressionEvaluation          bool                       `json:"skipExpressionEvaluation,omitempty"`
	Source                            string                     `json:"source,omitempty"`
	SendNotifications                 bool                       `json:"sendNotifications,omitempty"`
	StageEnabled                      *StageEnabled              `json:"stageEnabled,omitempty"`
	StageTimeoutMs                    int                        `json:"stageTimeoutMs,omitempty"`
	TrafficManagement                 *PipelineTrafficManagement `json:"trafficManagement,omitempty"`
	Type                              string                     `json:"type"`
}

// defaultPromoteStage return Stage object with default values for promote to stage pipeline
func defaultPromoteStage(stage string) *Stage {
	return &Stage{
		FailPipeline: true,
		Instructions: "Continue deploy docker image\n\n" +
			"\u003cb\u003e${ trigger['artifacts'].?[type == 'docker/image'].![reference] }\u003c/b\u003e\n\n" +
			"to " + stage + " environment?",
		JudgmentInputs:                 []string{},
		Name:                           "Manual Judgment",
		Notifications:                  []Notification{},
		PropagateAuthenticationContext: true,
		RefID:                          "1",
		RequisiteStageRefIds:           []string{},
		StageTimeoutMs:                 36000000,
		Type:                           "manualJudgment",
	}
}

// defaultJenkinsStage return Stage object with default values for jenkins build pipeline
func defaultJenkinsStage(jenkinsJobName, application string) *Stage {
	var Parameters = map[string]string{}

	if jenkinsJobName == "parametrised_job" {
		Parameters = map[string]string{
			"MODULE_NAME": application,
			"RELEASE_TAG": "origin/master",
		}
	}

	return &Stage{
		ContinuePipeline:                  false,
		FailPipeline:                      true,
		Job:                               jenkinsJobName,
		Master:                            "default-jenkins",
		Name:                              "Jenkins",
		Parameters:                        Parameters,
		RefID:                             "1",
		RequisiteStageRefIds:              []string{},
		RestrictExecutionDuringTimeWindow: true,
		RestrictedExecutionWindow:         defaultStageExecutionWindow(),
		Type:                              "jenkins",
	}
}

// defaultDeployManifestStage return Stage object with default values for deploy manifest pipeline
func defaultDeployManifestStage(cluster, application, namespace, manifestPath string, stageRefIds, stageArtifactIds []string) *Stage {
	return &Stage{
		Account:            cluster,
		CloudProvider:      "kubernetes",
		ManifestArtifactID: manifestPath,
		Moniker: &Moniker{
			App: application,
		},
		Name:                     "Deploy " + manifestPath,
		NamespaceOverride:        namespace,
		RefID:                    "Deploy " + manifestPath,
		RequiredArtifactIds:      stageArtifactIds,
		RequisiteStageRefIds:     stageRefIds,
		SkipExpressionEvaluation: false,
		Source:                   "artifact",
		TrafficManagement:        defaultPipelineTrafficManagement(),
		Type:                     "deployManifest",
	}
}
