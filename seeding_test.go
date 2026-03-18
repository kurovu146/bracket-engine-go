package bracket

import (
	"reflect"
	"testing"
)

func TestGenerateSeedOrder(t *testing.T) {
	tests := []struct {
		name        string
		bracketSize int
		want        []int
	}{
		{
			name:        "size 2",
			bracketSize: 2,
			want:        []int{1, 2},
		},
		{
			name:        "size 4",
			bracketSize: 4,
			want:        []int{1, 4, 2, 3},
		},
		{
			name:        "size 8",
			bracketSize: 8,
			want:        []int{1, 8, 4, 5, 2, 7, 3, 6},
		},
		{
			name:        "size 16",
			bracketSize: 16,
			want:        []int{1, 16, 8, 9, 4, 13, 5, 12, 2, 15, 7, 10, 3, 14, 6, 11},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateSeedOrder(tt.bracketSize)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenerateSeedOrder(%d) = %v, want %v", tt.bracketSize, got, tt.want)
			}
		})
	}
}

func TestStandardSeed_ExactPowerOf2(t *testing.T) {
	// 4 players in bracket size 4 — no byes
	ids := []string{"p1", "p2", "p3", "p4"}
	result := StandardSeed(ids, 4)

	if len(result) != 4 {
		t.Fatalf("expected 4 slots, got %d", len(result))
	}

	// Seed order for size 4: [1, 4, 2, 3]
	// Position 0 → seed 1 → p1
	// Position 1 → seed 4 → p4
	// Position 2 → seed 2 → p2
	// Position 3 → seed 3 → p3
	expected := []string{"p1", "p4", "p2", "p3"}
	for i, exp := range expected {
		if result[i] == nil {
			t.Errorf("position %d: expected %q, got nil (bye)", i, exp)
		} else if *result[i] != exp {
			t.Errorf("position %d: expected %q, got %q", i, exp, *result[i])
		}
	}

	// No byes
	for i, r := range result {
		if r == nil {
			t.Errorf("position %d: unexpected bye", i)
		}
	}
}

func TestStandardSeed_WithByes(t *testing.T) {
	// 5 players in bracket size 8 — 3 byes
	ids := []string{"p1", "p2", "p3", "p4", "p5"}
	result := StandardSeed(ids, 8)

	if len(result) != 8 {
		t.Fatalf("expected 8 slots, got %d", len(result))
	}

	// Seed order for size 8: [1, 8, 4, 5, 2, 7, 3, 6]
	// Seeds 6, 7, 8 are beyond 5 players → byes
	byeCount := 0
	playerCount := 0
	for _, r := range result {
		if r == nil {
			byeCount++
		} else {
			playerCount++
		}
	}

	if byeCount != 3 {
		t.Errorf("expected 3 byes, got %d", byeCount)
	}
	if playerCount != 5 {
		t.Errorf("expected 5 players, got %d", playerCount)
	}

	// Verify specific positions:
	// Position 0 → seed 1 → p1
	// Position 1 → seed 8 → nil (bye)
	// Position 2 → seed 4 → p4
	// Position 3 → seed 5 → p5
	// Position 4 → seed 2 → p2
	// Position 5 → seed 7 → nil (bye)
	// Position 6 → seed 3 → p3
	// Position 7 → seed 6 → nil (bye)
	expectNonNil := map[int]string{0: "p1", 2: "p4", 3: "p5", 4: "p2", 6: "p3"}
	expectNil := []int{1, 5, 7}

	for pos, expectedID := range expectNonNil {
		if result[pos] == nil {
			t.Errorf("position %d: expected %q, got nil", pos, expectedID)
		} else if *result[pos] != expectedID {
			t.Errorf("position %d: expected %q, got %q", pos, expectedID, *result[pos])
		}
	}

	for _, pos := range expectNil {
		if result[pos] != nil {
			t.Errorf("position %d: expected nil (bye), got %q", pos, *result[pos])
		}
	}
}

func TestStandardSeed_2Players(t *testing.T) {
	ids := []string{"alice", "bob"}
	result := StandardSeed(ids, 2)

	if len(result) != 2 {
		t.Fatalf("expected 2 slots, got %d", len(result))
	}

	// Seed order for size 2: [1, 2]
	if result[0] == nil || *result[0] != "alice" {
		t.Errorf("position 0: expected 'alice', got %v", result[0])
	}
	if result[1] == nil || *result[1] != "bob" {
		t.Errorf("position 1: expected 'bob', got %v", result[1])
	}
}
