package bracket

import (
	"fmt"
	"strings"
)

// BracketError represents an error during bracket generation.
type BracketError struct {
	Message string
}

func (e *BracketError) Error() string {
	return e.Message
}

// ValidateParticipants checks that participant IDs are valid:
// at least minParticipants entries (default 2), no empty strings, no duplicates.
func ValidateParticipants(ids []string, minOpt ...int) error {
	min := 2
	if len(minOpt) > 0 && minOpt[0] > 0 {
		min = minOpt[0]
	}
	if len(ids) < min {
		return &BracketError{
			Message: fmt.Sprintf("At least %d participants required, got %d", min, len(ids)),
		}
	}

	for i, id := range ids {
		if strings.TrimSpace(id) == "" {
			return &BracketError{
				Message: fmt.Sprintf("Invalid participant ID at index %d: must be a non-empty string", i),
			}
		}
	}

	seen := make(map[string]struct{}, len(ids))
	for i, id := range ids {
		if _, exists := seen[id]; exists {
			return &BracketError{
				Message: fmt.Sprintf("Duplicate participant ID %q at index %d", id, i),
			}
		}
		seen[id] = struct{}{}
	}

	return nil
}
