package bracket

import "math"

// GenerateSwiss generates a Swiss-system tournament schedule.
//
// Only Round 1 is fully paired (by seed order: ids[0] vs ids[1], etc.).
// Rounds 2+ are placeholder matches with nil players, to be filled by the
// application based on standings after each round.
//
// Default number of rounds: ceil(log2(n)).
func GenerateSwiss(participantIDs []string, options *SwissOptions) ([]MatchSeed, error) {
	if len(participantIDs) < 2 {
		return nil, nil
	}

	if err := ValidateParticipants(participantIDs, 2); err != nil {
		return nil, err
	}

	var numRounds *int
	var bestOfDefault *int

	if options != nil {
		numRounds = options.NumRounds
		if options.BestOf != nil {
			bestOfDefault = options.BestOf.Default
		}
	}

	n := len(participantIDs)

	totalRounds := 0
	if numRounds != nil {
		totalRounds = *numRounds
	} else {
		totalRounds = int(math.Max(1, math.Ceil(math.Log2(float64(n)))))
	}

	matchesPerRound := n / 2

	var matches []MatchSeed
	matchNumber := 1

	// Round 1: pair by seed order (1v2, 3v4, 5v6, ...)
	for i := 0; i < matchesPerRound; i++ {
		matchInRound := i + 1
		p1 := participantIDs[i*2]
		p2 := participantIDs[i*2+1]

		matches = append(matches, MatchSeed{
			MatchID:     GenerateMatchID("swiss", 1, matchInRound),
			Round:       1,
			MatchNumber: matchNumber,
			Player1ID:   &p1,
			Player2ID:   &p2,
			BracketType: "swiss",
			RoundName:   roundNamePtr("swiss", 1, totalRounds),
			IsBye:       false,
			BestOf:      bestOfDefault,
		})
		matchNumber++
	}

	// Rounds 2+: placeholder matches (app fills players based on standings)
	for round := 2; round <= totalRounds; round++ {
		for i := 0; i < matchesPerRound; i++ {
			matchInRound := i + 1

			matches = append(matches, MatchSeed{
				MatchID:     GenerateMatchID("swiss", round, matchInRound),
				Round:       round,
				MatchNumber: matchNumber,
				Player1ID:   nil,
				Player2ID:   nil,
				BracketType: "swiss",
				RoundName:   roundNamePtr("swiss", round, totalRounds),
				IsBye:       false,
				BestOf:      bestOfDefault,
			})
			matchNumber++
		}
	}

	return matches, nil
}
