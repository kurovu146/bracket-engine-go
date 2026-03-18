package bracket

import (
	"testing"
)

func TestRoundRobin_4Players(t *testing.T) {
	ids := makePlayerIDsUnsorted(4)
	matches, err := GenerateRoundRobin(ids, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 4 players (even): 3 rounds x 2 matches = 6 matches
	if len(matches) != 6 {
		t.Fatalf("expected 6 matches, got %d", len(matches))
	}

	// Verify 3 rounds
	rounds := make(map[int]int)
	for _, m := range matches {
		rounds[m.Round]++
	}
	if len(rounds) != 3 {
		t.Errorf("expected 3 rounds, got %d", len(rounds))
	}
	for r, count := range rounds {
		if count != 2 {
			t.Errorf("round %d: expected 2 matches, got %d", r, count)
		}
	}

	// Every pair plays exactly once
	assertAllPairsPlayOnce(t, matches, ids)
}

func TestRoundRobin_3Players_Odd(t *testing.T) {
	ids := makePlayerIDsUnsorted(3)
	matches, err := GenerateRoundRobin(ids, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 3 players (odd): bye sentinel added → 4 participants internally
	// 3 rounds x (2 matches per round - bye matches)
	// Each round has 1 bye match (filtered out) → 1 real match per round
	// Total: 3 matches
	if len(matches) != 3 {
		t.Fatalf("expected 3 matches, got %d", len(matches))
	}

	// Every pair plays exactly once: C(3,2)=3 pairs
	assertAllPairsPlayOnce(t, matches, ids)
}

func TestRoundRobin_5Players_Odd(t *testing.T) {
	ids := makePlayerIDsUnsorted(5)
	matches, err := GenerateRoundRobin(ids, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 5 players (odd): bye sentinel → 6 internal, 5 rounds x 3 matches - 5 bye matches
	// Each round has 1 bye → 2 real matches per round → 5*2=10
	if len(matches) != 10 {
		t.Fatalf("expected 10 matches, got %d", len(matches))
	}

	// Every pair plays exactly once: C(5,2)=10 pairs
	assertAllPairsPlayOnce(t, matches, ids)
}

func TestRoundRobin_DoubleRoundRobin(t *testing.T) {
	ids := makePlayerIDsUnsorted(4)
	opts := &RoundRobinOptions{DoubleRoundRobin: true}
	matches, err := GenerateRoundRobin(ids, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 4 players double: 6 first pass + 6 second pass = 12 matches
	if len(matches) != 12 {
		t.Fatalf("expected 12 matches, got %d", len(matches))
	}

	// Every pair plays exactly twice
	assertAllPairsPlayTwice(t, matches, ids)

	// Second pass should have swapped player1/player2
	// First pass matches are indices 0-5, second pass 6-11
	// Find a match in first pass and its counterpart in second pass
	firstPassPairs := make(map[string][2]string)
	for _, m := range matches[:6] {
		if m.Player1ID != nil && m.Player2ID != nil {
			key := sortedPair(*m.Player1ID, *m.Player2ID)
			firstPassPairs[key] = [2]string{*m.Player1ID, *m.Player2ID}
		}
	}

	secondPassPairs := make(map[string][2]string)
	for _, m := range matches[6:] {
		if m.Player1ID != nil && m.Player2ID != nil {
			key := sortedPair(*m.Player1ID, *m.Player2ID)
			secondPassPairs[key] = [2]string{*m.Player1ID, *m.Player2ID}
		}
	}

	// For each pair, player positions should be swapped in second pass
	for key, first := range firstPassPairs {
		second, ok := secondPassPairs[key]
		if !ok {
			t.Errorf("pair %s not found in second pass", key)
			continue
		}
		if first[0] != second[1] || first[1] != second[0] {
			t.Errorf("pair %s: first pass (%s vs %s), second pass (%s vs %s) — expected swapped",
				key, first[0], first[1], second[0], second[1])
		}
	}
}

func TestRoundRobin_NoMatchLinking(t *testing.T) {
	ids := makePlayerIDsUnsorted(4)
	matches, err := GenerateRoundRobin(ids, nil)
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

func TestRoundRobin_BracketType(t *testing.T) {
	ids := makePlayerIDsUnsorted(4)
	matches, err := GenerateRoundRobin(ids, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for i, m := range matches {
		if m.BracketType != BracketRoundRobin {
			t.Errorf("match %d: expected bracket_type 'round_robin', got %q", i, m.BracketType)
		}
	}
}

func TestRoundRobin_MatchIDFormat(t *testing.T) {
	ids := makePlayerIDsUnsorted(4)
	matches, err := GenerateRoundRobin(ids, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// First match should be RR-R1-M1
	if matches[0].MatchID != "RR-R1-M1" {
		t.Errorf("first match_id: expected 'RR-R1-M1', got %q", matches[0].MatchID)
	}
}

func TestRoundRobin_RoundNames(t *testing.T) {
	ids := makePlayerIDsUnsorted(4)
	matches, err := GenerateRoundRobin(ids, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, m := range matches {
		if m.RoundName == nil {
			t.Errorf("match %s: expected round name, got nil", m.MatchID)
			continue
		}
		// Round robin names should be "Round N"
		if m.Round == 1 && *m.RoundName != "Round 1" {
			t.Errorf("round 1 name: expected 'Round 1', got %q", *m.RoundName)
		}
	}
}

func TestRoundRobin_NoByes(t *testing.T) {
	ids := makePlayerIDsUnsorted(4)
	matches, err := GenerateRoundRobin(ids, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for i, m := range matches {
		if m.IsBye {
			t.Errorf("match %d: unexpected bye", i)
		}
		if m.Player1ID == nil || m.Player2ID == nil {
			t.Errorf("match %d: expected both players, got p1=%v p2=%v", i, m.Player1ID, m.Player2ID)
		}
	}
}

func TestRoundRobin_BestOf(t *testing.T) {
	ids := makePlayerIDsUnsorted(4)
	opts := &RoundRobinOptions{
		BestOf: &BestOfConfig{Default: intPtr(3)},
	}
	matches, err := GenerateRoundRobin(ids, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for i, m := range matches {
		if m.BestOf == nil || *m.BestOf != 3 {
			t.Errorf("match %d: expected best_of=3, got %v", i, m.BestOf)
		}
	}
}

func TestRoundRobin_LessThan2_ReturnsNil(t *testing.T) {
	matches, err := GenerateRoundRobin([]string{"p1"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if matches != nil {
		t.Errorf("expected nil for <2 participants, got %v", matches)
	}
}

func TestRoundRobin_2Players(t *testing.T) {
	ids := makePlayerIDsUnsorted(2)
	matches, err := GenerateRoundRobin(ids, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 2 players: 1 round x 1 match = 1 match
	if len(matches) != 1 {
		t.Fatalf("expected 1 match, got %d", len(matches))
	}

	assertAllPairsPlayOnce(t, matches, ids)
}

func TestRoundRobin_6Players(t *testing.T) {
	ids := makePlayerIDsUnsorted(6)
	matches, err := GenerateRoundRobin(ids, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 6 players (even): 5 rounds x 3 matches = 15 matches
	// C(6,2) = 15 pairs
	if len(matches) != 15 {
		t.Fatalf("expected 15 matches, got %d", len(matches))
	}

	assertAllPairsPlayOnce(t, matches, ids)
}
