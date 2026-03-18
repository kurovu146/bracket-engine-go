package bracket

import (
	"fmt"
	"strings"
	"testing"
)

func TestGroupStage_8Players_2Groups(t *testing.T) {
	ids := makePlayerIDsUnsorted(8)
	numGroups := 2
	opts := &GroupStageOptions{NumGroups: &numGroups}
	result, err := GenerateGroupStage(ids, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 2 groups of 4 players each
	if len(result.Groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(result.Groups))
	}
	if len(result.Groups[0]) != 4 {
		t.Errorf("group 0: expected 4 players, got %d", len(result.Groups[0]))
	}
	if len(result.Groups[1]) != 4 {
		t.Errorf("group 1: expected 4 players, got %d", len(result.Groups[1]))
	}

	// Each group of 4 → C(4,2) = 6 matches, total = 12
	if len(result.Matches) != 12 {
		t.Fatalf("expected 12 matches, got %d", len(result.Matches))
	}

	// Verify bracket types
	g0Count := countMatchesByBracketType(result.Matches, GroupBracketType(0))
	g1Count := countMatchesByBracketType(result.Matches, GroupBracketType(1))
	if g0Count != 6 {
		t.Errorf("group 0: expected 6 matches, got %d", g0Count)
	}
	if g1Count != 6 {
		t.Errorf("group 1: expected 6 matches, got %d", g1Count)
	}
}

func TestGroupStage_SequentialDistribution(t *testing.T) {
	ids := makePlayerIDsUnsorted(8)
	numGroups := 2
	opts := &GroupStageOptions{
		NumGroups:    &numGroups,
		Distribution: "sequential",
	}
	result, err := GenerateGroupStage(ids, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Sequential: p1→G0, p2→G1, p3→G0, p4→G1, p5→G0, p6→G1, p7→G0, p8→G1
	expectedG0 := []string{"p1", "p3", "p5", "p7"}
	expectedG1 := []string{"p2", "p4", "p6", "p8"}

	for i, exp := range expectedG0 {
		if i >= len(result.Groups[0]) || result.Groups[0][i] != exp {
			t.Errorf("G0[%d]: expected %q, got %q", i, exp, result.Groups[0][i])
		}
	}
	for i, exp := range expectedG1 {
		if i >= len(result.Groups[1]) || result.Groups[1][i] != exp {
			t.Errorf("G1[%d]: expected %q, got %q", i, exp, result.Groups[1][i])
		}
	}
}

func TestGroupStage_SnakeDistribution(t *testing.T) {
	ids := makePlayerIDsUnsorted(8)
	numGroups := 2
	opts := &GroupStageOptions{
		NumGroups:    &numGroups,
		Distribution: "snake",
	}
	result, err := GenerateGroupStage(ids, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Snake with 2 groups:
	// Cycle 0 (forward): p1→G0, p2→G1
	// Cycle 1 (reverse): p3→G1, p4→G0
	// Cycle 2 (forward): p5→G0, p6→G1
	// Cycle 3 (reverse): p7→G1, p8→G0
	expectedG0 := []string{"p1", "p4", "p5", "p8"}
	expectedG1 := []string{"p2", "p3", "p6", "p7"}

	for i, exp := range expectedG0 {
		if i >= len(result.Groups[0]) || result.Groups[0][i] != exp {
			t.Errorf("G0[%d]: expected %q, got %q", i, exp, result.Groups[0][i])
		}
	}
	for i, exp := range expectedG1 {
		if i >= len(result.Groups[1]) || result.Groups[1][i] != exp {
			t.Errorf("G1[%d]: expected %q, got %q", i, exp, result.Groups[1][i])
		}
	}
}

func TestGroupStage_9Players_3Groups(t *testing.T) {
	ids := makePlayerIDsUnsorted(9)
	numGroups := 3
	opts := &GroupStageOptions{NumGroups: &numGroups}
	result, err := GenerateGroupStage(ids, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 3 groups of 3 players each
	if len(result.Groups) != 3 {
		t.Fatalf("expected 3 groups, got %d", len(result.Groups))
	}
	for i, g := range result.Groups {
		if len(g) != 3 {
			t.Errorf("group %d: expected 3 players, got %d", i, len(g))
		}
	}

	// Each group of 3 → C(3,2) = 3 matches (odd: 3 rounds x 1 match)
	// Total = 9 matches
	if len(result.Matches) != 9 {
		t.Fatalf("expected 9 matches, got %d", len(result.Matches))
	}
}

func TestGroupStage_BracketTypes(t *testing.T) {
	ids := makePlayerIDsUnsorted(9)
	numGroups := 3
	opts := &GroupStageOptions{NumGroups: &numGroups}
	result, err := GenerateGroupStage(ids, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, m := range result.Matches {
		bt := string(m.BracketType)
		if !strings.HasPrefix(bt, "group_") {
			t.Errorf("match %s: expected bracket_type starting with 'group_', got %q", m.MatchID, bt)
		}
	}

	// Verify each group bracket type exists
	for i := 0; i < 3; i++ {
		bt := GroupBracketType(i)
		count := countMatchesByBracketType(result.Matches, bt)
		if count != 3 {
			t.Errorf("group_%d: expected 3 matches, got %d", i, count)
		}
	}
}

func TestGroupStage_MatchIDFormat(t *testing.T) {
	ids := makePlayerIDsUnsorted(8)
	numGroups := 2
	opts := &GroupStageOptions{NumGroups: &numGroups}
	result, err := GenerateGroupStage(ids, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Group 0 matches should start with "G0-"
	// Group 1 matches should start with "G1-"
	for _, m := range result.Matches {
		bt := string(m.BracketType)
		groupNum := bt[len("group_"):]
		expectedPrefix := fmt.Sprintf("G%s-", groupNum)
		if !strings.HasPrefix(m.MatchID, expectedPrefix) {
			t.Errorf("match %s: expected prefix %q for bracket_type %q", m.MatchID, expectedPrefix, bt)
		}
	}

	// First match of group 0 should be "G0-R1-M1"
	if result.Matches[0].MatchID != "G0-R1-M1" {
		t.Errorf("first match: expected 'G0-R1-M1', got %q", result.Matches[0].MatchID)
	}
}

func TestGroupStage_GroupsContainCorrectPlayers(t *testing.T) {
	ids := makePlayerIDsUnsorted(8)
	numGroups := 2
	opts := &GroupStageOptions{NumGroups: &numGroups}
	result, err := GenerateGroupStage(ids, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// All players should be assigned to exactly one group
	allPlayers := make(map[string]int)
	for gi, group := range result.Groups {
		for _, p := range group {
			allPlayers[p] = gi
		}
	}

	if len(allPlayers) != 8 {
		t.Errorf("expected 8 players across all groups, got %d", len(allPlayers))
	}

	// Each original player should appear
	for _, id := range ids {
		if _, ok := allPlayers[id]; !ok {
			t.Errorf("player %q not assigned to any group", id)
		}
	}
}

func TestGroupStage_NoMatchLinking(t *testing.T) {
	ids := makePlayerIDsUnsorted(8)
	numGroups := 2
	opts := &GroupStageOptions{NumGroups: &numGroups}
	result, err := GenerateGroupStage(ids, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for i, m := range result.Matches {
		if m.NextMatchIndex != nil {
			t.Errorf("match %d: expected nil next_match_index, got %d", i, *m.NextMatchIndex)
		}
		if m.LoserNextMatchIndex != nil {
			t.Errorf("match %d: expected nil loser_next_match_index, got %d", i, *m.LoserNextMatchIndex)
		}
	}
}

func TestGroupStage_AutoGroupCount(t *testing.T) {
	// 12 players with no numGroups specified: round(12/4) = 3 groups
	ids := makePlayerIDsUnsorted(12)
	result, err := GenerateGroupStage(ids, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Groups) != 3 {
		t.Errorf("auto group count: expected 3 groups for 12 players, got %d", len(result.Groups))
	}
}

func TestGroupStage_AutoGroupCount_SmallN(t *testing.T) {
	// 2 players: max(2, round(2/4)) = max(2, 1) = 2 groups
	// But 1 player per group → groups with < 2 members produce 0 matches
	ids := makePlayerIDsUnsorted(2)
	result, err := GenerateGroupStage(ids, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should have 2 groups but only groups with >= 2 members generate matches
	if len(result.Groups) < 1 {
		t.Errorf("expected at least 1 group, got %d", len(result.Groups))
	}
}

func TestGroupStage_DoubleRoundRobin(t *testing.T) {
	ids := makePlayerIDsUnsorted(8)
	numGroups := 2
	opts := &GroupStageOptions{
		NumGroups:        &numGroups,
		DoubleRoundRobin: true,
	}
	result, err := GenerateGroupStage(ids, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 2 groups of 4, double RR: each group has 6*2=12 matches, total=24
	if len(result.Matches) != 24 {
		t.Fatalf("expected 24 matches for double round-robin, got %d", len(result.Matches))
	}
}

func TestGroupStage_BestOf(t *testing.T) {
	ids := makePlayerIDsUnsorted(8)
	numGroups := 2
	opts := &GroupStageOptions{
		NumGroups: &numGroups,
		BestOf:    &BestOfConfig{Default: intPtr(3)},
	}
	result, err := GenerateGroupStage(ids, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for i, m := range result.Matches {
		if m.BestOf == nil || *m.BestOf != 3 {
			t.Errorf("match %d: expected best_of=3, got %v", i, m.BestOf)
		}
	}
}

func TestGroupStage_LessThan2_ReturnsEmpty(t *testing.T) {
	result, err := GenerateGroupStage([]string{"p1"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("expected non-nil result, got nil")
	}
	if len(result.Groups) != 0 {
		t.Errorf("expected 0 groups, got %d", len(result.Groups))
	}
	if len(result.Matches) != 0 {
		t.Errorf("expected 0 matches, got %d", len(result.Matches))
	}
}

func TestGroupStage_MatchNumbers_Sequential(t *testing.T) {
	ids := makePlayerIDsUnsorted(8)
	numGroups := 2
	opts := &GroupStageOptions{NumGroups: &numGroups}
	result, err := GenerateGroupStage(ids, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for i, m := range result.Matches {
		expected := i + 1
		if m.MatchNumber != expected {
			t.Errorf("match %d: expected match_number=%d, got %d", i, expected, m.MatchNumber)
		}
	}
}

func TestGroupStage_EachGroupPairsPlayOnce(t *testing.T) {
	ids := makePlayerIDsUnsorted(8)
	numGroups := 2
	opts := &GroupStageOptions{NumGroups: &numGroups}
	result, err := GenerateGroupStage(ids, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// For each group, extract matches and verify all pairs play once
	for gi, group := range result.Groups {
		bt := GroupBracketType(gi)
		var groupMatches []MatchSeed
		for _, m := range result.Matches {
			if m.BracketType == bt {
				groupMatches = append(groupMatches, m)
			}
		}
		assertAllPairsPlayOnce(t, groupMatches, group)
	}
}
