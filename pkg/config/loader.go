package config

import (
	"os"

	"github.com/ziollek/etherscription/pkg/logging"
	yaml "gopkg.in/yaml.v3"
)

const configPath = "./configuration/config.yaml"

func LoadConfig(path string) (*Config, error) {
	var appConfig Config
	// Load configuration from file
	yamlData, err := os.ReadFile(path)
	// let's use the default path if a provided path does not exist
	if err != nil && os.IsNotExist(err) {
		logging.Logger().Warn().Str("module", "config").Msgf("Configuration file not found at %s, using default path %s", path, configPath)
		yamlData, err = os.ReadFile(configPath)
	}
	if err != nil {
		return nil, err
	}
	if err = yaml.Unmarshal(yamlData, &appConfig); err != nil {
		return nil, err
	}
	logging.Logger().Info().Msgf("loaded config\n%s", string(yamlData))
	return &appConfig, nil
}
