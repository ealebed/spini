package github

import (
	"net/url"
	"testing"
)

func TestAddOptions(t *testing.T) {
	tests := []struct {
		name     string
		baseURL  string
		opt      interface{}
		wantErr  bool
		validate func(*testing.T, string, error)
	}{
		{
			name:    "nil pointer options",
			baseURL: "https://api.example.com/v1",
			opt:     (*struct{ Field string })(nil),
			wantErr: false,
			validate: func(t *testing.T, result string, err error) {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				if result != "https://api.example.com/v1" {
					t.Errorf("Expected URL unchanged, got %q", result)
				}
			},
		},
		{
			name:    "valid URL with options",
			baseURL: "https://api.example.com/v1",
			opt: struct {
				Page    int    `url:"page"`
				PerPage int    `url:"per_page"`
				Filter  string `url:"filter"`
			}{
				Page:    1,
				PerPage: 10,
				Filter:  "active",
			},
			wantErr: false,
			validate: func(t *testing.T, result string, err error) {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				parsed, parseErr := url.Parse(result)
				if parseErr != nil {
					t.Fatalf("Failed to parse result URL: %v", parseErr)
				}
				if parsed.Query().Get("page") != "1" {
					t.Errorf("Expected page=1, got %q", parsed.Query().Get("page"))
				}
				if parsed.Query().Get("per_page") != "10" {
					t.Errorf("Expected per_page=10, got %q", parsed.Query().Get("per_page"))
				}
				if parsed.Query().Get("filter") != "active" {
					t.Errorf("Expected filter=active, got %q", parsed.Query().Get("filter"))
				}
			},
		},
		{
			name:    "URL with existing query parameters",
			baseURL: "https://api.example.com/v1?existing=value",
			opt: struct {
				NewParam string `url:"new_param"`
			}{
				NewParam: "new_value",
			},
			wantErr: false,
			validate: func(t *testing.T, result string, err error) {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				parsed, parseErr := url.Parse(result)
				if parseErr != nil {
					t.Fatalf("Failed to parse result URL: %v", parseErr)
				}
				// Function replaces existing query parameters, doesn't merge them
				if parsed.Query().Get("existing") != "" {
					t.Errorf("Expected existing to be replaced, got %q", parsed.Query().Get("existing"))
				}
				if parsed.Query().Get("new_param") != "new_value" {
					t.Errorf("Expected new_param=new_value, got %q", parsed.Query().Get("new_param"))
				}
			},
		},
		{
			name:    "empty options struct",
			baseURL: "https://api.example.com/v1",
			opt:     struct{}{},
			wantErr: false,
			validate: func(t *testing.T, result string, err error) {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				parsed, parseErr := url.Parse(result)
				if parseErr != nil {
					t.Fatalf("Failed to parse result URL: %v", parseErr)
				}
				if len(parsed.Query()) != 0 {
					t.Errorf("Expected no query parameters, got %v", parsed.Query())
				}
			},
		},
		{
			name:    "invalid URL",
			baseURL: "://invalid-url",
			opt: struct {
				Param string `url:"param"`
			}{
				Param: "value",
			},
			wantErr: true, // addOptions returns error on parse failure
			validate: func(t *testing.T, result string, err error) {
				if err == nil {
					t.Error("Expected error for invalid URL, got nil")
				}
				if result != "://invalid-url" {
					t.Errorf("Expected original URL on error, got %q", result)
				}
			},
		},
		{
			name:    "zero values in options",
			baseURL: "https://api.example.com/v1",
			opt: struct {
				Page    int    `url:"page"`
				PerPage int    `url:"per_page"`
				Filter  string `url:"filter"`
			}{
				Page:    0,
				PerPage: 0,
				Filter:  "",
			},
			wantErr: false,
			validate: func(t *testing.T, result string, err error) {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				parsed, parseErr := url.Parse(result)
				if parseErr != nil {
					t.Fatalf("Failed to parse result URL: %v", parseErr)
				}
				// Zero values should still be included in query string
				if parsed.Query().Get("page") != "0" {
					t.Errorf("Expected page=0, got %q", parsed.Query().Get("page"))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := addOptions(tt.baseURL, tt.opt)
			if (err != nil) != tt.wantErr {
				t.Errorf("addOptions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			tt.validate(t, result, err)
		})
	}
}
