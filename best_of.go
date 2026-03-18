package bracket

// resolveBestOfSE resolves best-of for single elimination matches.
func resolveBestOfSE(bracketType BracketType, roundName *string, cfg *BestOfConfig) *int {
	if cfg == nil {
		return nil
	}
	if bracketType == BracketThirdPlace && cfg.ThirdPlace != nil {
		return cfg.ThirdPlace
	}
	if roundName != nil && *roundName == "Final" && cfg.Final != nil {
		return cfg.Final
	}
	return cfg.Default
}

// resolveBestOfDE resolves best-of for double elimination matches.
func resolveBestOfDE(bracketType BracketType, roundName *string, cfg *BestOfConfig) *int {
	if cfg == nil {
		return nil
	}
	if bracketType == BracketGrandFinal && cfg.GrandFinal != nil {
		return cfg.GrandFinal
	}
	if bracketType == BracketGrandFinalReset && cfg.GrandFinalReset != nil {
		return cfg.GrandFinalReset
	}
	if roundName != nil && *roundName == "Final" && cfg.Final != nil {
		return cfg.Final
	}
	return cfg.Default
}
