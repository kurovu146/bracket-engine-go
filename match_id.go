package bracket

import (
	"fmt"
	"strings"
)

// GenerateMatchID produces a stable, human-readable match identifier.
//
// Format by bracket type:
//   - winners:           "WB-R{round}-M{matchInRound}"
//   - losers:            "LB-R{round}-M{matchInRound}"
//   - grand_final:       "GF-M1"
//   - grand_final_reset: "GF-M2"
//   - third_place:       "3RD-M1"
//   - round_robin:       "RR-R{round}-M{matchInRound}"
//   - swiss:             "SW-R{round}-M{matchInRound}"
//   - group_N:           "G{N}-R{round}-M{matchInRound}"
func GenerateMatchID(bracketType BracketType, round, matchInRound int) string {
	switch bracketType {
	case BracketWinners:
		return fmt.Sprintf("WB-R%d-M%d", round, matchInRound)
	case BracketLosers:
		return fmt.Sprintf("LB-R%d-M%d", round, matchInRound)
	case BracketGrandFinal:
		return "GF-M1"
	case BracketGrandFinalReset:
		return "GF-M2"
	case BracketThirdPlace:
		return "3RD-M1"
	case BracketRoundRobin:
		return fmt.Sprintf("RR-R%d-M%d", round, matchInRound)
	case BracketSwiss:
		return fmt.Sprintf("SW-R%d-M%d", round, matchInRound)
	default:
		// group_N pattern
		bt := string(bracketType)
		if strings.HasPrefix(bt, "group_") {
			groupNum := bt[len("group_"):]
			return fmt.Sprintf("G%s-R%d-M%d", groupNum, round, matchInRound)
		}
		return fmt.Sprintf("%s-R%d-M%d", bt, round, matchInRound)
	}
}
