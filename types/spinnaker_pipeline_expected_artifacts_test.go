package types

import (
	"testing"
)

func TestNewDockerPipelineExpectedArtifact(t *testing.T) {
	tests := []struct {
		name         string
		organization string
		image        string
		version      string
		validate     func(*testing.T, *PipelineExpectedArtifact)
	}{
		{
			name:         "complete artifact",
			organization: "myorg",
			image:        "myapp",
			version:      "1.2.3",
			validate: func(t *testing.T, artifact *PipelineExpectedArtifact) {
				if artifact.DefaultArtifact == nil {
					t.Fatal("Expected DefaultArtifact to be set")
				}
				expectedName := dockerHubUrl + "myorg/myapp"
				if artifact.DefaultArtifact.Name != expectedName {
					t.Errorf("Expected DefaultArtifact.Name %q, got %q", expectedName, artifact.DefaultArtifact.Name)
				}
				expectedReference := dockerHubUrl + "myorg/myapp:1.2.3"
				if artifact.DefaultArtifact.Reference != expectedReference {
					t.Errorf("Expected DefaultArtifact.Reference %q, got %q", expectedReference, artifact.DefaultArtifact.Reference)
				}
				if artifact.DefaultArtifact.Type != "docker/image" {
					t.Errorf("Expected DefaultArtifact.Type 'docker/image', got %q", artifact.DefaultArtifact.Type)
				}
				if artifact.DefaultArtifact.Version != "1.2.3" {
					t.Errorf("Expected DefaultArtifact.Version '1.2.3', got %q", artifact.DefaultArtifact.Version)
				}
				if artifact.DisplayName != expectedName {
					t.Errorf("Expected DisplayName %q, got %q", expectedName, artifact.DisplayName)
				}
				if artifact.ID != "myorg/myapp" {
					t.Errorf("Expected ID 'myorg/myapp', got %q", artifact.ID)
				}
				if artifact.MatchArtifact == nil {
					t.Fatal("Expected MatchArtifact to be set")
				}
				if artifact.MatchArtifact.Name != expectedName {
					t.Errorf("Expected MatchArtifact.Name %q, got %q", expectedName, artifact.MatchArtifact.Name)
				}
				if !artifact.UseDefaultArtifact {
					t.Error("Expected UseDefaultArtifact to be true")
				}
				if !artifact.UsePriorArtifact {
					t.Error("Expected UsePriorArtifact to be true")
				}
			},
		},
		{
			name:         "empty version",
			organization: "myorg",
			image:        "myapp",
			version:      "",
			validate: func(t *testing.T, artifact *PipelineExpectedArtifact) {
				expectedReference := dockerHubUrl + "myorg/myapp:"
				if artifact.DefaultArtifact.Reference != expectedReference {
					t.Errorf("Expected DefaultArtifact.Reference %q, got %q", expectedReference, artifact.DefaultArtifact.Reference)
				}
				if artifact.DefaultArtifact.Version != "" {
					t.Errorf("Expected DefaultArtifact.Version '', got %q", artifact.DefaultArtifact.Version)
				}
			},
		},
		{
			name:         "empty organization",
			organization: "",
			image:        "myapp",
			version:      "1.0.0",
			validate: func(t *testing.T, artifact *PipelineExpectedArtifact) {
				expectedName := dockerHubUrl + "/myapp"
				if artifact.DefaultArtifact.Name != expectedName {
					t.Errorf("Expected DefaultArtifact.Name %q, got %q", expectedName, artifact.DefaultArtifact.Name)
				}
				if artifact.ID != "/myapp" {
					t.Errorf("Expected ID '/myapp', got %q", artifact.ID)
				}
			},
		},
		{
			name:         "empty image",
			organization: "myorg",
			image:        "",
			version:      "1.0.0",
			validate: func(t *testing.T, artifact *PipelineExpectedArtifact) {
				expectedName := dockerHubUrl + "myorg/"
				if artifact.DefaultArtifact.Name != expectedName {
					t.Errorf("Expected DefaultArtifact.Name %q, got %q", expectedName, artifact.DefaultArtifact.Name)
				}
				if artifact.ID != "myorg/" {
					t.Errorf("Expected ID 'myorg/', got %q", artifact.ID)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := newDockerPipelineExpectedArtifact(tt.organization, tt.image, tt.version)
			if result == nil {
				t.Fatal("newDockerPipelineExpectedArtifact returned nil")
			}
			tt.validate(t, result)
		})
	}
}

func TestNewManifestPipelineExpectedArtifact(t *testing.T) {
	tests := []struct {
		name             string
		githubContentUrl string
		relativePath     string
		validate         func(*testing.T, *PipelineExpectedArtifact)
	}{
		{
			name:             "complete artifact",
			githubContentUrl: "https://github.com/myorg/myrepo/blob/master/",
			relativePath:     "datacenters/gke1/default/myapp.yaml",
			validate: func(t *testing.T, artifact *PipelineExpectedArtifact) {
				if artifact.DefaultArtifact == nil {
					t.Fatal("Expected DefaultArtifact to be set")
				}
				if artifact.DefaultArtifact.Name != "datacenters/gke1/default/myapp.yaml" {
					t.Errorf("Expected DefaultArtifact.Name 'datacenters/gke1/default/myapp.yaml', got %q", artifact.DefaultArtifact.Name)
				}
				expectedReference := "https://github.com/myorg/myrepo/blob/master/datacenters/gke1/default/myapp.yaml"
				if artifact.DefaultArtifact.Reference != expectedReference {
					t.Errorf("Expected DefaultArtifact.Reference %q, got %q", expectedReference, artifact.DefaultArtifact.Reference)
				}
				if artifact.DefaultArtifact.Type != "github/file" {
					t.Errorf("Expected DefaultArtifact.Type 'github/file', got %q", artifact.DefaultArtifact.Type)
				}
				if artifact.DefaultArtifact.Version != "master" {
					t.Errorf("Expected DefaultArtifact.Version 'master', got %q", artifact.DefaultArtifact.Version)
				}
				if artifact.DisplayName != "datacenters/gke1/default/myapp.yaml" {
					t.Errorf("Expected DisplayName 'datacenters/gke1/default/myapp.yaml', got %q", artifact.DisplayName)
				}
				if artifact.ID != "datacenters/gke1/default/myapp.yaml" {
					t.Errorf("Expected ID 'datacenters/gke1/default/myapp.yaml', got %q", artifact.ID)
				}
				if artifact.MatchArtifact == nil {
					t.Fatal("Expected MatchArtifact to be set")
				}
				if artifact.MatchArtifact.Name != "datacenters/gke1/default/myapp.yaml" {
					t.Errorf("Expected MatchArtifact.Name 'datacenters/gke1/default/myapp.yaml', got %q", artifact.MatchArtifact.Name)
				}
				if !artifact.MatchArtifact.CustomKind {
					t.Error("Expected MatchArtifact.CustomKind to be true")
				}
				if !artifact.UseDefaultArtifact {
					t.Error("Expected UseDefaultArtifact to be true")
				}
				if artifact.UsePriorArtifact {
					t.Error("Expected UsePriorArtifact to be false")
				}
			},
		},
		{
			name:             "empty relative path",
			githubContentUrl: "https://github.com/myorg/myrepo/blob/master/",
			relativePath:     "",
			validate: func(t *testing.T, artifact *PipelineExpectedArtifact) {
				if artifact.DefaultArtifact.Name != "" {
					t.Errorf("Expected DefaultArtifact.Name '', got %q", artifact.DefaultArtifact.Name)
				}
				expectedReference := "https://github.com/myorg/myrepo/blob/master/"
				if artifact.DefaultArtifact.Reference != expectedReference {
					t.Errorf("Expected DefaultArtifact.Reference %q, got %q", expectedReference, artifact.DefaultArtifact.Reference)
				}
			},
		},
		{
			name:             "empty github content url",
			githubContentUrl: "",
			relativePath:     "datacenters/gke1/default/myapp.yaml",
			validate: func(t *testing.T, artifact *PipelineExpectedArtifact) {
				if artifact.DefaultArtifact.Reference != "datacenters/gke1/default/myapp.yaml" {
					t.Errorf("Expected DefaultArtifact.Reference 'datacenters/gke1/default/myapp.yaml', got %q", artifact.DefaultArtifact.Reference)
				}
			},
		},
		{
			name:             "root path",
			githubContentUrl: "https://github.com/myorg/myrepo/blob/master/",
			relativePath:     "/",
			validate: func(t *testing.T, artifact *PipelineExpectedArtifact) {
				if artifact.ID != "/" {
					t.Errorf("Expected ID '/', got %q", artifact.ID)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := newManifestPipelineExpectedArtifact(tt.githubContentUrl, tt.relativePath)
			if result == nil {
				t.Fatal("newManifestPipelineExpectedArtifact returned nil")
			}
			tt.validate(t, result)
		})
	}
}
