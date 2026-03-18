package bracket

import "fmt"

// BracketType represents the type of bracket a match belongs to.
type BracketType string

const (
	BracketWinners         BracketType = "winners"
	BracketLosers          BracketType = "losers"
	BracketGrandFinal      BracketType = "grand_final"
	BracketGrandFinalReset BracketType = "grand_final_reset"
	BracketThirdPlace      BracketType = "third_place"
	BracketRoundRobin      BracketType = "round_robin"
	BracketSwiss           BracketType = "swiss"
)

// GroupBracketType returns "group_N" BracketType for group index N.
func GroupBracketType(n int) BracketType {
	return BracketType(fmt.Sprintf("group_%d", n))
}

// MatchSlot indicates which slot the winner/loser fills in the next match.
type MatchSlot string

const (
	SlotPlayer1 MatchSlot = "player1"
	SlotPlayer2 MatchSlot = "player2"
)

// BestOfConfig configures best-of counts for different match types.
type BestOfConfig struct {
	Default         *int // Default best-of for all matches
	Final           *int // Best-of for the final match
	ThirdPlace      *int // Best-of for 3rd place match
	GrandFinal      *int // Best-of for grand final
	GrandFinalReset *int // Best-of for grand final reset
}

// MatchSeed represents a generated match in the bracket.
type MatchSeed struct {
	MatchID             string      `json:"match_id"`
	Round               int         `json:"round"`
	MatchNumber         int         `json:"match_number"`
	Player1ID           *string     `json:"player1_id"`
	Player2ID           *string     `json:"player2_id"`
	BracketType         BracketType `json:"bracket_type"`
	NextMatchIndex      *int        `json:"next_match_index"`
	LoserNextMatchIndex *int        `json:"loser_next_match_index"`
	NextMatchSlot       *MatchSlot  `json:"next_match_slot"`
	LoserNextMatchSlot  *MatchSlot  `json:"loser_next_match_slot"`
	RoundName           *string     `json:"round_name"`
	IsBye               bool        `json:"is_bye"`
	BestOf              *int        `json:"best_of"`
}

// SingleEliminationOptions configures single elimination bracket generation.
type SingleEliminationOptions struct {
	ThirdPlaceMatch bool
	BestOf          *BestOfConfig
}

// DoubleEliminationOptions configures double elimination bracket generation.
type DoubleEliminationOptions struct {
	GrandFinalReset bool
	BestOf          *BestOfConfig
}

// RoundRobinOptions configures round-robin schedule generation.
type RoundRobinOptions struct {
	DoubleRoundRobin bool
	BestOf           *BestOfConfig
}

// SwissOptions configures Swiss-system tournament generation.
type SwissOptions struct {
	NumRounds *int
	BestOf    *BestOfConfig
}

// GroupStageOptions configures group stage tournament generation.
type GroupStageOptions struct {
	NumGroups        *int
	Distribution     string // "sequential" (default) or "snake"
	DoubleRoundRobin bool
	BestOf           *BestOfConfig
}

// GroupStageResult contains group assignments and generated matches.
type GroupStageResult struct {
	Groups  [][]string  `json:"groups"`
	Matches []MatchSeed `json:"matches"`
}

func strPtr(s string) *string    { return &s }
func intPtr(i int) *int          { return &i }
func slotPtr(s MatchSlot) *MatchSlot { return &s }
