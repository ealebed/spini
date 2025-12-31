package types

import (
	"testing"

	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func TestGetIntOrDefault(t *testing.T) {
	tests := []struct {
		name         string
		value        int
		defaultValue int
		expected     int
	}{
		{
			name:         "positive value returns value",
			value:        8080,
			defaultValue: 80,
			expected:     8080,
		},
		{
			name:         "zero value returns default",
			value:        0,
			defaultValue: 80,
			expected:     80,
		},
		{
			name:         "negative value returns default",
			value:        -1,
			defaultValue: 80,
			expected:     80,
		},
		{
			name:         "value 1 returns value",
			value:        1,
			defaultValue: 80,
			expected:     1,
		},
		{
			name:         "large value returns value",
			value:        65535,
			defaultValue: 80,
			expected:     65535,
		},
		{
			name:         "zero default with zero value",
			value:        0,
			defaultValue: 0,
			expected:     0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getIntOrDefault(tt.value, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("getIntOrDefault(%d, %d) = %d, want %d", tt.value, tt.defaultValue, result, tt.expected)
			}
		})
	}
}

func TestNewLivenessProbe(t *testing.T) {
	tests := []struct {
		name     string
		tier     *Datacenter
		validate func(*testing.T, *apiv1.Probe)
	}{
		{
			name: "file type probe",
			tier: &Datacenter{
				LivenessProbe: &Probe{
					Type:             probeTypeFile,
					Delay:            10,
					Port:             8080,
					TimeoutSeconds:   5,
					PeriodSeconds:    10,
					SuccessThreshold: 1,
					FailureThreshold: 3,
				},
			},
			validate: func(t *testing.T, probe *apiv1.Probe) {
				if probe.ProbeHandler.Exec == nil {
					t.Error("Expected Exec probe handler for file type")
				}
				if probe.ProbeHandler.HTTPGet != nil {
					t.Error("Expected no HTTPGet probe handler for file type")
				}
				if len(probe.ProbeHandler.Exec.Command) != 2 || probe.ProbeHandler.Exec.Command[0] != "cat" {
					t.Errorf("Expected command ['cat', '/tmp/live'], got %v", probe.ProbeHandler.Exec.Command)
				}
				if probe.InitialDelaySeconds != 10 {
					t.Errorf("Expected InitialDelaySeconds 10, got %d", probe.InitialDelaySeconds)
				}
			},
		},
		{
			name: "http type probe with path",
			tier: &Datacenter{
				LivenessProbe: &Probe{
					Type:             "http",
					Path:             "/healthz",
					Delay:            15,
					Port:             9090,
					TimeoutSeconds:   10,
					PeriodSeconds:    20,
					SuccessThreshold: 1,
					FailureThreshold: 3,
				},
			},
			validate: func(t *testing.T, probe *apiv1.Probe) {
				if probe.ProbeHandler.HTTPGet == nil {
					t.Error("Expected HTTPGet probe handler for http type")
				}
				if probe.ProbeHandler.Exec != nil {
					t.Error("Expected no Exec probe handler for http type")
				}
				if probe.ProbeHandler.HTTPGet.Path != "/healthz" {
					t.Errorf("Expected path /healthz, got %s", probe.ProbeHandler.HTTPGet.Path)
				}
				if probe.ProbeHandler.HTTPGet.Port != intstr.FromInt(9090) {
					t.Errorf("Expected port 9090, got %v", probe.ProbeHandler.HTTPGet.Port)
				}
				if probe.InitialDelaySeconds != 15 {
					t.Errorf("Expected InitialDelaySeconds 15, got %d", probe.InitialDelaySeconds)
				}
			},
		},
		{
			name: "http type probe without path uses default",
			tier: &Datacenter{
				LivenessProbe: &Probe{
					Type:             "http",
					Path:             "",
					Delay:            20,
					Port:             0,
					TimeoutSeconds:   5,
					PeriodSeconds:    10,
					SuccessThreshold: 1,
					FailureThreshold: 3,
				},
			},
			validate: func(t *testing.T, probe *apiv1.Probe) {
				if probe.ProbeHandler.HTTPGet == nil {
					t.Error("Expected HTTPGet probe handler")
				}
				if probe.ProbeHandler.HTTPGet.Path != probePathHealth {
					t.Errorf("Expected default path %s, got %s", probePathHealth, probe.ProbeHandler.HTTPGet.Path)
				}
				if probe.ProbeHandler.HTTPGet.Port != intstr.FromInt(8080) {
					t.Errorf("Expected default port 8080, got %v", probe.ProbeHandler.HTTPGet.Port)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := newLivenessProbe(tt.tier)
			if result == nil {
				t.Fatal("newLivenessProbe returned nil")
			}
			tt.validate(t, result)
		})
	}
}

func TestNewReadinessProbe(t *testing.T) {
	tests := []struct {
		name     string
		tier     *Datacenter
		validate func(*testing.T, *apiv1.Probe)
	}{
		{
			name: "readiness probe uses liveness when nil",
			tier: &Datacenter{
				LivenessProbe: &Probe{
					Type:             probeTypeFile,
					Delay:            10,
					Port:             8080,
					TimeoutSeconds:   5,
					PeriodSeconds:    10,
					SuccessThreshold: 1,
					FailureThreshold: 3,
				},
				ReadinessProbe: nil,
			},
			validate: func(t *testing.T, probe *apiv1.Probe) {
				if probe.ProbeHandler.Exec == nil {
					t.Error("Expected Exec probe handler for file type")
				}
				if len(probe.ProbeHandler.Exec.Command) != 2 || probe.ProbeHandler.Exec.Command[1] != "/tmp/ready" {
					t.Errorf("Expected command ['cat', '/tmp/ready'], got %v", probe.ProbeHandler.Exec.Command)
				}
			},
		},
		{
			name: "readiness probe with custom config",
			tier: &Datacenter{
				LivenessProbe: &Probe{
					Type: "http",
					Path: "/live",
				},
				ReadinessProbe: &Probe{
					Type:             "http",
					Path:             "/ready",
					Delay:            5,
					Port:             8080,
					TimeoutSeconds:   3,
					PeriodSeconds:    5,
					SuccessThreshold: 1,
					FailureThreshold: 2,
				},
			},
			validate: func(t *testing.T, probe *apiv1.Probe) {
				if probe.ProbeHandler.HTTPGet == nil {
					t.Error("Expected HTTPGet probe handler")
				}
				if probe.ProbeHandler.HTTPGet.Path != "/ready" {
					t.Errorf("Expected path /ready, got %s", probe.ProbeHandler.HTTPGet.Path)
				}
				if probe.InitialDelaySeconds != 5 {
					t.Errorf("Expected InitialDelaySeconds 5, got %d", probe.InitialDelaySeconds)
				}
			},
		},
		{
			name: "readiness probe without path uses default",
			tier: &Datacenter{
				LivenessProbe: &Probe{
					Type: "http",
				},
				ReadinessProbe: &Probe{
					Type: "http",
					Path: "",
				},
			},
			validate: func(t *testing.T, probe *apiv1.Probe) {
				if probe.ProbeHandler.HTTPGet.Path != probePathHealth {
					t.Errorf("Expected default path %s, got %s", probePathHealth, probe.ProbeHandler.HTTPGet.Path)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := newReadinessProbe(tt.tier)
			if result == nil {
				t.Fatal("newReadinessProbe returned nil")
			}
			tt.validate(t, result)
		})
	}
}

func TestNewStartupProbe(t *testing.T) {
	tests := []struct {
		name     string
		tier     *Datacenter
		validate func(*testing.T, *apiv1.Probe)
	}{
		{
			name: "file type startup probe",
			tier: &Datacenter{
				StartupProbe: &Probe{
					Type:             probeTypeFile,
					Delay:            0,
					Port:             8080,
					TimeoutSeconds:   5,
					PeriodSeconds:    10,
					SuccessThreshold: 1,
					FailureThreshold: 30,
				},
			},
			validate: func(t *testing.T, probe *apiv1.Probe) {
				if probe.ProbeHandler.Exec == nil {
					t.Error("Expected Exec probe handler for file type")
				}
				if len(probe.ProbeHandler.Exec.Command) != 2 || probe.ProbeHandler.Exec.Command[1] != "/tmp/started" {
					t.Errorf("Expected command ['cat', '/tmp/started'], got %v", probe.ProbeHandler.Exec.Command)
				}
				if probe.FailureThreshold != 30 {
					t.Errorf("Expected FailureThreshold 30, got %d", probe.FailureThreshold)
				}
			},
		},
		{
			name: "http type startup probe with path",
			tier: &Datacenter{
				StartupProbe: &Probe{
					Type:             "http",
					Path:             "/startup",
					Delay:            0,
					Port:             9090,
					TimeoutSeconds:   5,
					PeriodSeconds:    10,
					SuccessThreshold: 1,
					FailureThreshold: 30,
				},
			},
			validate: func(t *testing.T, probe *apiv1.Probe) {
				if probe.ProbeHandler.HTTPGet == nil {
					t.Error("Expected HTTPGet probe handler for http type")
				}
				if probe.ProbeHandler.HTTPGet.Path != "/startup" {
					t.Errorf("Expected path /startup, got %s", probe.ProbeHandler.HTTPGet.Path)
				}
				if probe.ProbeHandler.HTTPGet.Port != intstr.FromInt(9090) {
					t.Errorf("Expected port 9090, got %v", probe.ProbeHandler.HTTPGet.Port)
				}
			},
		},
		{
			name: "http type startup probe without path uses default",
			tier: &Datacenter{
				StartupProbe: &Probe{
					Type: "http",
					Path: "",
					Port: 0,
				},
			},
			validate: func(t *testing.T, probe *apiv1.Probe) {
				if probe.ProbeHandler.HTTPGet.Path != probePathHealth {
					t.Errorf("Expected default path %s, got %s", probePathHealth, probe.ProbeHandler.HTTPGet.Path)
				}
				if probe.ProbeHandler.HTTPGet.Port != intstr.FromInt(8080) {
					t.Errorf("Expected default port 8080, got %v", probe.ProbeHandler.HTTPGet.Port)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := newStartupProbe(tt.tier)
			if result == nil {
				t.Fatal("newStartupProbe returned nil")
			}
			tt.validate(t, result)
		})
	}
}
