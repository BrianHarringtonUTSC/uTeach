// Package config provides functionality to store user specific info for the uTeach app to function.
package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Config stores user specific information to run the app.
// Google credentials should be obtained from the Google Developer Console (https://console.developers.google.com).
type Config struct {
	HTTPAddress             string `json:"http_address"`
	DBPath                  string `json:"db_path"`
	TemplatesPath           string `json:"templates_path"`
	StaticFilesPath         string `json:"static_files_path"`
	CookieAuthenticationKey string `json:"cookie_authentication_key"`
	CookieEncryptionKey     string `json:"cookie_encryption_key"`
	GoogleRedirectURL       string `json:"google_redirect_url"`
	GoogleClientID          string // env variable
	GoogleClientSecret      string // env variable
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

	config.GoogleClientID = os.Getenv("UTEACH_GOOGLE_CLIENT_ID")
	config.GoogleClientSecret = os.Getenv("UTEACH_GOOGLE_CLIENT_SECRET")

	if config.GoogleClientID == "" || config.GoogleClientSecret == "" {
		panic("UTEACH_GOOGLE_CLIENT_ID and/or UTEACH_GOOGLE_CLIENT_SECRET not set in environment.")
	}

	return config
}
