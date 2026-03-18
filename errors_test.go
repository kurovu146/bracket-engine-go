package bracket

import (
	"strings"
	"testing"
)

func TestValidateParticipants_LessThan2(t *testing.T) {
	tests := []struct {
		name string
		ids  []string
	}{
		{name: "empty slice", ids: []string{}},
		{name: "single participant", ids: []string{"p1"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateParticipants(tt.ids)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if !strings.Contains(err.Error(), "At least 2 participants required") {
				t.Errorf("unexpected error message: %s", err.Error())
			}
		})
	}
}

func TestValidateParticipants_EmptyStringID(t *testing.T) {
	tests := []struct {
		name string
		ids  []string
	}{
		{name: "empty string at index 0", ids: []string{"", "p2"}},
		{name: "empty string at index 1", ids: []string{"p1", ""}},
		{name: "whitespace only", ids: []string{"p1", "   "}},
		{name: "tab only", ids: []string{"\t", "p2"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateParticipants(tt.ids)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if !strings.Contains(err.Error(), "must be a non-empty string") {
				t.Errorf("unexpected error message: %s", err.Error())
			}
		})
	}
}

func TestValidateParticipants_DuplicateIDs(t *testing.T) {
	err := ValidateParticipants([]string{"p1", "p2", "p1"})
	if err == nil {
		t.Fatal("expected error for duplicate IDs, got nil")
	}
	if !strings.Contains(err.Error(), "Duplicate participant ID") {
		t.Errorf("unexpected error message: %s", err.Error())
	}
	if !strings.Contains(err.Error(), `"p1"`) {
		t.Errorf("error should mention the duplicate ID 'p1': %s", err.Error())
	}
}

func TestValidateParticipants_Valid(t *testing.T) {
	tests := []struct {
		name string
		ids  []string
	}{
		{name: "2 participants", ids: []string{"p1", "p2"}},
		{name: "5 participants", ids: []string{"a", "b", "c", "d", "e"}},
		{name: "uuid-like IDs", ids: []string{"abc-123", "def-456", "ghi-789"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateParticipants(tt.ids)
			if err != nil {
				t.Errorf("expected no error, got %s", err.Error())
			}
		})
	}
}

func TestValidateParticipants_CustomMinimum(t *testing.T) {
	// Custom minimum of 4
	err := ValidateParticipants([]string{"p1", "p2", "p3"}, 4)
	if err == nil {
		t.Fatal("expected error when count < custom min, got nil")
	}
	if !strings.Contains(err.Error(), "At least 4 participants required") {
		t.Errorf("unexpected error message: %s", err.Error())
	}

	// Passing with custom minimum
	err = ValidateParticipants([]string{"p1", "p2", "p3", "p4"}, 4)
	if err != nil {
		t.Errorf("expected no error, got %s", err.Error())
	}
}

func TestBracketError_ImplementsError(t *testing.T) {
	var err error = &BracketError{Message: "test error"}
	if err.Error() != "test error" {
		t.Errorf("expected 'test error', got %q", err.Error())
	}
}
