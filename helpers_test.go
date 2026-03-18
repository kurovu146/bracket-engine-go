package bracket

import (
	"fmt"
	"sort"
	"testing"
)

// sortedPair returns a canonical key for a pair of player IDs (alphabetically sorted).
func sortedPair(a, b string) string {
	if a > b {
		a, b = b, a
	}
	return fmt.Sprintf("%s:%s", a, b)
}

// assertAllPairsPlayOnce verifies every pair of players plays exactly once.
func assertAllPairsPlayOnce(t *testing.T, matches []MatchSeed, playerIDs []string) {
	t.Helper()
	pairs := make(map[string]int)
	for _, m := range matches {
		if m.Player1ID != nil && m.Player2ID != nil {
			key := sortedPair(*m.Player1ID, *m.Player2ID)
			pairs[key]++
		}
	}
	expectedPairs := len(playerIDs) * (len(playerIDs) - 1) / 2
	if len(pairs) != expectedPairs {
		t.Errorf("expected %d unique pairs, got %d", expectedPairs, len(pairs))
	}
	for key, count := range pairs {
		if count != 1 {
			t.Errorf("pair %s appeared %d times, expected 1", key, count)
		}
	}
}

// assertAllPairsPlayTwice verifies every pair plays exactly twice (double round-robin).
func assertAllPairsPlayTwice(t *testing.T, matches []MatchSeed, playerIDs []string) {
	t.Helper()
	pairs := make(map[string]int)
	for _, m := range matches {
		if m.Player1ID != nil && m.Player2ID != nil {
			key := sortedPair(*m.Player1ID, *m.Player2ID)
			pairs[key]++
		}
	}
	expectedPairs := len(playerIDs) * (len(playerIDs) - 1) / 2
	if len(pairs) != expectedPairs {
		t.Errorf("expected %d unique pairs, got %d", expectedPairs, len(pairs))
	}
	for key, count := range pairs {
		if count != 2 {
			t.Errorf("pair %s appeared %d times, expected 2", key, count)
		}
	}
}

// countMatchesByBracketType counts matches for a given bracket type.
func countMatchesByBracketType(matches []MatchSeed, bt BracketType) int {
	count := 0
	for _, m := range matches {
		if m.BracketType == bt {
			count++
		}
	}
	return count
}

// countMatchesByRound counts matches for a given round.
func countMatchesByRound(matches []MatchSeed, round int) int {
	count := 0
	for _, m := range matches {
		if m.Round == round {
			count++
		}
	}
	return count
}

// makePlayerIDs generates a slice of player ID strings: ["p1", "p2", ..., "pN"].
func makePlayerIDs(n int) []string {
	ids := make([]string, n)
	for i := 0; i < n; i++ {
		ids[i] = fmt.Sprintf("p%d", i+1)
	}
	sort.Strings(ids) // Keep them sorted for deterministic tests
	return ids
}

// makePlayerIDsUnsorted generates player IDs without sorting (preserves p1,p2,...,pN order).
func makePlayerIDsUnsorted(n int) []string {
	ids := make([]string, n)
	for i := 0; i < n; i++ {
		ids[i] = fmt.Sprintf("p%d", i+1)
	}
	return ids
}
