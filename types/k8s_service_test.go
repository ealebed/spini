package types

import (
	"testing"

	apiv1 "k8s.io/api/core/v1"
)

func TestNewService(t *testing.T) {
	tests := []struct {
		name        string
		application string
		stage       string
		namespace   string
		ports       []Port
		validate    func(*testing.T, *apiv1.Service)
	}{
		{
			name:        "production stage",
			application: "myapp",
			stage:       "production",
			namespace:   "default",
			ports: []Port{
				{Name: "http", ContainerPort: 8080},
				{Name: "metrics", ContainerPort: 9090},
			},
			validate: func(t *testing.T, svc *apiv1.Service) {
				if svc.Name != "myapp" {
					t.Errorf("Expected Name 'myapp', got %q", svc.Name)
				}
				if svc.Namespace != "default" {
					t.Errorf("Expected Namespace 'default', got %q", svc.Namespace)
				}
				if len(svc.Spec.Ports) != 2 {
					t.Fatalf("Expected 2 ports, got %d", len(svc.Spec.Ports))
				}
				if svc.Spec.ClusterIP != "None" {
					t.Errorf("Expected ClusterIP 'None', got %q", svc.Spec.ClusterIP)
				}
				if svc.Spec.Selector["app"] != "myapp" {
					t.Errorf("Expected Selector['app'] 'myapp', got %q", svc.Spec.Selector["app"])
				}
			},
		},
		{
			name:        "non-production stage",
			application: "myapp",
			stage:       "staging",
			namespace:   "staging",
			ports: []Port{
				{Name: "http", ContainerPort: 8080},
			},
			validate: func(t *testing.T, svc *apiv1.Service) {
				if svc.Name != "myapp-staging" {
					t.Errorf("Expected Name 'myapp-staging', got %q", svc.Name)
				}
				if svc.Spec.Selector["app"] != "myapp-staging" {
					t.Errorf("Expected Selector['app'] 'myapp-staging', got %q", svc.Spec.Selector["app"])
				}
			},
		},
		{
			name:        "empty ports",
			application: "myapp",
			stage:       "production",
			namespace:   "default",
			ports:       []Port{},
			validate: func(t *testing.T, svc *apiv1.Service) {
				if len(svc.Spec.Ports) != 0 {
					t.Errorf("Expected 0 ports, got %d", len(svc.Spec.Ports))
				}
			},
		},
		{
			name:        "nil ports",
			application: "myapp",
			stage:       "production",
			namespace:   "default",
			ports:       nil,
			validate: func(t *testing.T, svc *apiv1.Service) {
				if len(svc.Spec.Ports) != 0 {
					t.Errorf("Expected 0 ports, got %d", len(svc.Spec.Ports))
				}
			},
		},
		{
			name:        "empty application name",
			application: "",
			stage:       "production",
			namespace:   "default",
			ports:       []Port{{Name: "http", ContainerPort: 8080}},
			validate: func(t *testing.T, svc *apiv1.Service) {
				if svc.Name != "" {
					t.Errorf("Expected Name '', got %q", svc.Name)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewService(tt.application, tt.stage, tt.namespace, tt.ports)
			if result == nil {
				t.Fatal("NewService returned nil")
			}
			tt.validate(t, result)
		})
	}
}
