package bracket

import (
	"testing"
)

func TestDoubleElim_4Players(t *testing.T) {
	ids := makePlayerIDsUnsorted(4)
	matches, err := GenerateDoubleElimination(ids, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 4 players → bracketSize=4, totalWinnersRounds=2
	// Winners: R1(2) + R2(1) = 3
	// Losers: totalLosersRounds = (2-1)*2 = 2 rounds
	//   L1: bracketSize/4 = 1 match
	//   L2: 1 match (count halves after even round, but min 1)
	// Grand final: 1
	// Total: 3 + 2 + 1 = 6

	winnersCount := countMatchesByBracketType(matches, BracketWinners)
	losersCount := countMatchesByBracketType(matches, BracketLosers)
	gfCount := countMatchesByBracketType(matches, BracketGrandFinal)

	if winnersCount != 3 {
		t.Errorf("expected 3 winners matches, got %d", winnersCount)
	}
	if losersCount != 2 {
		t.Errorf("expected 2 losers matches, got %d", losersCount)
	}
	if gfCount != 1 {
		t.Errorf("expected 1 grand final, got %d", gfCount)
	}
	if len(matches) != 6 {
		t.Errorf("expected 6 total matches, got %d", len(matches))
	}
}

func TestDoubleElim_8Players(t *testing.T) {
	ids := makePlayerIDsUnsorted(8)
	matches, err := GenerateDoubleElimination(ids, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 8 players → bracketSize=8, totalWinnersRounds=3
	// Winners: R1(4) + R2(2) + R3(1) = 7
	// Losers: totalLosersRounds = (3-1)*2 = 4 rounds
	//   L1: bracketSize/4 = 2 matches
	//   L2: 2 matches (same count, no halving yet)
	//   L3: halves to 1 (halve after L2 which is even)
	//   L4: 1 match
	// Grand final: 1
	// Total: 7 + 6 + 1 = 14

	winnersCount := countMatchesByBracketType(matches, BracketWinners)
	losersCount := countMatchesByBracketType(matches, BracketLosers)
	gfCount := countMatchesByBracketType(matches, BracketGrandFinal)

	if winnersCount != 7 {
		t.Errorf("expected 7 winners matches, got %d", winnersCount)
	}
	if losersCount != 6 {
		t.Errorf("expected 6 losers matches, got %d", losersCount)
	}
	if gfCount != 1 {
		t.Errorf("expected 1 grand final, got %d", gfCount)
	}
	if len(matches) != 14 {
		t.Errorf("expected 14 total matches, got %d", len(matches))
	}
}

func TestDoubleElim_2Players(t *testing.T) {
	ids := makePlayerIDsUnsorted(2)
	matches, err := GenerateDoubleElimination(ids, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 2 players → bracketSize=2, totalWinnersRounds=1
	// Winners: 1 match
	// Losers: totalLosersRounds = (1-1)*2 = 0 (no losers bracket)
	// Grand final: 1
	// Total: 2

	winnersCount := countMatchesByBracketType(matches, BracketWinners)
	losersCount := countMatchesByBracketType(matches, BracketLosers)
	gfCount := countMatchesByBracketType(matches, BracketGrandFinal)

	if winnersCount != 1 {
		t.Errorf("expected 1 winners match, got %d", winnersCount)
	}
	if losersCount != 0 {
		t.Errorf("expected 0 losers matches, got %d", losersCount)
	}
	if gfCount != 1 {
		t.Errorf("expected 1 grand final, got %d", gfCount)
	}
	if len(matches) != 2 {
		t.Errorf("expected 2 total matches, got %d", len(matches))
	}

	// WB match loser should go to GF as player2
	wbMatch := matches[0]
	if wbMatch.LoserNextMatchIndex == nil || *wbMatch.LoserNextMatchIndex != 1 {
		t.Errorf("WB loser should link to GF (index 1), got %v", wbMatch.LoserNextMatchIndex)
	}
	if wbMatch.LoserNextMatchSlot == nil || *wbMatch.LoserNextMatchSlot != SlotPlayer2 {
		t.Errorf("WB loser slot should be player2, got %v", wbMatch.LoserNextMatchSlot)
	}

	// WB winner should go to GF as player1
	if wbMatch.NextMatchIndex == nil || *wbMatch.NextMatchIndex != 1 {
		t.Errorf("WB winner should link to GF (index 1), got %v", wbMatch.NextMatchIndex)
	}
	if wbMatch.NextMatchSlot == nil || *wbMatch.NextMatchSlot != SlotPlayer1 {
		t.Errorf("WB winner slot should be player1, got %v", wbMatch.NextMatchSlot)
	}
}

func TestDoubleElim_3Players_Byes(t *testing.T) {
	ids := makePlayerIDsUnsorted(3)
	matches, err := GenerateDoubleElimination(ids, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 3 players → bracketSize=4, same structure as 4 players but with 1 bye
	winnersCount := countMatchesByBracketType(matches, BracketWinners)
	if winnersCount != 3 {
		t.Errorf("expected 3 winners matches, got %d", winnersCount)
	}

	// Should have at least 1 bye in winners R1
	byeCount := 0
	for _, m := range matches {
		if m.BracketType == BracketWinners && m.Round == 1 && m.IsBye {
			byeCount++
		}
	}
	if byeCount != 1 {
		t.Errorf("expected 1 bye in WB R1, got %d", byeCount)
	}
}

func TestDoubleElim_GrandFinalReset(t *testing.T) {
	ids := makePlayerIDsUnsorted(4)
	opts := &DoubleEliminationOptions{GrandFinalReset: true}
	matches, err := GenerateDoubleElimination(ids, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Normal 4-player DE = 6 matches, with reset = 7
	if len(matches) != 7 {
		t.Fatalf("expected 7 matches with GF reset, got %d", len(matches))
	}

	// Should have grand_final_reset bracket type
	resetCount := countMatchesByBracketType(matches, BracketGrandFinalReset)
	if resetCount != 1 {
		t.Errorf("expected 1 grand final reset, got %d", resetCount)
	}

	// Find GF and GF Reset
	var gfIdx, resetIdx int
	for i, m := range matches {
		if m.BracketType == BracketGrandFinal {
			gfIdx = i
		}
		if m.BracketType == BracketGrandFinalReset {
			resetIdx = i
		}
	}

	// GF should link to reset
	gf := matches[gfIdx]
	if gf.NextMatchIndex == nil || *gf.NextMatchIndex != resetIdx {
		t.Errorf("GF should link to reset (index %d), got %v", resetIdx, gf.NextMatchIndex)
	}
	if gf.NextMatchSlot == nil || *gf.NextMatchSlot != SlotPlayer1 {
		t.Errorf("GF next slot should be player1, got %v", gf.NextMatchSlot)
	}

	// Reset match IDs
	reset := matches[resetIdx]
	if reset.MatchID != "GF-M2" {
		t.Errorf("reset match_id: expected 'GF-M2', got %q", reset.MatchID)
	}
	if reset.RoundName == nil || *reset.RoundName != "Grand Final Reset" {
		t.Errorf("reset round name: expected 'Grand Final Reset', got %v", reset.RoundName)
	}

	// Reset should have no next match
	if reset.NextMatchIndex != nil {
		t.Errorf("reset should have no next match, got %v", reset.NextMatchIndex)
	}
}

func TestDoubleElim_WBtoLB_DropoutLinking(t *testing.T) {
	ids := makePlayerIDsUnsorted(8)
	matches, err := GenerateDoubleElimination(ids, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// W1 (4 matches at indices 0-3) losers should pair 2:1 into L1 (2 matches)
	// W1 M0 → L1 M0 (player1)
	// W1 M1 → L1 M0 (player2)
	// W1 M2 → L1 M1 (player1)
	// W1 M3 → L1 M1 (player2)
	w1Matches := []int{0, 1, 2, 3}
	for _, idx := range w1Matches {
		m := matches[idx]
		if m.LoserNextMatchIndex == nil {
			t.Errorf("W1 match %d: expected loser link, got nil", idx)
			continue
		}
		target := matches[*m.LoserNextMatchIndex]
		if target.BracketType != BracketLosers {
			t.Errorf("W1 match %d: loser target should be losers bracket, got %q", idx, target.BracketType)
		}
	}

	// W1 pairing: even→player1, odd→player2
	if matches[0].LoserNextMatchSlot == nil || *matches[0].LoserNextMatchSlot != SlotPlayer1 {
		t.Errorf("W1 M0 loser slot: expected player1, got %v", matches[0].LoserNextMatchSlot)
	}
	if matches[1].LoserNextMatchSlot == nil || *matches[1].LoserNextMatchSlot != SlotPlayer2 {
		t.Errorf("W1 M1 loser slot: expected player2, got %v", matches[1].LoserNextMatchSlot)
	}

	// W2+ losers should go to even LB rounds as player2
	// W2 (2 matches at indices 4-5) → L2 (lRoundIndex=1, which is even LB round L2)
	for i := 4; i <= 5; i++ {
		m := matches[i]
		if m.LoserNextMatchIndex == nil {
			t.Errorf("W2 match %d: expected loser link, got nil", i)
			continue
		}
		target := matches[*m.LoserNextMatchIndex]
		if target.BracketType != BracketLosers {
			t.Errorf("W2 match %d: loser target should be losers bracket, got %q", i, target.BracketType)
		}
		if m.LoserNextMatchSlot == nil || *m.LoserNextMatchSlot != SlotPlayer2 {
			t.Errorf("W2 match %d: loser slot should be player2, got %v", i, m.LoserNextMatchSlot)
		}
	}
}

func TestDoubleElim_LB_InternalLinking(t *testing.T) {
	ids := makePlayerIDsUnsorted(8)
	matches, err := GenerateDoubleElimination(ids, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Collect losers bracket matches in order
	var lbMatches []int
	for i, m := range matches {
		if m.BracketType == BracketLosers {
			lbMatches = append(lbMatches, i)
		}
	}

	// For 8 players: L1(2), L2(2), L3(1), L4(1) = 6 LB matches
	if len(lbMatches) != 6 {
		t.Fatalf("expected 6 LB matches, got %d", len(lbMatches))
	}

	// L1→L2: same count (2→2), 1-to-1, player1 slot
	for i := 0; i < 2; i++ {
		m := matches[lbMatches[i]]
		if m.NextMatchIndex == nil {
			t.Errorf("L1 match %d: expected next link", i)
			continue
		}
		if m.NextMatchSlot == nil || *m.NextMatchSlot != SlotPlayer1 {
			t.Errorf("L1 match %d: same-count link should be player1, got %v", i, m.NextMatchSlot)
		}
	}

	// L2→L3: halved (2→1), 2-to-1
	for i := 2; i < 4; i++ {
		m := matches[lbMatches[i]]
		if m.NextMatchIndex == nil {
			t.Errorf("L2 match %d: expected next link", i)
			continue
		}
		// Both should point to same L3 match
		if *m.NextMatchIndex != lbMatches[4] {
			t.Errorf("L2 match %d: expected link to L3 (index %d), got %d", i, lbMatches[4], *m.NextMatchIndex)
		}
	}

	// L3→L4: same count (1→1), 1-to-1, player1
	m := matches[lbMatches[4]]
	if m.NextMatchIndex == nil {
		t.Errorf("L3: expected next link")
	} else if *m.NextMatchIndex != lbMatches[5] {
		t.Errorf("L3: expected link to L4 (index %d), got %d", lbMatches[5], *m.NextMatchIndex)
	}
	if m.NextMatchSlot == nil || *m.NextMatchSlot != SlotPlayer1 {
		t.Errorf("L3: same-count link should be player1, got %v", m.NextMatchSlot)
	}
}

func TestDoubleElim_GF_Linking(t *testing.T) {
	ids := makePlayerIDsUnsorted(8)
	matches, err := GenerateDoubleElimination(ids, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Find GF match
	var gfIdx int
	for i, m := range matches {
		if m.BracketType == BracketGrandFinal {
			gfIdx = i
			break
		}
	}

	// WB final should link to GF as player1
	// WB final is the last winners match (index 6 for 8 players: 4+2+1=7 matches, index 6)
	wbFinalIdx := -1
	for i := len(matches) - 1; i >= 0; i-- {
		if matches[i].BracketType == BracketWinners {
			wbFinalIdx = i
			break
		}
	}

	wbFinal := matches[wbFinalIdx]
	if wbFinal.NextMatchIndex == nil || *wbFinal.NextMatchIndex != gfIdx {
		t.Errorf("WB final should link to GF (index %d), got %v", gfIdx, wbFinal.NextMatchIndex)
	}
	if wbFinal.NextMatchSlot == nil || *wbFinal.NextMatchSlot != SlotPlayer1 {
		t.Errorf("WB final slot should be player1, got %v", wbFinal.NextMatchSlot)
	}

	// LB final should link to GF as player2
	lbFinalIdx := gfIdx - 1
	lbFinal := matches[lbFinalIdx]
	if lbFinal.BracketType != BracketLosers {
		t.Fatalf("expected LB final at index %d, got bracket_type=%q", lbFinalIdx, lbFinal.BracketType)
	}
	if lbFinal.NextMatchIndex == nil || *lbFinal.NextMatchIndex != gfIdx {
		t.Errorf("LB final should link to GF (index %d), got %v", gfIdx, lbFinal.NextMatchIndex)
	}
	if lbFinal.NextMatchSlot == nil || *lbFinal.NextMatchSlot != SlotPlayer2 {
		t.Errorf("LB final slot should be player2, got %v", lbFinal.NextMatchSlot)
	}

	// GF match_id should be "GF-M1"
	gf := matches[gfIdx]
	if gf.MatchID != "GF-M1" {
		t.Errorf("GF match_id: expected 'GF-M1', got %q", gf.MatchID)
	}
	if gf.RoundName == nil || *gf.RoundName != "Grand Final" {
		t.Errorf("GF round name: expected 'Grand Final', got %v", gf.RoundName)
	}
}

func TestDoubleElim_BestOf(t *testing.T) {
	ids := makePlayerIDsUnsorted(4)
	bestOf := &BestOfConfig{
		Default:         intPtr(3),
		Final:           intPtr(5),
		GrandFinal:      intPtr(7),
		GrandFinalReset: intPtr(7),
	}
	opts := &DoubleEliminationOptions{
		GrandFinalReset: true,
		BestOf:          bestOf,
	}
	matches, err := GenerateDoubleElimination(ids, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, m := range matches {
		switch m.BracketType {
		case BracketGrandFinal:
			if m.BestOf == nil || *m.BestOf != 7 {
				t.Errorf("GF best_of: expected 7, got %v", m.BestOf)
			}
		case BracketGrandFinalReset:
			if m.BestOf == nil || *m.BestOf != 7 {
				t.Errorf("GF Reset best_of: expected 7, got %v", m.BestOf)
			}
		case BracketWinners:
			if m.RoundName != nil && *m.RoundName == "Final" {
				if m.BestOf == nil || *m.BestOf != 5 {
					t.Errorf("WB Final best_of: expected 5, got %v", m.BestOf)
				}
			} else {
				if m.BestOf == nil || *m.BestOf != 3 {
					t.Errorf("WB non-final best_of: expected 3, got %v", m.BestOf)
				}
			}
		case BracketLosers:
			if m.BestOf == nil || *m.BestOf != 3 {
				t.Errorf("LB best_of: expected 3, got %v", m.BestOf)
			}
		}
	}
}

func TestDoubleElim_LessThan2_ReturnsNil(t *testing.T) {
	matches, err := GenerateDoubleElimination([]string{"p1"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if matches != nil {
		t.Errorf("expected nil for <2 participants, got %v", matches)
	}
}

func TestDoubleElim_MatchNumbers_Sequential(t *testing.T) {
	ids := makePlayerIDsUnsorted(8)
	matches, err := GenerateDoubleElimination(ids, nil)
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

func TestDoubleElim_16Players(t *testing.T) {
	ids := makePlayerIDsUnsorted(16)
	matches, err := GenerateDoubleElimination(ids, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 16 players → bracketSize=16, totalWinnersRounds=4
	// Winners: 8+4+2+1 = 15
	// Losers: totalLosersRounds = (4-1)*2 = 6
	//   L1: 16/4=4, L2: 4, L3: 2 (halved after L2), L4: 2, L5: 1 (halved after L4), L6: 1
	//   Total LB = 4+4+2+2+1+1 = 14
	// GF: 1
	// Total: 15 + 14 + 1 = 30

	winnersCount := countMatchesByBracketType(matches, BracketWinners)
	losersCount := countMatchesByBracketType(matches, BracketLosers)
	gfCount := countMatchesByBracketType(matches, BracketGrandFinal)

	if winnersCount != 15 {
		t.Errorf("expected 15 winners matches, got %d", winnersCount)
	}
	if losersCount != 14 {
		t.Errorf("expected 14 losers matches, got %d", losersCount)
	}
	if gfCount != 1 {
		t.Errorf("expected 1 grand final, got %d", gfCount)
	}
	if len(matches) != 30 {
		t.Errorf("expected 30 total matches, got %d", len(matches))
	}
}

func TestDoubleElim_LB_RoundNames(t *testing.T) {
	ids := makePlayerIDsUnsorted(8)
	matches, err := GenerateDoubleElimination(ids, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Collect LB matches by round
	lbByRound := make(map[int][]MatchSeed)
	for _, m := range matches {
		if m.BracketType == BracketLosers {
			lbByRound[m.Round] = append(lbByRound[m.Round], m)
		}
	}

	// For 8 players: 4 LB rounds
	// L4 (last round) should be "LB Final"
	// L3 should be "LB Semi-final"
	for _, m := range lbByRound[4] {
		if m.RoundName == nil || *m.RoundName != "LB Final" {
			t.Errorf("L4 round name: expected 'LB Final', got %v", m.RoundName)
		}
	}
	for _, m := range lbByRound[3] {
		if m.RoundName == nil || *m.RoundName != "LB Semi-final" {
			t.Errorf("L3 round name: expected 'LB Semi-final', got %v", m.RoundName)
		}
	}
}
