package bracket

// GenerateSeedOrder generates standard tournament seed positions for a given bracket size.
// For bracketSize=8: [1, 8, 4, 5, 2, 7, 3, 6]
func GenerateSeedOrder(bracketSize int) []int {
	seeds := []int{1, 2}

	for len(seeds) < bracketSize {
		roundSize := len(seeds) * 2
		expanded := make([]int, 0, roundSize)
		for _, seed := range seeds {
			expanded = append(expanded, seed, roundSize+1-seed)
		}
		seeds = expanded
	}

	return seeds
}

// StandardSeed seeds participants into bracket positions using standard tournament seeding.
// Top seeds get byes for non-power-of-2 participant counts.
// Returns array of participant IDs (or nil for byes) in bracket position order.
func StandardSeed(ids []string, bracketSize int) []*string {
	seedOrder := GenerateSeedOrder(bracketSize)
	result := make([]*string, bracketSize)

	for i := 0; i < bracketSize; i++ {
		seedNumber := seedOrder[i] // 1-based
		if seedNumber <= len(ids) {
			s := ids[seedNumber-1]
			result[i] = &s
		}
		// else: remains nil (bye)
	}

	return result
}
