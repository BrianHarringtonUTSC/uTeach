// Package config provides functionality to store user specific info for the uTeach app to function.
package config

import (
	"log"
	"path/filepath"

	"github.com/spf13/viper"
)

const (
	envPrefix = "uteach"
)

// Config stores user specific information required to run the app.
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

// Load loads the config file at path and environment variables into a Config.
// Thanks to viper, the config file can be in any standard format (JSON, YAML, etc). Viper will detect the file
// extension and parse the file appropriately. The keys in the config file should match the key listed in the
// "mapstructure" struct tag in the Config struct. Viper will even detect environment variables. Simply attach
// UTEACH_ to the beginning and capitalize the entire variable.
// NOTE: all keys must be set (even as an empty string) in the config file, even if set as env variable.
// This is so viper knows to look for the key in the environment variables when creating the config.
// Example: to set the HTTPAddress, you can create a JSON file and set the key "http_adress": ":8000" OR
// you can set the key to an empty string "http_address": "", and then export the environment variable
// UTEACH_HTTP_ADDRESS=":8000". Both will work (if both defined the environment variable takes precedence).
// This is useful for defining sensitive information such as Google client secret, etc. An example config can be found
// at ../sample/config.json. Note that the Google Client ID and Secret are not defined in the sample as we never want to
// check in sensitive info like that to the repo. Instead, you must define them as env variables or move the file out of
// the repo and set them in the config file.
func Load(path string) *Config {
	loadViper(path)

	c := &Config{}
	if err := viper.Unmarshal(c); err != nil {
		log.Fatal(err)
	}

	// for file paths, if it is relative we want the path to be relative to the config's path and not the cwd
	dir := filepath.Dir(path)
	c.DBPath = joinIfNotAbs(dir, c.DBPath)
	c.TemplatesPath = joinIfNotAbs(dir, c.TemplatesPath)
	c.StaticFilesPath = joinIfNotAbs(dir, c.StaticFilesPath)

	return c
}

func joinIfNotAbs(base string, path string) string {
	if !filepath.IsAbs(path) {
		path = filepath.Join(base, path)
	}
	return path
}

func loadViper(path string) {
	viper.SetEnvPrefix(envPrefix)
	viper.AutomaticEnv()
	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}
}
