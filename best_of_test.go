package bracket

import "testing"

func TestResolveBestOfSE_NilConfig(t *testing.T) {
	result := resolveBestOfSE(BracketWinners, strPtr("Final"), nil)
	if result != nil {
		t.Errorf("expected nil for nil config, got %v", result)
	}
}

func TestResolveBestOfSE_Default(t *testing.T) {
	cfg := &BestOfConfig{Default: intPtr(3)}
	result := resolveBestOfSE(BracketWinners, strPtr("Semi-final"), cfg)
	if result == nil || *result != 3 {
		t.Errorf("expected 3, got %v", result)
	}
}

func TestResolveBestOfSE_Final(t *testing.T) {
	cfg := &BestOfConfig{Default: intPtr(3), Final: intPtr(5)}
	result := resolveBestOfSE(BracketWinners, strPtr("Final"), cfg)
	if result == nil || *result != 5 {
		t.Errorf("expected 5, got %v", result)
	}
}

func TestResolveBestOfSE_ThirdPlace(t *testing.T) {
	cfg := &BestOfConfig{Default: intPtr(3), ThirdPlace: intPtr(1)}
	result := resolveBestOfSE(BracketThirdPlace, strPtr("3rd Place Match"), cfg)
	if result == nil || *result != 1 {
		t.Errorf("expected 1, got %v", result)
	}
}

func TestResolveBestOfSE_ThirdPlace_FallsBackToDefault(t *testing.T) {
	cfg := &BestOfConfig{Default: intPtr(3)}
	result := resolveBestOfSE(BracketThirdPlace, strPtr("3rd Place Match"), cfg)
	if result == nil || *result != 3 {
		t.Errorf("expected 3 (default), got %v", result)
	}
}

func TestResolveBestOfDE_NilConfig(t *testing.T) {
	result := resolveBestOfDE(BracketGrandFinal, strPtr("Grand Final"), nil)
	if result != nil {
		t.Errorf("expected nil for nil config, got %v", result)
	}
}

func TestResolveBestOfDE_GrandFinal(t *testing.T) {
	cfg := &BestOfConfig{Default: intPtr(3), GrandFinal: intPtr(7)}
	result := resolveBestOfDE(BracketGrandFinal, strPtr("Grand Final"), cfg)
	if result == nil || *result != 7 {
		t.Errorf("expected 7, got %v", result)
	}
}

func TestResolveBestOfDE_GrandFinalReset(t *testing.T) {
	cfg := &BestOfConfig{Default: intPtr(3), GrandFinalReset: intPtr(7)}
	result := resolveBestOfDE(BracketGrandFinalReset, strPtr("Grand Final Reset"), cfg)
	if result == nil || *result != 7 {
		t.Errorf("expected 7, got %v", result)
	}
}

func TestResolveBestOfDE_WinnersFinal(t *testing.T) {
	cfg := &BestOfConfig{Default: intPtr(3), Final: intPtr(5)}
	result := resolveBestOfDE(BracketWinners, strPtr("Final"), cfg)
	if result == nil || *result != 5 {
		t.Errorf("expected 5, got %v", result)
	}
}

func TestResolveBestOfDE_Default(t *testing.T) {
	cfg := &BestOfConfig{Default: intPtr(3)}
	result := resolveBestOfDE(BracketLosers, strPtr("LB Round 1"), cfg)
	if result == nil || *result != 3 {
		t.Errorf("expected 3, got %v", result)
	}
}
