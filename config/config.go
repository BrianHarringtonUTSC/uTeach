// Package config provides functionality to store user specific info for the uTeach app to function.
package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Config stores user specific information to run the app.
type Config struct {
	Host            string `json:"host"`
	Port            int    `json:"port"`
	DBPath          string `json:"db_path"`
	TemplatesPath   string `json:"templates_path"`
	StaticFilesPath string `json:"static_files_path"`
}

// Load loads the json formatted file at path into a Config. Panics if it cannot decode the file.
func Load(path string) *Config {
	config := &Config{}

	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	err = json.NewDecoder(file).Decode(config)
	if err != nil {
		panic(err)
	}

	// make all paths relative to the config file's path
	configDir := filepath.Dir(path)
	config.DBPath = filepath.Join(configDir, config.DBPath)
	config.TemplatesPath = filepath.Join(configDir, config.TemplatesPath)
	config.StaticFilesPath = filepath.Join(configDir, config.StaticFilesPath)

	return config
}
