package project

import (
	"os"
)

var (
	DefaultRegion    = "ap-southeast-1"
	DefaultNamespace = "kanthor"
	DefaultTier      = "default"
)

func Region() string {
	region := os.Getenv("KANTHOR_REGION")
	if region != "" {
		return region
	}
	return DefaultRegion
}

func Namespace() string {
	ns := os.Getenv("KANTHOR_NAMESPACE")
	if ns != "" {
		return ns
	}
	return DefaultNamespace
}

func Tier() string {
	tier := os.Getenv("KANTHOR_TIER")
	if tier != "" {
		return tier
	}
	return DefaultTier
}
