package xfaker

import "strings"

func StreamName() string {
	words := F.Lorem().Words(8)
	return strings.ToLower(strings.Join(words, "_"))
}

func ConsumerName() string {
	words := F.Lorem().Words(8)
	return strings.ToLower(strings.Join(words, "_"))
}

func Subject() string {
	words := F.Lorem().Words(8)
	return strings.ToLower(strings.Join(words, "."))
}

func SubjectWihtPattern(pattern string) string {
	patternTokens := strings.Split(pattern, ".")
	for i := 0; i < len(patternTokens); i++ {
		if patternTokens[i] == "*" {
			patternTokens[i] = F.Lorem().Word()
		}
	}

	if patternTokens[len(patternTokens)-1] == ">" {
		lastTokens := F.Lorem().Words(F.IntBetween(1, 5))
		patternTokens = append(patternTokens[:len(patternTokens)-1], lastTokens...)
	}

	return strings.ToLower(strings.Join(patternTokens, "."))
}
