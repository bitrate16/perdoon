package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// LoadConfig reads config from file or creates default if missing
func LoadConfig(path string) (*Config, error) {
	cfg := &Config{}

	// Create default if not found
	if _, err := os.Stat(path); os.IsNotExist(err) {
		cfg = defaultConfig()
		if err := writeConfig(path, cfg); err != nil {
			return nil, err
		}

		return cfg, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func writeConfig(path string, cfg *Config) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := yaml.NewEncoder(f)
	enc.SetIndent(2)
	defer enc.Close()

	return enc.Encode(cfg)
}
