package bracket

import (
	"strings"
	"testing"
)

func TestSingleElim_4Players(t *testing.T) {
	ids := makePlayerIDsUnsorted(4)
	matches, err := GenerateSingleElimination(ids, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 4 players → bracketSize=4, totalRounds=2
	// R1: 2 matches, R2 (Final): 1 match = 3 total
	if len(matches) != 3 {
		t.Fatalf("expected 3 matches, got %d", len(matches))
	}

	// R1 matches should have players seeded: seed order [1,4,2,3]
	// Match 0: p1 vs p4
	// Match 1: p2 vs p3
	r1 := matches[0]
	if r1.Player1ID == nil || *r1.Player1ID != "p1" {
		t.Errorf("R1 M1 player1: expected 'p1', got %v", r1.Player1ID)
	}
	if r1.Player2ID == nil || *r1.Player2ID != "p4" {
		t.Errorf("R1 M1 player2: expected 'p4', got %v", r1.Player2ID)
	}

	r1m2 := matches[1]
	if r1m2.Player1ID == nil || *r1m2.Player1ID != "p2" {
		t.Errorf("R1 M2 player1: expected 'p2', got %v", r1m2.Player1ID)
	}
	if r1m2.Player2ID == nil || *r1m2.Player2ID != "p3" {
		t.Errorf("R1 M2 player2: expected 'p3', got %v", r1m2.Player2ID)
	}

	// R2 (Final) should have nil players
	final := matches[2]
	if final.Player1ID != nil || final.Player2ID != nil {
		t.Errorf("Final should have nil players, got p1=%v p2=%v", final.Player1ID, final.Player2ID)
	}

	// Verify round names
	if r1.RoundName == nil || *r1.RoundName != "Semi-final" {
		t.Errorf("R1 round name: expected 'Semi-final', got %v", r1.RoundName)
	}
	if final.RoundName == nil || *final.RoundName != "Final" {
		t.Errorf("Final round name: expected 'Final', got %v", final.RoundName)
	}
}

func TestSingleElim_8Players(t *testing.T) {
	ids := makePlayerIDsUnsorted(8)
	matches, err := GenerateSingleElimination(ids, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 8 players → bracketSize=8, totalRounds=3
	// R1: 4, R2: 2, R3: 1 = 7 total
	if len(matches) != 7 {
		t.Fatalf("expected 7 matches, got %d", len(matches))
	}

	// Verify round distribution
	r1Count := countMatchesByRound(matches, 1)
	r2Count := countMatchesByRound(matches, 2)
	r3Count := countMatchesByRound(matches, 3)
	if r1Count != 4 {
		t.Errorf("R1: expected 4 matches, got %d", r1Count)
	}
	if r2Count != 2 {
		t.Errorf("R2: expected 2 matches, got %d", r2Count)
	}
	if r3Count != 1 {
		t.Errorf("R3: expected 1 match, got %d", r3Count)
	}

	// Verify seeding: seed order [1,8,4,5,2,7,3,6]
	// R1 M1: p1 vs p8
	// R1 M2: p4 vs p5
	// R1 M3: p2 vs p7
	// R1 M4: p3 vs p6
	expectedSeeds := [][2]string{
		{"p1", "p8"},
		{"p4", "p5"},
		{"p2", "p7"},
		{"p3", "p6"},
	}
	for i, exp := range expectedSeeds {
		m := matches[i]
		if m.Player1ID == nil || *m.Player1ID != exp[0] {
			t.Errorf("R1 M%d player1: expected %q, got %v", i+1, exp[0], m.Player1ID)
		}
		if m.Player2ID == nil || *m.Player2ID != exp[1] {
			t.Errorf("R1 M%d player2: expected %q, got %v", i+1, exp[1], m.Player2ID)
		}
	}

	// Verify round names
	if matches[0].RoundName == nil || *matches[0].RoundName != "Quarter-final" {
		t.Errorf("R1 round name: expected 'Quarter-final', got %v", matches[0].RoundName)
	}
	if matches[4].RoundName == nil || *matches[4].RoundName != "Semi-final" {
		t.Errorf("R2 round name: expected 'Semi-final', got %v", matches[4].RoundName)
	}
	if matches[6].RoundName == nil || *matches[6].RoundName != "Final" {
		t.Errorf("R3 round name: expected 'Final', got %v", matches[6].RoundName)
	}

	// Verify no byes
	for i, m := range matches[:4] {
		if m.IsBye {
			t.Errorf("R1 M%d: unexpected bye", i+1)
		}
	}
}

func TestSingleElim_3Players_Byes(t *testing.T) {
	ids := makePlayerIDsUnsorted(3)
	matches, err := GenerateSingleElimination(ids, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 3 players → bracketSize=4, totalRounds=2
	// R1: 2, R2: 1 = 3 total
	if len(matches) != 3 {
		t.Fatalf("expected 3 matches, got %d", len(matches))
	}

	// 1 bye (seed 4 is nil)
	byeCount := 0
	for _, m := range matches[:2] {
		if m.IsBye {
			byeCount++
		}
	}
	if byeCount != 1 {
		t.Errorf("expected 1 bye match, got %d", byeCount)
	}
}

func TestSingleElim_5Players_Byes(t *testing.T) {
	ids := makePlayerIDsUnsorted(5)
	matches, err := GenerateSingleElimination(ids, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 5 players → bracketSize=8, totalRounds=3
	// R1: 4, R2: 2, R3: 1 = 7 total
	if len(matches) != 7 {
		t.Fatalf("expected 7 matches, got %d", len(matches))
	}

	// 3 byes (seeds 6,7,8 are nil)
	byeCount := 0
	for _, m := range matches[:4] {
		if m.IsBye {
			byeCount++
		}
	}
	if byeCount != 3 {
		t.Errorf("expected 3 bye matches, got %d", byeCount)
	}
}

func TestSingleElim_2Players(t *testing.T) {
	ids := makePlayerIDsUnsorted(2)
	matches, err := GenerateSingleElimination(ids, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 2 players → bracketSize=2, totalRounds=1
	// Only 1 match
	if len(matches) != 1 {
		t.Fatalf("expected 1 match, got %d", len(matches))
	}

	m := matches[0]
	if m.Player1ID == nil || *m.Player1ID != "p1" {
		t.Errorf("player1: expected 'p1', got %v", m.Player1ID)
	}
	if m.Player2ID == nil || *m.Player2ID != "p2" {
		t.Errorf("player2: expected 'p2', got %v", m.Player2ID)
	}
	if m.RoundName == nil || *m.RoundName != "Final" {
		t.Errorf("round name: expected 'Final', got %v", m.RoundName)
	}
	if m.NextMatchIndex != nil {
		t.Errorf("2-player final should have no next match, got %v", m.NextMatchIndex)
	}
	if m.IsBye {
		t.Error("2-player match should not be a bye")
	}
}

func TestSingleElim_WithThirdPlace(t *testing.T) {
	ids := makePlayerIDsUnsorted(4)
	opts := &SingleEliminationOptions{ThirdPlaceMatch: true}
	matches, err := GenerateSingleElimination(ids, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 4 players + 3rd place = 3 + 1 = 4 matches
	if len(matches) != 4 {
		t.Fatalf("expected 4 matches, got %d", len(matches))
	}

	// Last match should be third place
	thirdPlace := matches[3]
	if thirdPlace.BracketType != BracketThirdPlace {
		t.Errorf("expected bracket_type 'third_place', got %q", thirdPlace.BracketType)
	}
	if thirdPlace.MatchID != "3RD-M1" {
		t.Errorf("expected match_id '3RD-M1', got %q", thirdPlace.MatchID)
	}
	if thirdPlace.RoundName == nil || *thirdPlace.RoundName != "3rd Place Match" {
		t.Errorf("expected round name '3rd Place Match', got %v", thirdPlace.RoundName)
	}

	// Semi-final losers should link to 3rd place match
	semi1 := matches[0]
	semi2 := matches[1]
	if semi1.LoserNextMatchIndex == nil || *semi1.LoserNextMatchIndex != 3 {
		t.Errorf("semi1 loser should link to index 3, got %v", semi1.LoserNextMatchIndex)
	}
	if semi2.LoserNextMatchIndex == nil || *semi2.LoserNextMatchIndex != 3 {
		t.Errorf("semi2 loser should link to index 3, got %v", semi2.LoserNextMatchIndex)
	}

	// Loser slots: first semi → player1, second semi → player2
	if semi1.LoserNextMatchSlot == nil || *semi1.LoserNextMatchSlot != SlotPlayer1 {
		t.Errorf("semi1 loser slot: expected 'player1', got %v", semi1.LoserNextMatchSlot)
	}
	if semi2.LoserNextMatchSlot == nil || *semi2.LoserNextMatchSlot != SlotPlayer2 {
		t.Errorf("semi2 loser slot: expected 'player2', got %v", semi2.LoserNextMatchSlot)
	}
}

func TestSingleElim_WithThirdPlace_8Players(t *testing.T) {
	ids := makePlayerIDsUnsorted(8)
	opts := &SingleEliminationOptions{ThirdPlaceMatch: true}
	matches, err := GenerateSingleElimination(ids, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 7 normal + 1 third place = 8
	if len(matches) != 8 {
		t.Fatalf("expected 8 matches, got %d", len(matches))
	}

	// Verify third place match exists
	thirdPlaceCount := countMatchesByBracketType(matches, BracketThirdPlace)
	if thirdPlaceCount != 1 {
		t.Errorf("expected 1 third place match, got %d", thirdPlaceCount)
	}

	// Semi-finals are in round 2 (totalRounds=3, semi=round 2)
	semiIndices := []int{}
	for i, m := range matches {
		if m.Round == 2 && m.BracketType == BracketWinners {
			semiIndices = append(semiIndices, i)
		}
	}
	if len(semiIndices) != 2 {
		t.Fatalf("expected 2 semi-final matches, got %d", len(semiIndices))
	}

	// Both semi-finals should have loser links to third place match (index 7)
	for _, idx := range semiIndices {
		m := matches[idx]
		if m.LoserNextMatchIndex == nil || *m.LoserNextMatchIndex != 7 {
			t.Errorf("semi at index %d: loser should link to index 7, got %v", idx, m.LoserNextMatchIndex)
		}
	}
}

func TestSingleElim_MatchLinking(t *testing.T) {
	ids := makePlayerIDsUnsorted(8)
	matches, err := GenerateSingleElimination(ids, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// R1 (indices 0-3) should link to R2 (indices 4-5)
	// Match 0 → index 4 (player1)
	// Match 1 → index 4 (player2)
	// Match 2 → index 5 (player1)
	// Match 3 → index 5 (player2)
	expectedLinks := []struct {
		fromIdx       int
		toIdx         int
		slot          MatchSlot
	}{
		{0, 4, SlotPlayer1},
		{1, 4, SlotPlayer2},
		{2, 5, SlotPlayer1},
		{3, 5, SlotPlayer2},
	}

	for _, link := range expectedLinks {
		m := matches[link.fromIdx]
		if m.NextMatchIndex == nil || *m.NextMatchIndex != link.toIdx {
			t.Errorf("match %d: expected next_match_index=%d, got %v", link.fromIdx, link.toIdx, m.NextMatchIndex)
		}
		if m.NextMatchSlot == nil || *m.NextMatchSlot != link.slot {
			t.Errorf("match %d: expected next_match_slot=%q, got %v", link.fromIdx, link.slot, m.NextMatchSlot)
		}
	}

	// R2 (indices 4-5) should link to R3/Final (index 6)
	// Match 4 → index 6 (player1)
	// Match 5 → index 6 (player2)
	if matches[4].NextMatchIndex == nil || *matches[4].NextMatchIndex != 6 {
		t.Errorf("match 4: expected next_match_index=6, got %v", matches[4].NextMatchIndex)
	}
	if matches[4].NextMatchSlot == nil || *matches[4].NextMatchSlot != SlotPlayer1 {
		t.Errorf("match 4: expected slot player1, got %v", matches[4].NextMatchSlot)
	}
	if matches[5].NextMatchIndex == nil || *matches[5].NextMatchIndex != 6 {
		t.Errorf("match 5: expected next_match_index=6, got %v", matches[5].NextMatchIndex)
	}
	if matches[5].NextMatchSlot == nil || *matches[5].NextMatchSlot != SlotPlayer2 {
		t.Errorf("match 5: expected slot player2, got %v", matches[5].NextMatchSlot)
	}

	// Final (index 6) should have no next match
	if matches[6].NextMatchIndex != nil {
		t.Errorf("final: expected no next match, got %v", matches[6].NextMatchIndex)
	}
}

func TestSingleElim_SlotAssignment(t *testing.T) {
	ids := makePlayerIDsUnsorted(4)
	matches, err := GenerateSingleElimination(ids, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Match 0 (even index in round) → player1 slot
	// Match 1 (odd index in round) → player2 slot
	if matches[0].NextMatchSlot == nil || *matches[0].NextMatchSlot != SlotPlayer1 {
		t.Errorf("match 0: expected slot player1, got %v", matches[0].NextMatchSlot)
	}
	if matches[1].NextMatchSlot == nil || *matches[1].NextMatchSlot != SlotPlayer2 {
		t.Errorf("match 1: expected slot player2, got %v", matches[1].NextMatchSlot)
	}
}

func TestSingleElim_MatchIDFormat(t *testing.T) {
	ids := makePlayerIDsUnsorted(8)
	matches, err := GenerateSingleElimination(ids, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// R1: WB-R1-M1 through WB-R1-M4
	for i := 0; i < 4; i++ {
		if !strings.HasPrefix(matches[i].MatchID, "WB-R1-") {
			t.Errorf("match %d: expected prefix 'WB-R1-', got %q", i, matches[i].MatchID)
		}
	}
	// R2: WB-R2-M1, WB-R2-M2
	if matches[4].MatchID != "WB-R2-M1" {
		t.Errorf("match 4: expected 'WB-R2-M1', got %q", matches[4].MatchID)
	}
	if matches[5].MatchID != "WB-R2-M2" {
		t.Errorf("match 5: expected 'WB-R2-M2', got %q", matches[5].MatchID)
	}
	// R3 Final: WB-R3-M1
	if matches[6].MatchID != "WB-R3-M1" {
		t.Errorf("match 6: expected 'WB-R3-M1', got %q", matches[6].MatchID)
	}
}

func TestSingleElim_AllWinnersBracketType(t *testing.T) {
	ids := makePlayerIDsUnsorted(8)
	matches, err := GenerateSingleElimination(ids, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for i, m := range matches {
		if m.BracketType != BracketWinners {
			t.Errorf("match %d: expected bracket_type 'winners', got %q", i, m.BracketType)
		}
	}
}

func TestSingleElim_BestOf(t *testing.T) {
	ids := makePlayerIDsUnsorted(4)
	bestOf := &BestOfConfig{
		Default: intPtr(3),
		Final:   intPtr(5),
	}
	opts := &SingleEliminationOptions{BestOf: bestOf}
	matches, err := GenerateSingleElimination(ids, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// R1 matches should have best_of = 3
	for i := 0; i < 2; i++ {
		if matches[i].BestOf == nil || *matches[i].BestOf != 3 {
			t.Errorf("R1 match %d: expected best_of=3, got %v", i, matches[i].BestOf)
		}
	}

	// Final should have best_of = 5
	final := matches[2]
	if final.BestOf == nil || *final.BestOf != 5 {
		t.Errorf("final: expected best_of=5, got %v", final.BestOf)
	}
}

func TestSingleElim_BestOfThirdPlace(t *testing.T) {
	ids := makePlayerIDsUnsorted(4)
	bestOf := &BestOfConfig{
		Default:    intPtr(3),
		Final:      intPtr(5),
		ThirdPlace: intPtr(1),
	}
	opts := &SingleEliminationOptions{ThirdPlaceMatch: true, BestOf: bestOf}
	matches, err := GenerateSingleElimination(ids, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Find third place match
	for _, m := range matches {
		if m.BracketType == BracketThirdPlace {
			if m.BestOf == nil || *m.BestOf != 1 {
				t.Errorf("third place: expected best_of=1, got %v", m.BestOf)
			}
		}
	}
}

func TestSingleElim_LessThan2_ReturnsNil(t *testing.T) {
	matches, err := GenerateSingleElimination([]string{"p1"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if matches != nil {
		t.Errorf("expected nil for <2 participants, got %v", matches)
	}
}

func TestSingleElim_16Players(t *testing.T) {
	ids := makePlayerIDsUnsorted(16)
	matches, err := GenerateSingleElimination(ids, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 16 players → 15 matches (8+4+2+1)
	if len(matches) != 15 {
		t.Fatalf("expected 15 matches, got %d", len(matches))
	}

	// Round 1 name should be "Round 1" (4 rounds from final)
	if matches[0].RoundName == nil || *matches[0].RoundName != "Round 1" {
		t.Errorf("R1 name: expected 'Round 1', got %v", matches[0].RoundName)
	}
}

func TestSingleElim_MatchNumbers_Sequential(t *testing.T) {
	ids := makePlayerIDsUnsorted(8)
	matches, err := GenerateSingleElimination(ids, nil)
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
