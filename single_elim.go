package bracket

import "math"

// nextPowerOf2 returns the smallest power of 2 >= n.
func nextPowerOf2(n int) int {
	if n <= 1 {
		return 1
	}
	return 1 << int(math.Ceil(math.Log2(float64(n))))
}

// log2Int returns the integer log base 2 of n.
func log2Int(n int) int {
	return int(math.Log2(float64(n)))
}

// GenerateSingleElimination generates a single elimination (knockout) bracket.
// Supports byes for non-power-of-2 participant counts.
// Top seeds receive byes per standard tournament seeding.
func GenerateSingleElimination(participantIDs []string, options *SingleEliminationOptions) ([]MatchSeed, error) {
	// n < 2 returns empty slice without error (backward compat)
	if len(participantIDs) < 2 {
		return nil, nil
	}

	if err := ValidateParticipants(participantIDs, 2); err != nil {
		return nil, err
	}

	n := len(participantIDs)
	bracketSize := nextPowerOf2(n)
	totalRounds := log2Int(bracketSize)

	seeded := StandardSeed(participantIDs, bracketSize)
	matches := make([]MatchSeed, 0)
	matchNumber := 1

	// --- Round 1 ---
	firstRoundCount := bracketSize / 2

	for i := 0; i < firstRoundCount; i++ {
		p1 := seeded[i*2]
		p2 := seeded[i*2+1]
		roundName := ResolveRoundName(BracketWinners, 1, totalRounds)
		matchInRound := i + 1

		var bestOfCfg *BestOfConfig
		if options != nil {
			bestOfCfg = options.BestOf
		}

		matches = append(matches, MatchSeed{
			MatchID:             GenerateMatchID(BracketWinners, 1, matchInRound),
			Round:               1,
			MatchNumber:         matchNumber,
			Player1ID:           p1,
			Player2ID:           p2,
			BracketType:         BracketWinners,
			NextMatchIndex:      nil,
			LoserNextMatchIndex: nil,
			NextMatchSlot:       nil,
			LoserNextMatchSlot:  nil,
			RoundName:           strPtr(roundName),
			IsBye:               p1 == nil || p2 == nil,
			BestOf:              resolveBestOfSE(BracketWinners, strPtr(roundName), bestOfCfg),
		})
		matchNumber++
	}

	// --- Subsequent rounds ---
	for round := 2; round <= totalRounds; round++ {
		matchesInRound := bracketSize / int(math.Pow(2, float64(round)))
		roundName := ResolveRoundName(BracketWinners, round, totalRounds)

		var bestOfCfg *BestOfConfig
		if options != nil {
			bestOfCfg = options.BestOf
		}

		for i := 0; i < matchesInRound; i++ {
			matchInRound := i + 1

			matches = append(matches, MatchSeed{
				MatchID:             GenerateMatchID(BracketWinners, round, matchInRound),
				Round:               round,
				MatchNumber:         matchNumber,
				Player1ID:           nil,
				Player2ID:           nil,
				BracketType:         BracketWinners,
				NextMatchIndex:      nil,
				LoserNextMatchIndex: nil,
				NextMatchSlot:       nil,
				LoserNextMatchSlot:  nil,
				RoundName:           strPtr(roundName),
				IsBye:               false,
				BestOf:              resolveBestOfSE(BracketWinners, strPtr(roundName), bestOfCfg),
			})
			matchNumber++
		}
	}

	// --- Link winners to next round ---
	prevRoundStart := 0
	for round := 1; round < totalRounds; round++ {
		matchesInRound := bracketSize / int(math.Pow(2, float64(round)))
		nextRoundStart := prevRoundStart + matchesInRound

		for i := 0; i < matchesInRound; i++ {
			nextMatchIndex := nextRoundStart + i/2
			var slot MatchSlot
			if i%2 == 0 {
				slot = SlotPlayer1
			} else {
				slot = SlotPlayer2
			}

			matches[prevRoundStart+i].NextMatchIndex = intPtr(nextMatchIndex)
			matches[prevRoundStart+i].NextMatchSlot = slotPtr(slot)
		}

		prevRoundStart = nextRoundStart
	}

	// --- 3rd place match ---
	if options != nil && options.ThirdPlaceMatch && totalRounds >= 2 {
		semiRound := totalRounds - 1

		// Find the two semi-final matches
		semiIndices := make([]int, 0, 2)
		for i := 0; i < len(matches); i++ {
			if matches[i].Round == semiRound && matches[i].BracketType == BracketWinners {
				semiIndices = append(semiIndices, i)
			}
		}

		thirdPlaceIndex := len(matches)
		thirdPlaceRoundName := ResolveRoundName(BracketThirdPlace, totalRounds, totalRounds)

		matches = append(matches, MatchSeed{
			MatchID:             GenerateMatchID(BracketThirdPlace, totalRounds, 1),
			Round:               totalRounds,
			MatchNumber:         matchNumber,
			Player1ID:           nil,
			Player2ID:           nil,
			BracketType:         BracketThirdPlace,
			NextMatchIndex:      nil,
			LoserNextMatchIndex: nil,
			NextMatchSlot:       nil,
			LoserNextMatchSlot:  nil,
			RoundName:           strPtr(thirdPlaceRoundName),
			IsBye:               false,
			BestOf:              resolveBestOfSE(BracketThirdPlace, strPtr(thirdPlaceRoundName), options.BestOf),
		})
		matchNumber++

		// Link semi-final losers to the 3rd place match
		for si, idx := range semiIndices {
			matches[idx].LoserNextMatchIndex = intPtr(thirdPlaceIndex)
			if si == 0 {
				matches[idx].LoserNextMatchSlot = slotPtr(SlotPlayer1)
			} else {
				matches[idx].LoserNextMatchSlot = slotPtr(SlotPlayer2)
			}
		}
	}

	return matches, nil
}
