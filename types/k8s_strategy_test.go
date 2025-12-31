package types

import (
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func TestDefaultDeploymentStrategy(t *testing.T) {
	result := defaultDeploymentStrategy()

	if result.Type != appsv1.RollingUpdateDeploymentStrategyType {
		t.Errorf("Expected Type 'RollingUpdate', got %q", result.Type)
	}
	if result.RollingUpdate == nil {
		t.Fatal("Expected RollingUpdate to be set")
	}
	if result.RollingUpdate.MaxSurge == nil {
		t.Fatal("Expected MaxSurge to be set")
	}
	if result.RollingUpdate.MaxSurge.Type != intstr.String {
		t.Errorf("Expected MaxSurge.Type String, got %v", result.RollingUpdate.MaxSurge.Type)
	}
	if result.RollingUpdate.MaxSurge.StrVal != "25%" {
		t.Errorf("Expected MaxSurge.StrVal '25%%', got %q", result.RollingUpdate.MaxSurge.StrVal)
	}
	if result.RollingUpdate.MaxUnavailable == nil {
		t.Fatal("Expected MaxUnavailable to be set")
	}
	if result.RollingUpdate.MaxUnavailable.StrVal != "25%" {
		t.Errorf("Expected MaxUnavailable.StrVal '25%%', got %q", result.RollingUpdate.MaxUnavailable.StrVal)
	}
}

func TestNewDeploymentStrategy(t *testing.T) {
	tests := []struct {
		name     string
		strategy *DeployStrategy
		validate func(*testing.T, appsv1.DeploymentStrategy)
	}{
		{
			name:     "nil strategy uses default",
			strategy: nil,
			validate: func(t *testing.T, result appsv1.DeploymentStrategy) {
				if result.Type != appsv1.RollingUpdateDeploymentStrategyType {
					t.Errorf("Expected Type 'RollingUpdate', got %q", result.Type)
				}
				if result.RollingUpdate == nil {
					t.Fatal("Expected RollingUpdate to be set")
				}
			},
		},
		{
			name: "RollingUpdate strategy",
			strategy: &DeployStrategy{
				Type: "RollingUpdate",
				RollingUpdate: &RollingUpdateDeployment{
					MaxSurge:       "50%",
					MaxUnavailable: "30%",
				},
			},
			validate: func(t *testing.T, result appsv1.DeploymentStrategy) {
				if result.Type != appsv1.RollingUpdateDeploymentStrategyType {
					t.Errorf("Expected Type 'RollingUpdate', got %q", result.Type)
				}
				if result.RollingUpdate.MaxSurge.StrVal != "50%" {
					t.Errorf("Expected MaxSurge '50%%', got %q", result.RollingUpdate.MaxSurge.StrVal)
				}
				if result.RollingUpdate.MaxUnavailable.StrVal != "30%" {
					t.Errorf("Expected MaxUnavailable '30%%', got %q", result.RollingUpdate.MaxUnavailable.StrVal)
				}
			},
		},
		{
			name: "Recreate strategy",
			strategy: &DeployStrategy{
				Type: "Recreate",
			},
			validate: func(t *testing.T, result appsv1.DeploymentStrategy) {
				if result.Type != appsv1.RecreateDeploymentStrategyType {
					t.Errorf("Expected Type 'Recreate', got %q", result.Type)
				}
				if result.RollingUpdate != nil {
					t.Error("Expected RollingUpdate to be nil for Recreate strategy")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewDeploymentStrategy(tt.strategy)
			tt.validate(t, result)
		})
	}
}
