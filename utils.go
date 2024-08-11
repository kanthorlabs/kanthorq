package kanthorq

import "strings"

func MatchSubject(pattern, subject string) bool {
	if len(subject) == 0 || len(pattern) == 0 {
		return false
	}

	patternTokens := strings.Split(pattern, ".")
	subjectTokens := strings.Split(subject, ".")

	pIdx, sIdx := 0, 0

	for pIdx < len(patternTokens) && sIdx < len(subjectTokens) {
		p := patternTokens[pIdx]

		// '>' matches the rest of the tokens only when it's the last token
		if p == ">" {
			return pIdx == len(patternTokens)-1
		}

		// '*' matches exactly one token
		if p == "*" {
			pIdx++
			sIdx++
			continue
		}

		// Literal token should match exactly
		if p != subjectTokens[sIdx] {
			return false
		}

		// Both literal tokens match, move to the next one
		sIdx++
		pIdx++
	}

	// If we've processed all pattern tokens, check if the subject tokens are also fully consumed
	return pIdx == len(patternTokens) && sIdx == len(subjectTokens)
}
