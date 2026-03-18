package bracket

import "math"

// GenerateGroupStage divides participants into groups and generates round-robin
// matches within each group.
//
// Distribution methods:
//   - "sequential" (default): participants assigned 0,1,2,...,0,1,2,...
//   - "snake": alternating direction each cycle (0,1,2,3,3,2,1,0,...)
//
// Groups with fewer than 2 members are preserved in the output (index stability)
// but produce 0 matches.
func GenerateGroupStage(participantIDs []string, options *GroupStageOptions) (*GroupStageResult, error) {
	if len(participantIDs) < 2 {
		return &GroupStageResult{
			Groups:  [][]string{},
			Matches: []MatchSeed{},
		}, nil
	}

	if err := ValidateParticipants(participantIDs, 2); err != nil {
		return nil, err
	}

	n := len(participantIDs)

	// Resolve options
	var groupCount int
	distribution := "sequential"
	var doubleRR bool
	var bestOfCfg *BestOfConfig

	if options != nil {
		if options.NumGroups != nil {
			groupCount = *options.NumGroups
		}
		if options.Distribution != "" {
			distribution = options.Distribution
		}
		doubleRR = options.DoubleRoundRobin
		bestOfCfg = options.BestOf
	}

	if groupCount == 0 {
		// Auto-calculate: aim for ~4 players per group, minimum 2 groups
		groupCount = int(math.Max(2, math.Round(float64(n)/4.0)))
	}

	// Distribute participants into groups
	groups := make([][]string, groupCount)
	for i := range groups {
		groups[i] = []string{}
	}

	if distribution == "snake" {
		for i := 0; i < n; i++ {
			cycle := i / groupCount
			pos := i % groupCount
			groupIndex := pos
			if cycle%2 != 0 {
				groupIndex = groupCount - 1 - pos
			}
			groups[groupIndex] = append(groups[groupIndex], participantIDs[i])
		}
	} else {
		// Sequential (default)
		for i := 0; i < n; i++ {
			groups[i%groupCount] = append(groups[i%groupCount], participantIDs[i])
		}
	}

	// Build round-robin options for each group
	rrOpts := &RoundRobinOptions{
		DoubleRoundRobin: doubleRR,
		BestOf:           bestOfCfg,
	}

	var allMatches []MatchSeed
	matchNumber := 1

	for gi := 0; gi < len(groups); gi++ {
		// Skip groups with < 2 players but preserve index in groups array
		if len(groups[gi]) < 2 {
			continue
		}

		bracketType := fmtGroupBracketType(gi)

		groupMatches, err := GenerateRoundRobin(groups[gi], rrOpts)
		if err != nil {
			return nil, err
		}

		// Track per-round match counter for stable match IDs
		matchInRound := make(map[int]int)

		for _, m := range groupMatches {
			matchInRound[m.Round]++

			allMatches = append(allMatches, MatchSeed{
				MatchID:             GenerateMatchID(bracketType, m.Round, matchInRound[m.Round]),
				Round:               m.Round,
				MatchNumber:         matchNumber,
				Player1ID:           m.Player1ID,
				Player2ID:           m.Player2ID,
				BracketType:         bracketType,
				NextMatchIndex:      nil,
				LoserNextMatchIndex: nil,
				NextMatchSlot:       nil,
				LoserNextMatchSlot:  nil,
				RoundName:           roundNamePtr(bracketType, m.Round, m.Round),
				IsBye:               m.IsBye,
				BestOf:              m.BestOf,
			})
			matchNumber++
		}
	}

	return &GroupStageResult{
		Groups:  groups,
		Matches: allMatches,
	}, nil
}
