package types

import (
	"testing"
)

func TestDependencyContains(t *testing.T) {
	tests := []struct {
		name         string
		dependencies []DependsOn
		str          string
		expected     bool
	}{
		{
			name: "dependency found",
			dependencies: []DependsOn{
				{Name: "redis"},
				{Name: "postgres"},
				{Name: "elasticsearch"},
			},
			str:      "postgres",
			expected: true,
		},
		{
			name: "dependency not found",
			dependencies: []DependsOn{
				{Name: "redis"},
				{Name: "postgres"},
			},
			str:      "mongodb",
			expected: false,
		},
		{
			name:         "empty dependencies",
			dependencies: []DependsOn{},
			str:          "redis",
			expected:     false,
		},
		{
			name:         "nil dependencies",
			dependencies: nil,
			str:          "redis",
			expected:     false,
		},
		{
			name: "empty string search",
			dependencies: []DependsOn{
				{Name: "redis"},
				{Name: ""},
			},
			str:      "",
			expected: true,
		},
		{
			name: "case sensitive match",
			dependencies: []DependsOn{
				{Name: "Redis"},
				{Name: "Postgres"},
			},
			str:      "redis",
			expected: false,
		},
		{
			name: "single dependency match",
			dependencies: []DependsOn{
				{Name: "redis"},
			},
			str:      "redis",
			expected: true,
		},
		{
			name: "dependency with URL",
			dependencies: []DependsOn{
				{Name: "redis", URL: "http://redis.example.com"},
				{Name: "postgres"},
			},
			str:      "redis",
			expected: true,
		},
		{
			name: "partial name match fails",
			dependencies: []DependsOn{
				{Name: "redis-cluster"},
			},
			str:      "redis",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := dependencyContains(tt.dependencies, tt.str)
			if result != tt.expected {
				t.Errorf("dependencyContains(%v, %q) = %v, want %v", tt.dependencies, tt.str, result, tt.expected)
			}
		})
	}
}
