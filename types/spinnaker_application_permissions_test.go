package types

import (
	"testing"
)

func TestDefaultApplicationPermissions(t *testing.T) {
	result := defaultApplicationPermissions()
	if result == nil {
		t.Fatal("defaultApplicationPermissions returned nil")
	}

	if len(result.Execute) != 1 || result.Execute[0] != "devops" {
		t.Errorf("Expected Execute ['devops'], got %v", result.Execute)
	}
	if len(result.Read) != 1 || result.Read[0] != "devops" {
		t.Errorf("Expected Read ['devops'], got %v", result.Read)
	}
	if len(result.Write) != 1 || result.Write[0] != "devops" {
		t.Errorf("Expected Write ['devops'], got %v", result.Write)
	}
}

func TestAppendApplicationPermissions(t *testing.T) {
	tests := []struct {
		name      string
		teamOwner string
		validate  func(*testing.T, *ApplicationPermissions)
	}{
		{
			name:      "standard team owner",
			teamOwner: "backend-team",
			validate: func(t *testing.T, perms *ApplicationPermissions) {
				if len(perms.Execute) != 2 {
					t.Fatalf("Expected 2 Execute permissions, got %d", len(perms.Execute))
				}
				if perms.Execute[0] != "devops" {
					t.Errorf("Expected Execute[0] 'devops', got %q", perms.Execute[0])
				}
				if perms.Execute[1] != "backend-team" {
					t.Errorf("Expected Execute[1] 'backend-team', got %q", perms.Execute[1])
				}
				if len(perms.Read) != 2 || perms.Read[1] != "backend-team" {
					t.Errorf("Expected Read[1] 'backend-team', got %q", perms.Read[1])
				}
				if len(perms.Write) != 2 || perms.Write[1] != "backend-team" {
					t.Errorf("Expected Write[1] 'backend-team', got %q", perms.Write[1])
				}
			},
		},
		{
			name:      "empty team owner",
			teamOwner: "",
			validate: func(t *testing.T, perms *ApplicationPermissions) {
				if len(perms.Execute) != 2 {
					t.Fatalf("Expected 2 Execute permissions, got %d", len(perms.Execute))
				}
				if perms.Execute[1] != "" {
					t.Errorf("Expected Execute[1] '', got %q", perms.Execute[1])
				}
			},
		},
		{
			name:      "special characters in team owner",
			teamOwner: "team-name_123",
			validate: func(t *testing.T, perms *ApplicationPermissions) {
				if perms.Execute[1] != "team-name_123" {
					t.Errorf("Expected Execute[1] 'team-name_123', got %q", perms.Execute[1])
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := AppendApplicationPermissions(tt.teamOwner)
			if result == nil {
				t.Fatal("AppendApplicationPermissions returned nil")
			}
			tt.validate(t, result)
		})
	}
}
