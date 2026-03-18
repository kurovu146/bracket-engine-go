package bracket

import (
	"testing"
)

func TestSwiss_8Players(t *testing.T) {
	ids := makePlayerIDsUnsorted(8)
	matches, err := GenerateSwiss(ids, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 8 players: ceil(log2(8)) = 3 rounds, 4 matches per round = 12 total
	if len(matches) != 12 {
		t.Fatalf("expected 12 matches, got %d", len(matches))
	}

	// Verify round distribution
	r1Count := countMatchesByRound(matches, 1)
	r2Count := countMatchesByRound(matches, 2)
	r3Count := countMatchesByRound(matches, 3)
	if r1Count != 4 {
		t.Errorf("R1: expected 4 matches, got %d", r1Count)
	}
	if r2Count != 4 {
		t.Errorf("R2: expected 4 matches, got %d", r2Count)
	}
	if r3Count != 4 {
		t.Errorf("R3: expected 4 matches, got %d", r3Count)
	}
}

func TestSwiss_Round1_HasRealPlayers(t *testing.T) {
	ids := makePlayerIDsUnsorted(8)
	matches, err := GenerateSwiss(ids, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Round 1: paired by seed order (p1 vs p2, p3 vs p4, p5 vs p6, p7 vs p8)
	r1Matches := []MatchSeed{}
	for _, m := range matches {
		if m.Round == 1 {
			r1Matches = append(r1Matches, m)
		}
	}

	expectedPairs := [][2]string{
		{"p1", "p2"},
		{"p3", "p4"},
		{"p5", "p6"},
		{"p7", "p8"},
	}

	for i, exp := range expectedPairs {
		m := r1Matches[i]
		if m.Player1ID == nil || *m.Player1ID != exp[0] {
			t.Errorf("R1 M%d player1: expected %q, got %v", i+1, exp[0], m.Player1ID)
		}
		if m.Player2ID == nil || *m.Player2ID != exp[1] {
			t.Errorf("R1 M%d player2: expected %q, got %v", i+1, exp[1], m.Player2ID)
		}
	}
}

func TestSwiss_Round2Plus_NilPlayers(t *testing.T) {
	ids := makePlayerIDsUnsorted(8)
	matches, err := GenerateSwiss(ids, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, m := range matches {
		if m.Round >= 2 {
			if m.Player1ID != nil {
				t.Errorf("round %d match %s: expected nil player1, got %q", m.Round, m.MatchID, *m.Player1ID)
			}
			if m.Player2ID != nil {
				t.Errorf("round %d match %s: expected nil player2, got %q", m.Round, m.MatchID, *m.Player2ID)
			}
		}
	}
}

func TestSwiss_CustomRounds(t *testing.T) {
	ids := makePlayerIDsUnsorted(8)
	numRounds := 5
	opts := &SwissOptions{NumRounds: &numRounds}
	matches, err := GenerateSwiss(ids, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 5 rounds x 4 matches = 20
	if len(matches) != 20 {
		t.Fatalf("expected 20 matches, got %d", len(matches))
	}

	// Verify all 5 rounds exist
	rounds := make(map[int]bool)
	for _, m := range matches {
		rounds[m.Round] = true
	}
	for r := 1; r <= 5; r++ {
		if !rounds[r] {
			t.Errorf("round %d missing", r)
		}
	}
}

func TestSwiss_3Players(t *testing.T) {
	ids := makePlayerIDsUnsorted(3)
	matches, err := GenerateSwiss(ids, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 3 players: ceil(log2(3)) = 2 rounds, n/2 = 1 match per round = 2 total
	if len(matches) != 2 {
		t.Fatalf("expected 2 matches, got %d", len(matches))
	}

	// R1: 1 match with players
	r1 := matches[0]
	if r1.Round != 1 {
		t.Errorf("first match should be round 1, got %d", r1.Round)
	}
	if r1.Player1ID == nil || r1.Player2ID == nil {
		t.Error("R1 should have both players")
	}

	// R2: 1 placeholder match
	r2 := matches[1]
	if r2.Round != 2 {
		t.Errorf("second match should be round 2, got %d", r2.Round)
	}
	if r2.Player1ID != nil || r2.Player2ID != nil {
		t.Error("R2 should have nil players")
	}
}

func TestSwiss_BracketType(t *testing.T) {
	ids := makePlayerIDsUnsorted(8)
	matches, err := GenerateSwiss(ids, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for i, m := range matches {
		if m.BracketType != BracketSwiss {
			t.Errorf("match %d: expected bracket_type 'swiss', got %q", i, m.BracketType)
		}
	}
}

func TestSwiss_MatchIDFormat(t *testing.T) {
	ids := makePlayerIDsUnsorted(8)
	matches, err := GenerateSwiss(ids, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if matches[0].MatchID != "SW-R1-M1" {
		t.Errorf("first match_id: expected 'SW-R1-M1', got %q", matches[0].MatchID)
	}
	if matches[4].MatchID != "SW-R2-M1" {
		t.Errorf("R2 M1 match_id: expected 'SW-R2-M1', got %q", matches[4].MatchID)
	}
}

func TestSwiss_NoMatchLinking(t *testing.T) {
	ids := makePlayerIDsUnsorted(8)
	matches, err := GenerateSwiss(ids, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for i, m := range matches {
		if m.NextMatchIndex != nil {
			t.Errorf("match %d: expected nil next_match_index, got %d", i, *m.NextMatchIndex)
		}
		if m.LoserNextMatchIndex != nil {
			t.Errorf("match %d: expected nil loser_next_match_index, got %d", i, *m.LoserNextMatchIndex)
		}
	}
}

func TestSwiss_RoundNames(t *testing.T) {
	ids := makePlayerIDsUnsorted(8)
	matches, err := GenerateSwiss(ids, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, m := range matches {
		if m.RoundName == nil {
			t.Errorf("match %s: expected round name", m.MatchID)
			continue
		}
		// Swiss uses "Round N" format
		if m.Round == 1 && *m.RoundName != "Round 1" {
			t.Errorf("R1 name: expected 'Round 1', got %q", *m.RoundName)
		}
		if m.Round == 3 && *m.RoundName != "Round 3" {
			t.Errorf("R3 name: expected 'Round 3', got %q", *m.RoundName)
		}
	}
}

func TestSwiss_BestOf(t *testing.T) {
	ids := makePlayerIDsUnsorted(4)
	opts := &SwissOptions{
		BestOf: &BestOfConfig{Default: intPtr(5)},
	}
	matches, err := GenerateSwiss(ids, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for i, m := range matches {
		if m.BestOf == nil || *m.BestOf != 5 {
			t.Errorf("match %d: expected best_of=5, got %v", i, m.BestOf)
		}
	}
}

func TestSwiss_LessThan2_ReturnsNil(t *testing.T) {
	matches, err := GenerateSwiss([]string{"p1"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if matches != nil {
		t.Errorf("expected nil for <2 participants, got %v", matches)
	}
}

func TestSwiss_4Players(t *testing.T) {
	ids := makePlayerIDsUnsorted(4)
	matches, err := GenerateSwiss(ids, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 4 players: ceil(log2(4)) = 2 rounds, 2 matches per round = 4 total
	if len(matches) != 4 {
		t.Fatalf("expected 4 matches, got %d", len(matches))
	}
}

func TestSwiss_MatchNumbers_Sequential(t *testing.T) {
	ids := makePlayerIDsUnsorted(8)
	matches, err := GenerateSwiss(ids, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for i, m := range matches {
		expected := i + 1
		if m.MatchNumber != expected {
			t.Errorf("match %d: expected match_number=%d, got %d", i, expected, m.MatchNumber)
		}
	}
}
