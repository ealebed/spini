package types

import (
	"strings"
	"testing"
)

func TestDefaultPromoteStage(t *testing.T) {
	tests := []struct {
		name     string
		stage    string
		validate func(*testing.T, *Stage)
	}{
		{
			name:  "production stage",
			stage: "production",
			validate: func(t *testing.T, stage *Stage) {
				if stage.Type != "manualJudgment" {
					t.Errorf("Expected type 'manualJudgment', got %q", stage.Type)
				}
				if stage.Name != "Manual Judgment" {
					t.Errorf("Expected name 'Manual Judgment', got %q", stage.Name)
				}
				if !stage.FailPipeline {
					t.Error("Expected FailPipeline to be true")
				}
				if stage.RefID != "1" {
					t.Errorf("Expected RefID '1', got %q", stage.RefID)
				}
				if !stage.PropagateAuthenticationContext {
					t.Error("Expected PropagateAuthenticationContext to be true")
				}
				if stage.StageTimeoutMs != 36000000 {
					t.Errorf("Expected StageTimeoutMs 36000000, got %d", stage.StageTimeoutMs)
				}
				if !strings.Contains(stage.Instructions, "production") {
					t.Errorf("Expected instructions to contain 'production', got %q", stage.Instructions)
				}
			},
		},
		{
			name:  "staging stage",
			stage: "staging",
			validate: func(t *testing.T, stage *Stage) {
				if stage.Type != "manualJudgment" {
					t.Errorf("Expected type 'manualJudgment', got %q", stage.Type)
				}
				if !strings.Contains(stage.Instructions, "staging") {
					t.Errorf("Expected instructions to contain 'staging', got %q", stage.Instructions)
				}
			},
		},
		{
			name:  "empty stage name",
			stage: "",
			validate: func(t *testing.T, stage *Stage) {
				if stage.Type != "manualJudgment" {
					t.Errorf("Expected type 'manualJudgment', got %q", stage.Type)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := defaultPromoteStage(tt.stage)
			if result == nil {
				t.Fatal("defaultPromoteStage returned nil")
			}
			tt.validate(t, result)
		})
	}
}

func TestDefaultJenkinsStage(t *testing.T) {
	tests := []struct {
		name           string
		jenkinsJobName string
		application    string
		validate       func(*testing.T, *Stage)
	}{
		{
			name:           "parametrised job",
			jenkinsJobName: "parametrised_job",
			application:    "myapp",
			validate: func(t *testing.T, stage *Stage) {
				if stage.Type != "jenkins" {
					t.Errorf("Expected type 'jenkins', got %q", stage.Type)
				}
				if stage.Job != "parametrised_job" {
					t.Errorf("Expected Job 'parametrised_job', got %q", stage.Job)
				}
				if stage.Master != "default-jenkins" {
					t.Errorf("Expected Master 'default-jenkins', got %q", stage.Master)
				}
				if !stage.RestrictExecutionDuringTimeWindow {
					t.Error("Expected RestrictExecutionDuringTimeWindow to be true")
				}
				if stage.Parameters == nil {
					t.Fatal("Expected Parameters to be set")
				}
				if stage.Parameters["MODULE_NAME"] != "myapp" {
					t.Errorf("Expected MODULE_NAME 'myapp', got %q", stage.Parameters["MODULE_NAME"])
				}
				if stage.Parameters["RELEASE_TAG"] != "origin/master" {
					t.Errorf("Expected RELEASE_TAG 'origin/master', got %q", stage.Parameters["RELEASE_TAG"])
				}
			},
		},
		{
			name:           "non-parametrised job",
			jenkinsJobName: "simple_job",
			application:    "myapp",
			validate: func(t *testing.T, stage *Stage) {
				if stage.Type != "jenkins" {
					t.Errorf("Expected type 'jenkins', got %q", stage.Type)
				}
				if stage.Job != "simple_job" {
					t.Errorf("Expected Job 'simple_job', got %q", stage.Job)
				}
				if len(stage.Parameters) != 0 {
					t.Errorf("Expected empty Parameters, got %v", stage.Parameters)
				}
			},
		},
		{
			name:           "empty application name",
			jenkinsJobName: "parametrised_job",
			application:    "",
			validate: func(t *testing.T, stage *Stage) {
				if stage.Parameters["MODULE_NAME"] != "" {
					t.Errorf("Expected empty MODULE_NAME, got %q", stage.Parameters["MODULE_NAME"])
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := defaultJenkinsStage(tt.jenkinsJobName, tt.application)
			if result == nil {
				t.Fatal("defaultJenkinsStage returned nil")
			}
			tt.validate(t, result)
		})
	}
}

func TestDefaultDeployManifestStage(t *testing.T) {
	tests := []struct {
		name             string
		cluster          string
		application      string
		namespace        string
		manifestPath     string
		stageRefIds      []string
		stageArtifactIds []string
		validate         func(*testing.T, *Stage)
	}{
		{
			name:             "complete stage",
			cluster:          "gke1",
			application:      "myapp",
			namespace:        "default",
			manifestPath:     "datacenters/gke1/default/myapp.yaml",
			stageRefIds:      []string{"1", "2"},
			stageArtifactIds: []string{"artifact1", "artifact2"},
			validate: func(t *testing.T, stage *Stage) {
				if stage.Type != "deployManifest" {
					t.Errorf("Expected type 'deployManifest', got %q", stage.Type)
				}
				if stage.Account != "gke1" {
					t.Errorf("Expected Account 'gke1', got %q", stage.Account)
				}
				if stage.CloudProvider != "kubernetes" {
					t.Errorf("Expected CloudProvider 'kubernetes', got %q", stage.CloudProvider)
				}
				if stage.NamespaceOverride != "default" {
					t.Errorf("Expected NamespaceOverride 'default', got %q", stage.NamespaceOverride)
				}
				if stage.ManifestArtifactID != "datacenters/gke1/default/myapp.yaml" {
					t.Errorf("Expected ManifestArtifactID 'datacenters/gke1/default/myapp.yaml', got %q", stage.ManifestArtifactID)
				}
				if stage.Source != "artifact" {
					t.Errorf("Expected Source 'artifact', got %q", stage.Source)
				}
				if stage.Moniker == nil {
					t.Fatal("Expected Moniker to be set")
				}
				if stage.Moniker.App != "myapp" {
					t.Errorf("Expected Moniker.App 'myapp', got %q", stage.Moniker.App)
				}
				if len(stage.RequisiteStageRefIds) != 2 {
					t.Errorf("Expected 2 RequisiteStageRefIds, got %d", len(stage.RequisiteStageRefIds))
				}
				if len(stage.RequiredArtifactIds) != 2 {
					t.Errorf("Expected 2 RequiredArtifactIds, got %d", len(stage.RequiredArtifactIds))
				}
				if stage.TrafficManagement == nil {
					t.Fatal("Expected TrafficManagement to be set")
				}
			},
		},
		{
			name:             "empty ref ids and artifact ids",
			cluster:          "gke2",
			application:      "otherapp",
			namespace:        "production",
			manifestPath:     "datacenters/gke2/production/otherapp.yaml",
			stageRefIds:      []string{},
			stageArtifactIds: []string{},
			validate: func(t *testing.T, stage *Stage) {
				if len(stage.RequisiteStageRefIds) != 0 {
					t.Errorf("Expected empty RequisiteStageRefIds, got %v", stage.RequisiteStageRefIds)
				}
				if len(stage.RequiredArtifactIds) != 0 {
					t.Errorf("Expected empty RequiredArtifactIds, got %v", stage.RequiredArtifactIds)
				}
			},
		},
		{
			name:             "nil ref ids and artifact ids",
			cluster:          "gke3",
			application:      "testapp",
			namespace:        "test",
			manifestPath:     "datacenters/gke3/test/testapp.yaml",
			stageRefIds:      nil,
			stageArtifactIds: nil,
			validate: func(t *testing.T, stage *Stage) {
				// In Go, nil slices are valid and equivalent to empty slices
				if stage.RequisiteStageRefIds != nil && len(stage.RequisiteStageRefIds) != 0 {
					t.Errorf("Expected RequisiteStageRefIds to be nil or empty, got %v", stage.RequisiteStageRefIds)
				}
				if stage.RequiredArtifactIds != nil && len(stage.RequiredArtifactIds) != 0 {
					t.Errorf("Expected RequiredArtifactIds to be nil or empty, got %v", stage.RequiredArtifactIds)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := defaultDeployManifestStage(tt.cluster, tt.application, tt.namespace, tt.manifestPath, tt.stageRefIds, tt.stageArtifactIds)
			if result == nil {
				t.Fatal("defaultDeployManifestStage returned nil")
			}
			tt.validate(t, result)
		})
	}
}
