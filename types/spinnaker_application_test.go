package types

import (
	"testing"
)

func TestDefaultCustomBanner(t *testing.T) {
	result := defaultCustomBanner()
	if result == nil {
		t.Fatal("defaultCustomBanner returned nil")
	}

	if result.BackgroundColor != "var(--color-accessory-light)" {
		t.Errorf("Expected BackgroundColor 'var(--color-accessory-light)', got %q", result.BackgroundColor)
	}
	if !result.Enabled {
		t.Error("Expected Enabled to be true")
	}
	if result.Text != "Default Custom Banner Text" {
		t.Errorf("Expected Text 'Default Custom Banner Text', got %q", result.Text)
	}
	if result.TextColor != "var(--color-text-primary)" {
		t.Errorf("Expected TextColor 'var(--color-text-primary)', got %q", result.TextColor)
	}
}

func TestDefaultApplication(t *testing.T) {
	result := defaultApplication()
	if result == nil {
		t.Fatal("defaultApplication returned nil")
	}

	if result.CloudProviders != "kubernetes" {
		t.Errorf("Expected CloudProviders 'kubernetes', got %q", result.CloudProviders)
	}
	if len(result.CustomBanners) != 1 {
		t.Fatalf("Expected 1 CustomBanner, got %d", len(result.CustomBanners))
	}
	if result.CustomBanners[0] == nil {
		t.Fatal("Expected CustomBanners[0] to be set")
	}
	if result.DataSources == nil {
		t.Fatal("Expected DataSources to be set")
	}
	if len(result.DataSources.Disabled) != 1 || result.DataSources.Disabled[0] != "securityGroups" {
		t.Errorf("Expected DataSources.Disabled ['securityGroups'], got %v", result.DataSources.Disabled)
	}
	if len(result.DataSources.Enabled) != 0 {
		t.Errorf("Expected empty DataSources.Enabled, got %v", result.DataSources.Enabled)
	}
	if len(result.TrafficGuards) != 0 {
		t.Errorf("Expected empty TrafficGuards, got %v", result.TrafficGuards)
	}
}
