package types

import (
	"testing"
)

func TestDefaultPipelineTrafficManagement(t *testing.T) {
	result := defaultPipelineTrafficManagement()
	if result == nil {
		t.Fatal("defaultPipelineTrafficManagement returned nil")
	}

	if result.Enabled {
		t.Error("Expected Enabled to be false")
	}
	if result.Options == nil {
		t.Fatal("Expected Options to be set")
	}
	if result.Options.EnableTraffic {
		t.Error("Expected Options.EnableTraffic to be false")
	}
	if len(result.Options.Services) != 0 {
		t.Errorf("Expected empty Services, got %v", result.Options.Services)
	}
}

func TestDefaultPipelineTrafficManagementOptions(t *testing.T) {
	result := defaultPipelineTrafficManagementOptions()
	if result == nil {
		t.Fatal("defaultPipelineTrafficManagementOptions returned nil")
	}

	if result.EnableTraffic {
		t.Error("Expected EnableTraffic to be false")
	}
	if len(result.Services) != 0 {
		t.Errorf("Expected empty Services, got %v", result.Services)
	}
}
