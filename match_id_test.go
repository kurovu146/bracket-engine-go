package bracket

import "testing"

func TestGenerateMatchID(t *testing.T) {
	tests := []struct {
		name         string
		bracketType  BracketType
		round        int
		matchInRound int
		want         string
	}{
		// Winners bracket
		{name: "winners R1 M1", bracketType: BracketWinners, round: 1, matchInRound: 1, want: "WB-R1-M1"},
		{name: "winners R2 M3", bracketType: BracketWinners, round: 2, matchInRound: 3, want: "WB-R2-M3"},
		{name: "winners R3 M1", bracketType: BracketWinners, round: 3, matchInRound: 1, want: "WB-R3-M1"},

		// Losers bracket
		{name: "losers R1 M1", bracketType: BracketLosers, round: 1, matchInRound: 1, want: "LB-R1-M1"},
		{name: "losers R4 M2", bracketType: BracketLosers, round: 4, matchInRound: 2, want: "LB-R4-M2"},

		// Grand final
		{name: "grand final", bracketType: BracketGrandFinal, round: 1, matchInRound: 1, want: "GF-M1"},
		// Grand final ignores round/matchInRound args
		{name: "grand final arbitrary args", bracketType: BracketGrandFinal, round: 5, matchInRound: 9, want: "GF-M1"},

		// Grand final reset
		{name: "grand final reset", bracketType: BracketGrandFinalReset, round: 1, matchInRound: 1, want: "GF-M2"},

		// Third place
		{name: "third place", bracketType: BracketThirdPlace, round: 1, matchInRound: 1, want: "3RD-M1"},

		// Round robin
		{name: "round robin R1 M1", bracketType: BracketRoundRobin, round: 1, matchInRound: 1, want: "RR-R1-M1"},
		{name: "round robin R3 M2", bracketType: BracketRoundRobin, round: 3, matchInRound: 2, want: "RR-R3-M2"},

		// Swiss
		{name: "swiss R1 M1", bracketType: BracketSwiss, round: 1, matchInRound: 1, want: "SW-R1-M1"},
		{name: "swiss R2 M4", bracketType: BracketSwiss, round: 2, matchInRound: 4, want: "SW-R2-M4"},

		// Group stage
		{name: "group 0 R1 M1", bracketType: GroupBracketType(0), round: 1, matchInRound: 1, want: "G0-R1-M1"},
		{name: "group 2 R3 M4", bracketType: GroupBracketType(2), round: 3, matchInRound: 4, want: "G2-R3-M4"},
		{name: "group 5 R1 M2", bracketType: GroupBracketType(5), round: 1, matchInRound: 2, want: "G5-R1-M2"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateMatchID(tt.bracketType, tt.round, tt.matchInRound)
			if got != tt.want {
				t.Errorf("GenerateMatchID(%q, %d, %d) = %q, want %q",
					tt.bracketType, tt.round, tt.matchInRound, got, tt.want)
			}
		})
	}
}
