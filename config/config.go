// Package config provides functionality to store user specific info for the uTeach app to function.
package config

import (
	"github.com/spf13/viper"
	"log"
	"path/filepath"
)

// Config stores user specific information to run the app.
// Google credentials should be obtained from the Google Developer Console (https://console.developers.google.com).
type Config struct {
	HTTPAddress             string `mapstructure:"http_address"`
	DBPath                  string `mapstructure:"db_path"`
	TemplatesPath           string `mapstructure:"templates_path"`
	StaticFilesPath         string `mapstructure:"static_files_path"`
	CookieAuthenticationKey string `mapstructure:"cookie_authentication_key"`
	CookieEncryptionKey     string `mapstructure:"cookie_encryption_key"`
	GoogleRedirectURL       string `mapstructure:"google_redirect_url"`
	GoogleClientID          string `mapstructure:"google_client_id"`
	GoogleClientSecret      string `mapstructure:"google_client_secret"`
}

func joinIfNotAbs(base string, path string) string {
	if !filepath.IsAbs(path) {
		path = filepath.Join(base, path)
	}
	return path
}

func loadViper(path string) {
	viper.SetEnvPrefix("uteach")
	viper.AutomaticEnv()
	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}
}

// Load loads the config file at path and environment variables into a Config.
// NOTE: all keys must be set (even if given an empty value) in the files, even if set as env variable.
func Load(path string) *Config {
	loadViper(path)

	c := &Config{}
	if err := viper.Unmarshal(c); err != nil {
		log.Fatal(err)
	}

	// for file paths, we want the path to be relative to the config's path
	dir := filepath.Dir(path)
	c.DBPath = joinIfNotAbs(dir, c.DBPath)
	c.TemplatesPath = joinIfNotAbs(dir, c.TemplatesPath)
	c.StaticFilesPath = joinIfNotAbs(dir, c.StaticFilesPath)

	return c
}
