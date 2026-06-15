package symbols

import "regexp"

var bsePattern = regexp.MustCompile(`^[48]`)

// InferAShareExchange infers an A-share exchange from a plain numeric code.
func InferAShareExchange(code string) Exchange {
	if len(code) >= 2 && code[:2] == "92" || bsePattern.MatchString(code) {
		return ExchangeBSE
	}
	if code == "" {
		return ExchangeSZSE
	}
	switch code[0] {
	case '6', '5', '9':
		return ExchangeSSE
	default:
		return ExchangeSZSE
	}
}
