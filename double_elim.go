package bracket

import "math"

// GenerateDoubleElimination generates a double elimination bracket.
//
// Players start in the Winners Bracket. Losing once drops you to the
// Losers Bracket. Losing twice eliminates you. The winners of each bracket
// meet in the Grand Final.
//
// Losers bracket flow:
//   - Odd losers rounds (L1, L3, L5...): losers from a winners round
//     play among themselves (or are paired for L1)
//   - Even losers rounds (L2, L4, L6...): survivors from the previous
//     losers round play against new losers dropping from the next winners round
func GenerateDoubleElimination(participantIDs []string, options *DoubleEliminationOptions) ([]MatchSeed, error) {
	if len(participantIDs) < 2 {
		return nil, nil
	}

	if err := ValidateParticipants(participantIDs, 2); err != nil {
		return nil, err
	}

	n := len(participantIDs)
	bracketSize := nextPowerOf2(n)
	totalWinnersRounds := log2Int(bracketSize)
	seeded := StandardSeed(participantIDs, bracketSize)

	matches := make([]MatchSeed, 0)
	matchNumber := 1

	var bestOfCfg *BestOfConfig
	if options != nil {
		bestOfCfg = options.BestOf
	}

	// ============ WINNERS BRACKET ============
	winnersRoundStart := make([]int, 0, totalWinnersRounds)

	for round := 1; round <= totalWinnersRounds; round++ {
		matchesInRound := bracketSize / int(math.Pow(2, float64(round)))
		winnersRoundStart = append(winnersRoundStart, len(matches))
		roundName := ResolveRoundName(BracketWinners, round, totalWinnersRounds)

		for i := 0; i < matchesInRound; i++ {
			isFirstRound := round == 1
			var p1, p2 *string
			if isFirstRound {
				p1 = seeded[i*2]
				p2 = seeded[i*2+1]
			}
			isBye := isFirstRound && (p1 == nil || p2 == nil)

			matches = append(matches, MatchSeed{
				MatchID:             GenerateMatchID(BracketWinners, round, i+1),
				Round:               round,
				MatchNumber:         matchNumber,
				Player1ID:           p1,
				Player2ID:           p2,
				BracketType:         BracketWinners,
				NextMatchIndex:      nil,
				LoserNextMatchIndex: nil,
				NextMatchSlot:       nil,
				LoserNextMatchSlot:  nil,
				RoundName:           strPtr(roundName),
				IsBye:               isBye,
				BestOf:              resolveBestOfDE(BracketWinners, strPtr(roundName), bestOfCfg),
			})
			matchNumber++
		}
	}

	// Link winners bracket internally
	for r := 0; r < totalWinnersRounds-1; r++ {
		start := winnersRoundStart[r]
		nextStart := winnersRoundStart[r+1]
		count := bracketSize / int(math.Pow(2, float64(r+1)))

		for i := 0; i < count; i++ {
			matches[start+i].NextMatchIndex = intPtr(nextStart + i/2)
			if i%2 == 0 {
				matches[start+i].NextMatchSlot = slotPtr(SlotPlayer1)
			} else {
				matches[start+i].NextMatchSlot = slotPtr(SlotPlayer2)
			}
		}
	}

	winnersMatchCount := len(matches)

	// ============ LOSERS BRACKET ============
	totalLosersRounds := (totalWinnersRounds - 1) * 2
	losersRoundStart := make([]int, 0, totalLosersRounds)
	losersRoundCounts := make([]int, 0, totalLosersRounds)

	// L1 count = bracketSize/4, halves every 2 rounds (on even->odd transition)
	currentCount := bracketSize / 4

	for lRound := 1; lRound <= totalLosersRounds; lRound++ {
		losersRoundStart = append(losersRoundStart, len(matches))
		losersRoundCounts = append(losersRoundCounts, currentCount)
		lRoundName := ResolveRoundName(BracketLosers, lRound, totalLosersRounds)

		for i := 0; i < currentCount; i++ {
			matches = append(matches, MatchSeed{
				MatchID:             GenerateMatchID(BracketLosers, lRound, i+1),
				Round:               lRound,
				MatchNumber:         matchNumber,
				Player1ID:           nil,
				Player2ID:           nil,
				BracketType:         BracketLosers,
				NextMatchIndex:      nil,
				LoserNextMatchIndex: nil,
				NextMatchSlot:       nil,
				LoserNextMatchSlot:  nil,
				RoundName:           strPtr(lRoundName),
				IsBye:               false,
				BestOf:              resolveBestOfDE(BracketLosers, strPtr(lRoundName), bestOfCfg),
			})
			matchNumber++
		}

		// Count halves on even->odd transition (after even rounds)
		if lRound%2 == 0 {
			currentCount = currentCount / 2
			if currentCount < 1 {
				currentCount = 1
			}
		}
	}

	// Link losers bracket internally
	for lr := 0; lr < totalLosersRounds-1; lr++ {
		start := losersRoundStart[lr]
		nextStart := losersRoundStart[lr+1]
		count := losersRoundCounts[lr]
		nextCount := losersRoundCounts[lr+1]

		for i := 0; i < count; i++ {
			if count == nextCount {
				// Same count: 1-to-1 (LB survivor fills player1)
				matches[start+i].NextMatchIndex = intPtr(nextStart + i)
				matches[start+i].NextMatchSlot = slotPtr(SlotPlayer1)
			} else {
				// Halved: 2-to-1 (even pos -> player1, odd -> player2)
				matches[start+i].NextMatchIndex = intPtr(nextStart + i/2)
				if i%2 == 0 {
					matches[start+i].NextMatchSlot = slotPtr(SlotPlayer1)
				} else {
					matches[start+i].NextMatchSlot = slotPtr(SlotPlayer2)
				}
			}
		}
	}

	// ============ LINK WINNERS LOSERS -> LOSERS BRACKET ============
	for wRound := 0; wRound < totalWinnersRounds; wRound++ {
		wStart := winnersRoundStart[wRound]
		wCount := bracketSize / int(math.Pow(2, float64(wRound+1)))

		if wRound == 0 {
			// W1 losers -> L1: pair them 2-to-1
			if len(losersRoundStart) > 0 {
				lStart := losersRoundStart[0]
				for i := 0; i < wCount; i++ {
					matches[wStart+i].LoserNextMatchIndex = intPtr(lStart + i/2)
					if i%2 == 0 {
						matches[wStart+i].LoserNextMatchSlot = slotPtr(SlotPlayer1)
					} else {
						matches[wStart+i].LoserNextMatchSlot = slotPtr(SlotPlayer2)
					}
				}
			}
		} else {
			// W(r+1) losers -> L(r*2) for r=1,2,...
			// lRoundIndex is 0-based: wRound*2 - 1
			lRoundIndex := wRound*2 - 1
			if lRoundIndex < totalLosersRounds {
				lStart := losersRoundStart[lRoundIndex]
				for i := 0; i < wCount; i++ {
					if lStart+i < len(matches) && matches[lStart+i].BracketType == BracketLosers {
						matches[wStart+i].LoserNextMatchIndex = intPtr(lStart + i)
						matches[wStart+i].LoserNextMatchSlot = slotPtr(SlotPlayer2)
					}
				}
			}
		}
	}

	// ============ GRAND FINAL ============
	grandFinalIndex := len(matches)
	gfRoundName := ResolveRoundName(BracketGrandFinal, 1, 1)

	matches = append(matches, MatchSeed{
		MatchID:             GenerateMatchID(BracketGrandFinal, 1, 1),
		Round:               1,
		MatchNumber:         matchNumber,
		Player1ID:           nil,
		Player2ID:           nil,
		BracketType:         BracketGrandFinal,
		NextMatchIndex:      nil,
		LoserNextMatchIndex: nil,
		NextMatchSlot:       nil,
		LoserNextMatchSlot:  nil,
		RoundName:           strPtr(gfRoundName),
		IsBye:               false,
		BestOf:              resolveBestOfDE(BracketGrandFinal, strPtr(gfRoundName), bestOfCfg),
	})
	matchNumber++

	// Link WB Final -> Grand Final as player1
	matches[winnersMatchCount-1].NextMatchIndex = intPtr(grandFinalIndex)
	matches[winnersMatchCount-1].NextMatchSlot = slotPtr(SlotPlayer1)

	if totalLosersRounds == 0 {
		// 2-player case: WB match loser goes directly to GF as player2
		matches[0].LoserNextMatchIndex = intPtr(grandFinalIndex)
		matches[0].LoserNextMatchSlot = slotPtr(SlotPlayer2)
	} else {
		// Link LB Final -> Grand Final as player2
		lbFinalIndex := grandFinalIndex - 1
		if lbFinalIndex >= 0 && matches[lbFinalIndex].BracketType == BracketLosers {
			matches[lbFinalIndex].NextMatchIndex = intPtr(grandFinalIndex)
			matches[lbFinalIndex].NextMatchSlot = slotPtr(SlotPlayer2)
		}
	}

	// ============ GRAND FINAL RESET (optional) ============
	if options != nil && options.GrandFinalReset {
		resetIndex := len(matches)
		resetRoundName := ResolveRoundName(BracketGrandFinalReset, 1, 1)

		// GF winner links to reset match
		matches[grandFinalIndex].NextMatchIndex = intPtr(resetIndex)
		matches[grandFinalIndex].NextMatchSlot = slotPtr(SlotPlayer1)

		matches = append(matches, MatchSeed{
			MatchID:             GenerateMatchID(BracketGrandFinalReset, 1, 1),
			Round:               1,
			MatchNumber:         matchNumber,
			Player1ID:           nil,
			Player2ID:           nil,
			BracketType:         BracketGrandFinalReset,
			NextMatchIndex:      nil,
			LoserNextMatchIndex: nil,
			NextMatchSlot:       nil,
			LoserNextMatchSlot:  nil,
			RoundName:           strPtr(resetRoundName),
			IsBye:               false,
			BestOf:              resolveBestOfDE(BracketGrandFinalReset, strPtr(resetRoundName), bestOfCfg),
		})
		matchNumber++
	}

	return matches, nil
}
