package bracket

// copyIntPtr returns a copy of the int value behind the pointer (avoids shared aliasing).
func copyIntPtr(p *int) *int {
	if p == nil {
		return nil
	}
	v := *p
	return &v
}

// resolveBestOfSE resolves best-of for single elimination matches.
func resolveBestOfSE(bracketType BracketType, roundName *string, cfg *BestOfConfig) *int {
	if cfg == nil {
		return nil
	}
	if bracketType == BracketThirdPlace && cfg.ThirdPlace != nil {
		return copyIntPtr(cfg.ThirdPlace)
	}
	if roundName != nil && *roundName == "Final" && cfg.Final != nil {
		return copyIntPtr(cfg.Final)
	}
	return copyIntPtr(cfg.Default)
}

// resolveBestOfDE resolves best-of for double elimination matches.
func resolveBestOfDE(bracketType BracketType, roundName *string, cfg *BestOfConfig) *int {
	if cfg == nil {
		return nil
	}
	if bracketType == BracketGrandFinal && cfg.GrandFinal != nil {
		return copyIntPtr(cfg.GrandFinal)
	}
	if bracketType == BracketGrandFinalReset && cfg.GrandFinalReset != nil {
		return copyIntPtr(cfg.GrandFinalReset)
	}
	if roundName != nil && *roundName == "Final" && cfg.Final != nil {
		return copyIntPtr(cfg.Final)
	}
	return copyIntPtr(cfg.Default)
}
