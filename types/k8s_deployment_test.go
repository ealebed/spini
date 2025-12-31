package types

import (
	"testing"
)

func TestInt32Ptr(t *testing.T) {
	tests := []struct {
		name     string
		input    int32
		expected int32
	}{
		{
			name:     "zero value",
			input:    0,
			expected: 0,
		},
		{
			name:     "positive value",
			input:    42,
			expected: 42,
		},
		{
			name:     "negative value",
			input:    -10,
			expected: -10,
		},
		{
			name:     "max int32",
			input:    2147483647,
			expected: 2147483647,
		},
		{
			name:     "min int32",
			input:    -2147483648,
			expected: -2147483648,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := int32Ptr(tt.input)
			if result == nil {
				t.Fatal("int32Ptr returned nil")
			}
			if *result != tt.expected {
				t.Errorf("int32Ptr(%d) = %d, want %d", tt.input, *result, tt.expected)
			}
		})
	}
}

func TestInt64Ptr(t *testing.T) {
	tests := []struct {
		name     string
		input    int64
		expected int64
	}{
		{
			name:     "zero value",
			input:    0,
			expected: 0,
		},
		{
			name:     "positive value",
			input:    42,
			expected: 42,
		},
		{
			name:     "negative value",
			input:    -10,
			expected: -10,
		},
		{
			name:     "large positive value",
			input:    9223372036854775807,
			expected: 9223372036854775807,
		},
		{
			name:     "large negative value",
			input:    -9223372036854775808,
			expected: -9223372036854775808,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := int64Ptr(tt.input)
			if result == nil {
				t.Fatal("int64Ptr returned nil")
			}
			if *result != tt.expected {
				t.Errorf("int64Ptr(%d) = %d, want %d", tt.input, *result, tt.expected)
			}
		})
	}
}
