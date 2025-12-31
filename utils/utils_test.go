package utils

import (
	"testing"
)

func TestSliceContains(t *testing.T) {
	tests := []struct {
		name     string
		slice    []string
		str      string
		expected bool
	}{
		{
			name:     "string found in slice",
			slice:    []string{"alpha", "beta", "gamma"},
			str:      "beta",
			expected: true,
		},
		{
			name:     "string not found in slice",
			slice:    []string{"alpha", "beta", "gamma"},
			str:      "delta",
			expected: false,
		},
		{
			name:     "empty slice",
			slice:    []string{},
			str:      "alpha",
			expected: false,
		},
		{
			name:     "nil slice",
			slice:    nil,
			str:      "alpha",
			expected: false,
		},
		{
			name:     "empty string in slice",
			slice:    []string{"alpha", "", "gamma"},
			str:      "",
			expected: true,
		},
		{
			name:     "case sensitive match",
			slice:    []string{"Alpha", "Beta", "Gamma"},
			str:      "alpha",
			expected: false,
		},
		{
			name:     "single element slice match",
			slice:    []string{"alpha"},
			str:      "alpha",
			expected: true,
		},
		{
			name:     "single element slice no match",
			slice:    []string{"alpha"},
			str:      "beta",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sliceContains(tt.slice, tt.str)
			if result != tt.expected {
				t.Errorf("sliceContains(%v, %q) = %v, want %v", tt.slice, tt.str, result, tt.expected)
			}
		})
	}
}

func TestFillPipelineConfig(t *testing.T) {
	tests := []struct {
		name        string
		stage       string
		pList       []string
		pipelineIDs map[string]string
		expected    map[string]interface{}
	}{
		{
			name:  "production stage with beta and nightly",
			stage: stageProduction,
			pList: []string{"beta", stageNightly},
			pipelineIDs: map[string]string{
				"promote-to-production": "promote-prod-id",
				"beta-gke1":             "beta-id",
				"nightly-gke1":          "nightly-id",
			},
			expected: map[string]interface{}{
				"stage":                   stageProduction,
				"dockerTriggerEnabled":    true,
				"id":                      "promote-prod-id",
				"parentPipelineId":        "nightly-id",
				"GeneratePromotePipeline": true,
				"pipelineTriggerEnabled":  true,
			},
		},
		{
			name:  "nightly stage with beta",
			stage: stageNightly,
			pList: []string{"beta", stageNightly},
			pipelineIDs: map[string]string{
				"promote-to-nightly": "promote-nightly-id",
				"beta-gke1":          "beta-id",
			},
			expected: map[string]interface{}{
				"stage":                   stageNightly,
				"dockerTriggerEnabled":    true, // Line 137 always sets this to true
				"id":                      "promote-nightly-id",
				"parentPipelineId":        "beta-id",
				"GeneratePromotePipeline": true,
				"pipelineTriggerEnabled":  true,
			},
		},
		{
			name:        "beta stage without nightly",
			stage:       "beta",
			pList:       []string{"beta"},
			pipelineIDs: map[string]string{},
			expected: map[string]interface{}{
				"stage":                   "beta",
				"dockerTriggerEnabled":    true,
				"GeneratePromotePipeline": false,
				"pipelineTriggerEnabled":  false,
			},
		},
		{
			name:        "dev stage with beta and nightly",
			stage:       "dev",
			pList:       []string{"beta", stageNightly},
			pipelineIDs: map[string]string{},
			expected: map[string]interface{}{
				"stage":                   "dev",
				"dockerTriggerEnabled":    true, // Line 137 always sets this to true
				"GeneratePromotePipeline": false,
				"pipelineTriggerEnabled":  false,
			},
		},
		{
			name:        "nightly stage without beta",
			stage:       stageNightly,
			pList:       []string{stageNightly},
			pipelineIDs: map[string]string{},
			expected: map[string]interface{}{
				"stage":                   stageNightly,
				"dockerTriggerEnabled":    true,
				"GeneratePromotePipeline": false,
				"pipelineTriggerEnabled":  false,
			},
		},
		{
			name:  "production stage with only nightly",
			stage: stageProduction,
			pList: []string{stageNightly},
			pipelineIDs: map[string]string{
				"promote-to-production": "promote-prod-id",
				"nightly-gke1":          "nightly-id",
			},
			expected: map[string]interface{}{
				"stage":                   stageProduction,
				"dockerTriggerEnabled":    true,
				"id":                      "promote-prod-id",
				"parentPipelineId":        "nightly-id",
				"GeneratePromotePipeline": true,
				"pipelineTriggerEnabled":  true,
			},
		},
		{
			name:        "empty profile list",
			stage:       "staging",
			pList:       []string{},
			pipelineIDs: map[string]string{},
			expected: map[string]interface{}{
				"stage":                   "staging",
				"dockerTriggerEnabled":    true,
				"GeneratePromotePipeline": false,
				"pipelineTriggerEnabled":  false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := fillPipelineConfig(tt.stage, tt.pList, tt.pipelineIDs)

			// Check all expected keys
			for key, expectedValue := range tt.expected {
				if result[key] != expectedValue {
					t.Errorf("fillPipelineConfig(%q, %v, %v)[%q] = %v, want %v",
						tt.stage, tt.pList, tt.pipelineIDs, key, result[key], expectedValue)
				}
			}

			// Ensure result has stage
			if result["stage"] != tt.stage {
				t.Errorf("fillPipelineConfig(%q, %v, %v)[\"stage\"] = %v, want %v",
					tt.stage, tt.pList, tt.pipelineIDs, result["stage"], tt.stage)
			}
		})
	}
}
