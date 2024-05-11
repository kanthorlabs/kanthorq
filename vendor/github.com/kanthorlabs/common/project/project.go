package project

import (
	"fmt"
	"os"
	"strings"
)

var (
	version string
	DevEnv  string = "development"
	ProdEnv string = "production"
)

func SetVersion(v string) {
	version = v
}

func GetVersion() string {
	return version
}

func IsDev() bool {
	return strings.EqualFold(Env(), DevEnv)
}

func Env() string {
	if env := os.Getenv("KANTHOR_ENV"); env != "" {
		return env
	}
	return ProdEnv
}

func Name(name string) string {
	return fmt.Sprintf("%s_%s", Namespace(), name)
}

func Topic(segments ...string) string {
	parts := []string{}
	for i := range segments {
		if segments[i] != "" {
			parts = append(parts, segments[i])
		}
	}
	return strings.Join(parts, ".")
}

func IsTopic(subject, topic string) bool {
	return strings.HasPrefix(subject, Subject(topic))
}

func Subject(segments ...string) string {
	segments = append([]string{Namespace(), Region(), Tier()}, segments...)
	return Topic(segments...)
}
