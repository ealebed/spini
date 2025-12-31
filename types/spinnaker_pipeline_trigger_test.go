package types

import (
	"testing"
)

func TestNewDockerTrigger(t *testing.T) {
	tests := []struct {
		name         string
		organization string
		dockerImage  string
		owner        string
		enabled      bool
		validate     func(*testing.T, *Trigger)
	}{
		{
			name:         "enabled trigger",
			organization: "myorg",
			dockerImage:  "myapp",
			owner:        "john",
			enabled:      true,
			validate: func(t *testing.T, trigger *Trigger) {
				if trigger.Type != "docker" {
					t.Errorf("Expected Type 'docker', got %q", trigger.Type)
				}
				if !trigger.Enabled {
					t.Error("Expected Enabled to be true")
				}
				if trigger.Account != "myorg" {
					t.Errorf("Expected Account 'myorg', got %q", trigger.Account)
				}
				if trigger.Organization != "myorg" {
					t.Errorf("Expected Organization 'myorg', got %q", trigger.Organization)
				}
				if trigger.Registry != "index.docker.io" {
					t.Errorf("Expected Registry 'index.docker.io', got %q", trigger.Registry)
				}
				if trigger.Repository != "myorg/myapp" {
					t.Errorf("Expected Repository 'myorg/myapp', got %q", trigger.Repository)
				}
				expectedRunAsUser := "john-service-account@myorg.com"
				if trigger.RunAsUser != expectedRunAsUser {
					t.Errorf("Expected RunAsUser %q, got %q", expectedRunAsUser, trigger.RunAsUser)
				}
				if len(trigger.ExpectedArtifactIds) != 1 {
					t.Fatalf("Expected 1 ExpectedArtifactIds, got %d", len(trigger.ExpectedArtifactIds))
				}
				if trigger.ExpectedArtifactIds[0] != "myorg/myapp" {
					t.Errorf("Expected ExpectedArtifactIds[0] 'myorg/myapp', got %q", trigger.ExpectedArtifactIds[0])
				}
				expectedTag := "^\\d{2}\\.\\d{2}\\.\\d{2}\\-\\d{2}\\.\\d{2}$"
				if trigger.Tag != expectedTag {
					t.Errorf("Expected Tag %q, got %q", expectedTag, trigger.Tag)
				}
			},
		},
		{
			name:         "disabled trigger",
			organization: "myorg",
			dockerImage:  "myapp",
			owner:        "jane",
			enabled:      false,
			validate: func(t *testing.T, trigger *Trigger) {
				if trigger.Enabled {
					t.Error("Expected Enabled to be false")
				}
				if trigger.Type != "docker" {
					t.Errorf("Expected Type 'docker', got %q", trigger.Type)
				}
			},
		},
		{
			name:         "empty organization",
			organization: "",
			dockerImage:  "myapp",
			owner:        "john",
			enabled:      true,
			validate: func(t *testing.T, trigger *Trigger) {
				if trigger.Repository != "/myapp" {
					t.Errorf("Expected Repository '/myapp', got %q", trigger.Repository)
				}
				if trigger.ExpectedArtifactIds[0] != "/myapp" {
					t.Errorf("Expected ExpectedArtifactIds[0] '/myapp', got %q", trigger.ExpectedArtifactIds[0])
				}
			},
		},
		{
			name:         "empty docker image",
			organization: "myorg",
			dockerImage:  "",
			owner:        "john",
			enabled:      true,
			validate: func(t *testing.T, trigger *Trigger) {
				if trigger.Repository != "myorg/" {
					t.Errorf("Expected Repository 'myorg/', got %q", trigger.Repository)
				}
			},
		},
		{
			name:         "empty owner",
			organization: "myorg",
			dockerImage:  "myapp",
			owner:        "",
			enabled:      true,
			validate: func(t *testing.T, trigger *Trigger) {
				expectedRunAsUser := "-service-account@myorg.com"
				if trigger.RunAsUser != expectedRunAsUser {
					t.Errorf("Expected RunAsUser %q, got %q", expectedRunAsUser, trigger.RunAsUser)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := newDockerTrigger(tt.organization, tt.dockerImage, tt.owner, tt.enabled)
			if result == nil {
				t.Fatal("newDockerTrigger returned nil")
			}
			tt.validate(t, result)
		})
	}
}

func TestNewGitTrigger(t *testing.T) {
	tests := []struct {
		name              string
		organization      string
		repositoryName    string
		owner             string
		expectedArtifacts []string
		validate          func(*testing.T, *Trigger)
	}{
		{
			name:              "complete trigger",
			organization:      "myorg",
			repositoryName:    "myrepo",
			owner:             "john",
			expectedArtifacts: []string{"artifact1", "artifact2"},
			validate: func(t *testing.T, trigger *Trigger) {
				if trigger.Type != "git" {
					t.Errorf("Expected Type 'git', got %q", trigger.Type)
				}
				if !trigger.Enabled {
					t.Error("Expected Enabled to be true")
				}
				if trigger.Branch != "master" {
					t.Errorf("Expected Branch 'master', got %q", trigger.Branch)
				}
				if trigger.Project != "myorg" {
					t.Errorf("Expected Project 'myorg', got %q", trigger.Project)
				}
				if trigger.Slug != "myrepo" {
					t.Errorf("Expected Slug 'myrepo', got %q", trigger.Slug)
				}
				if trigger.Source != "github" {
					t.Errorf("Expected Source 'github', got %q", trigger.Source)
				}
				expectedRunAsUser := "john-service-account@myorg.com"
				if trigger.RunAsUser != expectedRunAsUser {
					t.Errorf("Expected RunAsUser %q, got %q", expectedRunAsUser, trigger.RunAsUser)
				}
				if len(trigger.ExpectedArtifactIds) != 2 {
					t.Fatalf("Expected 2 ExpectedArtifactIds, got %d", len(trigger.ExpectedArtifactIds))
				}
				if trigger.ExpectedArtifactIds[0] != "artifact1" {
					t.Errorf("Expected ExpectedArtifactIds[0] 'artifact1', got %q", trigger.ExpectedArtifactIds[0])
				}
			},
		},
		{
			name:              "empty expected artifacts",
			organization:      "myorg",
			repositoryName:    "myrepo",
			owner:             "jane",
			expectedArtifacts: []string{},
			validate: func(t *testing.T, trigger *Trigger) {
				if len(trigger.ExpectedArtifactIds) != 0 {
					t.Errorf("Expected empty ExpectedArtifactIds, got %v", trigger.ExpectedArtifactIds)
				}
			},
		},
		{
			name:              "nil expected artifacts",
			organization:      "myorg",
			repositoryName:    "myrepo",
			owner:             "jane",
			expectedArtifacts: nil,
			validate: func(t *testing.T, trigger *Trigger) {
				// In Go, nil slices are valid and equivalent to empty slices
				if trigger.ExpectedArtifactIds != nil && len(trigger.ExpectedArtifactIds) != 0 {
					t.Errorf("Expected ExpectedArtifactIds to be nil or empty, got %v", trigger.ExpectedArtifactIds)
				}
			},
		},
		{
			name:              "empty organization",
			organization:      "",
			repositoryName:    "myrepo",
			owner:             "john",
			expectedArtifacts: []string{"artifact1"},
			validate: func(t *testing.T, trigger *Trigger) {
				if trigger.Project != "" {
					t.Errorf("Expected Project '', got %q", trigger.Project)
				}
				expectedRunAsUser := "john-service-account@.com"
				if trigger.RunAsUser != expectedRunAsUser {
					t.Errorf("Expected RunAsUser %q, got %q", expectedRunAsUser, trigger.RunAsUser)
				}
			},
		},
		{
			name:              "empty repository name",
			organization:      "myorg",
			repositoryName:    "",
			owner:             "john",
			expectedArtifacts: []string{"artifact1"},
			validate: func(t *testing.T, trigger *Trigger) {
				if trigger.Slug != "" {
					t.Errorf("Expected Slug '', got %q", trigger.Slug)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := newGitTrigger(tt.organization, tt.repositoryName, tt.owner, tt.expectedArtifacts)
			if result == nil {
				t.Fatal("newGitTrigger returned nil")
			}
			tt.validate(t, result)
		})
	}
}

func TestNewPipelineTrigger(t *testing.T) {
	tests := []struct {
		name             string
		organization     string
		application      string
		owner            string
		parentPipelineId string
		validate         func(*testing.T, *Trigger)
	}{
		{
			name:             "complete trigger",
			organization:     "myorg",
			application:      "myapp",
			owner:            "john",
			parentPipelineId: "parent-pipeline-id",
			validate: func(t *testing.T, trigger *Trigger) {
				if trigger.Type != "pipeline" {
					t.Errorf("Expected Type 'pipeline', got %q", trigger.Type)
				}
				if !trigger.Enabled {
					t.Error("Expected Enabled to be true")
				}
				if trigger.Application != "myapp" {
					t.Errorf("Expected Application 'myapp', got %q", trigger.Application)
				}
				if trigger.Pipeline != "parent-pipeline-id" {
					t.Errorf("Expected Pipeline 'parent-pipeline-id', got %q", trigger.Pipeline)
				}
				expectedRunAsUser := "john-service-account@myorg.com"
				if trigger.RunAsUser != expectedRunAsUser {
					t.Errorf("Expected RunAsUser %q, got %q", expectedRunAsUser, trigger.RunAsUser)
				}
				if len(trigger.Status) != 1 {
					t.Fatalf("Expected 1 Status, got %d", len(trigger.Status))
				}
				if trigger.Status[0] != "successful" {
					t.Errorf("Expected Status[0] 'successful', got %q", trigger.Status[0])
				}
				if len(trigger.ExpectedArtifactIds) != 0 {
					t.Errorf("Expected empty ExpectedArtifactIds, got %v", trigger.ExpectedArtifactIds)
				}
			},
		},
		{
			name:             "empty parent pipeline id",
			organization:     "myorg",
			application:      "myapp",
			owner:            "jane",
			parentPipelineId: "",
			validate: func(t *testing.T, trigger *Trigger) {
				if trigger.Pipeline != "" {
					t.Errorf("Expected Pipeline '', got %q", trigger.Pipeline)
				}
			},
		},
		{
			name:             "empty application",
			organization:     "myorg",
			application:      "",
			owner:            "john",
			parentPipelineId: "parent-pipeline-id",
			validate: func(t *testing.T, trigger *Trigger) {
				if trigger.Application != "" {
					t.Errorf("Expected Application '', got %q", trigger.Application)
				}
			},
		},
		{
			name:             "empty owner",
			organization:     "myorg",
			application:      "myapp",
			owner:            "",
			parentPipelineId: "parent-pipeline-id",
			validate: func(t *testing.T, trigger *Trigger) {
				expectedRunAsUser := "-service-account@myorg.com"
				if trigger.RunAsUser != expectedRunAsUser {
					t.Errorf("Expected RunAsUser %q, got %q", expectedRunAsUser, trigger.RunAsUser)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := newPipelineTrigger(tt.organization, tt.application, tt.owner, tt.parentPipelineId)
			if result == nil {
				t.Fatal("newPipelineTrigger returned nil")
			}
			tt.validate(t, result)
		})
	}
}
