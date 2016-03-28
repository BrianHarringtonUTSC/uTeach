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

func makeAbs(base string, path string) string {
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

// Load loads the json formatted file at path into a Config. Panics if it cannot decode the file.
func Load(path string) *Config {
	loadViper(path)

	dir := filepath.Dir(path)

	// Note: viper has an unmarshal which will convert to a struct automaticallly
	// however, it wasn't setting env variables when i unmarshalled, so I'm doing getting each field manually for now
	c := &Config{}
	c.HTTPAddress = viper.GetString("http_address")
	c.DBPath = makeAbs(dir, viper.GetString("db_path"))
	c.TemplatesPath = makeAbs(dir, viper.GetString("templates_path"))
	c.StaticFilesPath = makeAbs(dir, viper.GetString("static_files_path"))
	c.CookieAuthenticationKey = viper.GetString("cookie_authentication_key")
	c.CookieEncryptionKey = viper.GetString("cookie_encryption_key")
	c.GoogleRedirectURL = viper.GetString("google_redirect_url")
	c.GoogleClientID = viper.GetString("google_client_id")
	c.GoogleClientSecret = viper.GetString("google_client_secret")

	return c
}
