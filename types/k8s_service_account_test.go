package types

import (
	"testing"

	apiv1 "k8s.io/api/core/v1"
)

func TestNewServiceAccount(t *testing.T) {
	tests := []struct {
		name                string
		application         string
		stage               string
		namespace           string
		imagePullSecretName string
		validate            func(*testing.T, *apiv1.ServiceAccount)
	}{
		{
			name:                "production stage",
			application:         "myapp",
			stage:               "production",
			namespace:           "default",
			imagePullSecretName: "dockerhubkey",
			validate: func(t *testing.T, sa *apiv1.ServiceAccount) {
				if sa.Name != "myapp" {
					t.Errorf("Expected Name 'myapp', got %q", sa.Name)
				}
				if sa.Namespace != "default" {
					t.Errorf("Expected Namespace 'default', got %q", sa.Namespace)
				}
				if len(sa.ImagePullSecrets) != 1 {
					t.Fatalf("Expected 1 ImagePullSecret, got %d", len(sa.ImagePullSecrets))
				}
				if sa.ImagePullSecrets[0].Name != "dockerhubkey" {
					t.Errorf("Expected ImagePullSecrets[0].Name 'dockerhubkey', got %q", sa.ImagePullSecrets[0].Name)
				}
			},
		},
		{
			name:                "non-production stage",
			application:         "myapp",
			stage:               "staging",
			namespace:           "staging",
			imagePullSecretName: "dockerhubkey",
			validate: func(t *testing.T, sa *apiv1.ServiceAccount) {
				if sa.Name != "myapp-staging" {
					t.Errorf("Expected Name 'myapp-staging', got %q", sa.Name)
				}
				if sa.Namespace != "staging" {
					t.Errorf("Expected Namespace 'staging', got %q", sa.Namespace)
				}
			},
		},
		{
			name:                "empty image pull secret name",
			application:         "myapp",
			stage:               "production",
			namespace:           "default",
			imagePullSecretName: "",
			validate: func(t *testing.T, sa *apiv1.ServiceAccount) {
				if len(sa.ImagePullSecrets) != 1 {
					t.Fatalf("Expected 1 ImagePullSecret, got %d", len(sa.ImagePullSecrets))
				}
				if sa.ImagePullSecrets[0].Name != "" {
					t.Errorf("Expected ImagePullSecrets[0].Name '', got %q", sa.ImagePullSecrets[0].Name)
				}
			},
		},
		{
			name:                "empty application name",
			application:         "",
			stage:               "production",
			namespace:           "default",
			imagePullSecretName: "dockerhubkey",
			validate: func(t *testing.T, sa *apiv1.ServiceAccount) {
				if sa.Name != "" {
					t.Errorf("Expected Name '', got %q", sa.Name)
				}
			},
		},
		{
			name:                "empty namespace",
			application:         "myapp",
			stage:               "production",
			namespace:           "",
			imagePullSecretName: "dockerhubkey",
			validate: func(t *testing.T, sa *apiv1.ServiceAccount) {
				if sa.Namespace != "" {
					t.Errorf("Expected Namespace '', got %q", sa.Namespace)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewServiceAccount(tt.application, tt.stage, tt.namespace, tt.imagePullSecretName)
			if result == nil {
				t.Fatal("NewServiceAccount returned nil")
			}
			tt.validate(t, result)
		})
	}
}
