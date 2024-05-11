package configuration

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/kanthorlabs/common/safe"
	"github.com/kanthorlabs/common/utils"
	"github.com/spf13/viper"
)

var FileLookingDirs = []string{"$KANTHOR_HOME/", "$HOME/.kanthor/", "./"}
var FileName = "configs"
var FileExt = "yaml"

// NewFile creates a configuration provider that
//   - reads configuraton with name `configs.yaml` from the given directories
//   - reads environment variables with prefix `ns`
//
// the configuration inside the file will be merged and overriden by the environment variables
func NewFile(ns string, dirs []string) (Provider, error) {
	if len(dirs) == 0 {
		return nil, fmt.Errorf("CONFIGURATION.FILE.NO_DIRECTOY.ERROR")
	}

	instance := viper.New()
	instance.SetConfigName(FileName) // name of config file (without extension)
	instance.SetConfigType(FileExt)  // extension

	sources := []Source{}
	for _, dir := range dirs {
		dir = strings.Trim(dir, " ")
		filename := FileName + "." + FileExt

		source := Source{Dir: dir, Looking: path.Join(dir, filename), Found: path.Join(utils.AbsPathify(dir), filename)}
		sources = append(sources, source)
	}

	for i := range sources {
		if _, err := os.Stat(sources[i].Found); err == nil {
			sources[i].Used = true
			instance.AddConfigPath(sources[i].Dir)

			if err := instance.MergeInConfig(); err != nil {
				// ignore not found files, otherwise return error
				if _, notfound := err.(viper.ConfigFileNotFoundError); !notfound {
					return nil, fmt.Errorf("CONFIGURATION.FILE.ERROR: %w", err)
				}
			}

			break
		}
	}

	instance.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))
	instance.SetEnvPrefix(ns)
	instance.AutomaticEnv()

	return &file{viper: instance, sources: sources}, nil
}

type file struct {
	viper   *viper.Viper
	sources []Source
}

func (provider *file) Unmarshal(dest any) error {
	return provider.viper.Unmarshal(
		dest,
		viper.DecodeHook(safe.MetadataMapstructureHook()),
	)
}

func (provider *file) Sources() []Source {
	return provider.sources
}

func (provider *file) SetDefault(key string, value any) {
	provider.viper.SetDefault(key, value)
}

func (provider *file) Set(key string, value any) {
	provider.viper.Set(key, value)
}
