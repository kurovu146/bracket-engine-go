package bracket

import "testing"

func TestResolveRoundName_Winners(t *testing.T) {
	tests := []struct {
		name        string
		round       int
		totalRounds int
		want        string
	}{
		{name: "Final", round: 3, totalRounds: 3, want: "Final"},
		{name: "Semi-final", round: 2, totalRounds: 3, want: "Semi-final"},
		{name: "Quarter-final", round: 1, totalRounds: 3, want: "Quarter-final"},
		{name: "Round 1 of 4", round: 1, totalRounds: 4, want: "Round 1"},
		{name: "Final of 1", round: 1, totalRounds: 1, want: "Final"},
		{name: "Semi-final of 2", round: 1, totalRounds: 2, want: "Semi-final"},
		{name: "Final of 2", round: 2, totalRounds: 2, want: "Final"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ResolveRoundName(BracketWinners, tt.round, tt.totalRounds)
			if got != tt.want {
				t.Errorf("ResolveRoundName(winners, %d, %d) = %q, want %q", tt.round, tt.totalRounds, got, tt.want)
			}
		})
	}
}

func TestResolveRoundName_Losers(t *testing.T) {
	tests := []struct {
		name        string
		round       int
		totalRounds int
		want        string
	}{
		{name: "LB Final", round: 4, totalRounds: 4, want: "LB Final"},
		{name: "LB Semi-final", round: 3, totalRounds: 4, want: "LB Semi-final"},
		{name: "LB Round 1", round: 1, totalRounds: 4, want: "LB Round 1"},
		{name: "LB Round 2", round: 2, totalRounds: 4, want: "LB Round 2"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ResolveRoundName(BracketLosers, tt.round, tt.totalRounds)
			if got != tt.want {
				t.Errorf("ResolveRoundName(losers, %d, %d) = %q, want %q", tt.round, tt.totalRounds, got, tt.want)
			}
		})
	}
}

func TestResolveRoundName_SpecialTypes(t *testing.T) {
	tests := []struct {
		name       string
		bracket    BracketType
		want       string
	}{
		{name: "Grand Final", bracket: BracketGrandFinal, want: "Grand Final"},
		{name: "Grand Final Reset", bracket: BracketGrandFinalReset, want: "Grand Final Reset"},
		{name: "3rd Place Match", bracket: BracketThirdPlace, want: "3rd Place Match"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ResolveRoundName(tt.bracket, 1, 1)
			if got != tt.want {
				t.Errorf("ResolveRoundName(%q, 1, 1) = %q, want %q", tt.bracket, got, tt.want)
			}
		})
	}
}

func TestResolveRoundName_RoundRobin(t *testing.T) {
	got := ResolveRoundName(BracketRoundRobin, 3, 5)
	if got != "Round 3" {
		t.Errorf("expected 'Round 3', got %q", got)
	}
}

func TestResolveRoundName_Swiss(t *testing.T) {
	got := ResolveRoundName(BracketSwiss, 2, 4)
	if got != "Round 2" {
		t.Errorf("expected 'Round 2', got %q", got)
	}
}

func TestResolveRoundName_Group(t *testing.T) {
	got := ResolveRoundName(GroupBracketType(1), 2, 3)
	if got != "Round 2" {
		t.Errorf("expected 'Round 2', got %q", got)
	}
}
