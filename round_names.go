package bracket

import (
	"fmt"
	"strings"
)

// ResolveRoundName returns a human-readable round name based on bracket type and position.
//
// Elimination brackets use positional names (Final, Semi-final, Quarter-final).
// Round-robin, Swiss, and group stages use "Round N".
func ResolveRoundName(bracketType BracketType, round, totalRounds int) string {
	switch bracketType {
	case "grand_final":
		return "Grand Final"
	case "grand_final_reset":
		return "Grand Final Reset"
	case "third_place":
		return "3rd Place Match"

	case "winners":
		fromFinal := totalRounds - round
		if fromFinal == 0 {
			return "Final"
		}
		if fromFinal == 1 {
			return "Semi-final"
		}
		if fromFinal == 2 {
			return "Quarter-final"
		}
		return fmt.Sprintf("Round %d", round)

	case "losers":
		fromFinal := totalRounds - round
		if fromFinal == 0 {
			return "LB Final"
		}
		if fromFinal == 1 {
			return "LB Semi-final"
		}
		return fmt.Sprintf("LB Round %d", round)

	case "round_robin", "swiss":
		return fmt.Sprintf("Round %d", round)

	default:
		if strings.HasPrefix(string(bracketType), "group_") {
			return fmt.Sprintf("Round %d", round)
		}
		return fmt.Sprintf("Round %d", round)
	}
}
