package bracket

import "fmt"

// byeSentinel is an empty string used as a marker for bye slots.
// Matches containing a bye participant are filtered out of the output.
const byeSentinel = ""

// generateRRPass generates one pass of round-robin rounds using the circle method.
// Fix the first participant, rotate the rest each round.
//
// ids must have even length; bye sentinel ("") should already be appended if needed.
// roundOffset is added to the loop index so second pass continues numbering.
// matchNumberStart is the match_number for the first match in this pass.
func generateRRPass(
	ids []string,
	roundsPerPass int,
	totalRoundsAll int,
	bestOf *int,
	roundOffset int,
	matchNumberStart int,
) []MatchSeed {
	totalParticipants := len(ids)
	matchesPerRound := totalParticipants / 2

	var matches []MatchSeed
	matchNumber := matchNumberStart

	fixed := ids[0]
	// Copy rotating slice so mutations don't affect the caller
	rotating := make([]string, len(ids)-1)
	copy(rotating, ids[1:])

	for round := 0; round < roundsPerPass; round++ {
		absoluteRound := round + 1 + roundOffset

		// Build current order: [fixed, rotating...]
		currentOrder := make([]string, totalParticipants)
		currentOrder[0] = fixed
		copy(currentOrder[1:], rotating)

		matchInRound := 0

		for i := 0; i < matchesPerRound; i++ {
			p1 := currentOrder[i]
			p2 := currentOrder[totalParticipants-1-i]

			// Skip bye matches
			if p1 == byeSentinel || p2 == byeSentinel {
				continue
			}

			matchInRound++
			roundName := ResolveRoundName("round_robin", absoluteRound, totalRoundsAll)

			matches = append(matches, MatchSeed{
				MatchID:     GenerateMatchID("round_robin", absoluteRound, matchInRound),
				Round:       absoluteRound,
				MatchNumber: matchNumber,
				Player1ID:   strPtr(p1),
				Player2ID:   strPtr(p2),
				BracketType: "round_robin",
				RoundName:   &roundName,
				IsBye:       false,
				BestOf:      bestOf,
			})
			matchNumber++
		}

		// Rotate: move last element to front
		last := rotating[len(rotating)-1]
		copy(rotating[1:], rotating[:len(rotating)-1])
		rotating[0] = last
	}

	return matches
}

// GenerateRoundRobin generates a round-robin schedule where every participant
// plays every other participant. Uses the circle method (fix first, rotate rest).
//
// Handles odd counts by adding a bye sentinel. Bye matches are filtered out.
// Optional double round-robin plays a second pass with swapped player1/player2.
func GenerateRoundRobin(participantIDs []string, options *RoundRobinOptions) ([]MatchSeed, error) {
	if len(participantIDs) < 2 {
		return nil, nil
	}

	if err := ValidateParticipants(participantIDs, 2); err != nil {
		return nil, err
	}

	var bestOf *int
	isDouble := false
	if options != nil {
		isDouble = options.DoubleRoundRobin
		if options.BestOf != nil {
			bestOf = options.BestOf.Default
		}
	}

	// Copy IDs; add bye sentinel if odd
	ids := make([]string, len(participantIDs))
	copy(ids, participantIDs)
	if len(ids)%2 != 0 {
		ids = append(ids, byeSentinel)
	}

	roundsPerPass := len(ids) - 1
	totalRoundsAll := roundsPerPass
	if isDouble {
		totalRoundsAll = roundsPerPass * 2
	}

	firstPass := generateRRPass(ids, roundsPerPass, totalRoundsAll, bestOf, 0, 1)

	if !isDouble {
		return firstPass, nil
	}

	// Second pass: same schedule but player1/player2 swapped
	secondPassRaw := generateRRPass(ids, roundsPerPass, totalRoundsAll, bestOf, roundsPerPass, len(firstPass)+1)

	secondPass := make([]MatchSeed, len(secondPassRaw))
	for i, m := range secondPassRaw {
		secondPass[i] = m
		secondPass[i].Player1ID = m.Player2ID
		secondPass[i].Player2ID = m.Player1ID
	}

	result := make([]MatchSeed, 0, len(firstPass)+len(secondPass))
	result = append(result, firstPass...)
	result = append(result, secondPass...)

	return result, nil
}

// roundNamePtr returns a pointer to the resolved round name.
func roundNamePtr(bracketType BracketType, round, totalRounds int) *string {
	name := ResolveRoundName(bracketType, round, totalRounds)
	return &name
}

// fmtGroupBracketType returns "group_N" BracketType for a given group index.
func fmtGroupBracketType(groupIndex int) BracketType {
	return BracketType(fmt.Sprintf("group_%d", groupIndex))
}
