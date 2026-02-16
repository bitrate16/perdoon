package config

import "flag"

// Get config path from commandline arguments
func GetConfigPath() string {
	var ConfigPath string

	flag.StringVar(
		&ConfigPath,
		"config",
		"config.yml",
		"Path to config file",
	)

	flag.Parse()

	return ConfigPath
}
