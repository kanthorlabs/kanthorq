package faker

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
