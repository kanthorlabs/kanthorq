package kanthorq

import "strings"

func MatchSubject(filter, subject string) bool {
	if len(subject) == 0 || len(filter) == 0 {
		return false
	}

	filterTokens := strings.Split(filter, ".")
	subjectTokens := strings.Split(subject, ".")

	pIdx, sIdx := 0, 0

	for pIdx < len(filterTokens) && sIdx < len(subjectTokens) {
		p := filterTokens[pIdx]

		// '>' matches the rest of the tokens only when it's the last token
		if p == ">" {
			return pIdx == len(filterTokens)-1
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

	// If we've processed all filter tokens, check if the subject tokens are also fully consumed
	return pIdx == len(filterTokens) && sIdx == len(subjectTokens)
}
