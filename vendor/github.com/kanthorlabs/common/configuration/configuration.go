package configuration

type Provider interface {
	Unmarshal(dest any) error
	Sources() []Source
	SetDefault(key string, value any)
	Set(key string, value any)
}

type Source struct {
	Dir     string
	Looking string
	Found   string
	Used    bool
}

// New creates a new configuration provider with the given namespace
// by default, it will create a file configuration provider
func New(ns string) (Provider, error) {
	return NewFile(ns, FileLookingDirs)
}
