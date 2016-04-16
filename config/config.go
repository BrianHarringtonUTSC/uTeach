// Package config provides functionality to create app specific configuration.
package config

import (
	"path/filepath"

	"github.com/spf13/viper"
)

const (
	envPrefix = "uteach"
)

// Config stores user specific information required to run the app.
// Google credentials should be obtained from the Google Developer Console (https://console.developers.google.com).
type Config struct {
	HTTPAddress                   string `mapstructure:"http_address"`
	DBPath                        string `mapstructure:"db_path"`
	TemplatesPath                 string `mapstructure:"templates_path"`
	StaticFilesPath               string `mapstructure:"static_files_path"`
	CookieAuthenticationKeyBase64 string `mapstructure:"cookie_authentication_key_base64"` // must be a base 64 encoded string of a 64 byte array
	CookieEncryptionKeyBase64     string `mapstructure:"cookie_encryption_key_base64"`     // must be a base 64 encoded string of a 32 byte array
	GoogleRedirectURL             string `mapstructure:"google_redirect_url"`
	GoogleClientID                string `mapstructure:"google_client_id"`
	GoogleClientSecret            string `mapstructure:"google_client_secret"`
}

// Load loads the config file at path and environment variables into a Config.
// The config file can be in any standard format (JSON, YAML, etc).
// The keys in the config file should match the key listed in the "mapstructure" struct tag in the Config struct.
// Keys can also be set as environment variables. Simply attach UTEACH_ to the beginning and capitalize the entire key.
// All keys must be set (even as an empty string) in the config file, even if the key is set as an env variable.
func Load(path string) (*Config, error) {
	if err := loadViper(path); err != nil {
		return nil, err
	}

	c := new(Config)
	if err := viper.Unmarshal(c); err != nil {
		return nil, err
	}

	// for file paths, if it is relative we want the path to be relative to the config's path and not the cwd
	dir := filepath.Dir(path)
	c.DBPath = joinIfNotAbs(dir, c.DBPath)
	c.TemplatesPath = joinIfNotAbs(dir, c.TemplatesPath)
	c.StaticFilesPath = joinIfNotAbs(dir, c.StaticFilesPath)

	return c, nil
}

func joinIfNotAbs(base string, path string) string {
	if !filepath.IsAbs(path) {
		path = filepath.Join(base, path)
	}
	return path
}

func loadViper(path string) error {
	viper.SetEnvPrefix(envPrefix)
	viper.AutomaticEnv()
	viper.SetConfigFile(path)
	return viper.ReadInConfig()
}
