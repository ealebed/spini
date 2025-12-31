package types

import (
	"testing"
)

func TestDefaultExecutionWindowWhitelist(t *testing.T) {
	result := defaultExecutionWindowWhitelist()
	if result == nil {
		t.Fatal("defaultExecutionWindowWhitelist returned nil")
	}

	if result.StartHour != 7 {
		t.Errorf("Expected StartHour 7, got %d", result.StartHour)
	}
	if result.StartMin != 0 {
		t.Errorf("Expected StartMin 0, got %d", result.StartMin)
	}
	if result.EndHour != 11 {
		t.Errorf("Expected EndHour 11, got %d", result.EndHour)
	}
	if result.EndMin != 0 {
		t.Errorf("Expected EndMin 0, got %d", result.EndMin)
	}
}

func TestDefaultStageExecutionWindow(t *testing.T) {
	result := defaultStageExecutionWindow()
	if result == nil {
		t.Fatal("defaultStageExecutionWindow returned nil")
	}

	expectedDays := []int{2, 3, 4, 5, 6}
	if len(result.Days) != len(expectedDays) {
		t.Fatalf("Expected %d days, got %d", len(expectedDays), len(result.Days))
	}
	for i, day := range expectedDays {
		if result.Days[i] != day {
			t.Errorf("Expected Days[%d] %d, got %d", i, day, result.Days[i])
		}
	}

	if result.Whitelist == nil {
		t.Fatal("Expected Whitelist to be set")
	}
	if len(*result.Whitelist) != 1 {
		t.Fatalf("Expected 1 whitelist entry, got %d", len(*result.Whitelist))
	}
	if (*result.Whitelist)[0] == nil {
		t.Fatal("Expected Whitelist[0] to be set")
	}
	if (*result.Whitelist)[0].StartHour != 7 {
		t.Errorf("Expected Whitelist[0].StartHour 7, got %d", (*result.Whitelist)[0].StartHour)
	}
}
