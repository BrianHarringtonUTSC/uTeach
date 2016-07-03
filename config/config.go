// Package config provides functionality to create app specific configuration.
package config

import (
	"encoding/base64"
	"path/filepath"

	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

const (
	envPrefix = "uteach"

	// Auth0 URLs for OAuth2
	base        = "https://uteach.auth0.com"
	authURL     = base + "/authorize"
	tokenURL    = base + "/oauth/token"
	userInfoURL = base + "/userinfo"
)

// Config stores context information required to run the app.
type Config struct {
	HTTPAddress             string
	DBPath                  string
	TemplatesPath           string
	StaticFilesPath         string
	CookieAuthenticationKey []byte
	CookieEncryptionKey     []byte
	OAuth2                  *oauth2.Config
	OAuth2UserInfoURL       string
}

// preprocessedConfig is created using env variables and config files and should be used to create the final "Config" above.
type preprocessedConfig struct {
	HTTPAddress                   string `mapstructure:"http_address"`
	DBPath                        string `mapstructure:"db_path"`
	TemplatesPath                 string `mapstructure:"templates_path"`
	StaticFilesPath               string `mapstructure:"static_files_path"`
	CookieAuthenticationKeyBase64 string `mapstructure:"cookie_authentication_key_base64"` // must be a base64 encoded string of a 64 byte array
	CookieEncryptionKeyBase64     string `mapstructure:"cookie_encryption_key_base64"`     // must be a base64 encoded string of a 32 byte array
	OAuth2ClientID                string `mapstructure:"oauth2_client_id"`
	OAuth2ClientSecret            string `mapstructure:"oauth2_client_secret"`
	OAuth2RedirectURL             string `mapstructure:"oauth2_redirect_url"`
}

// Load loads the config file at path and environment variables into a Config.
// The config file can be in any standard format (JSON, YAML, etc).
// The keys in the config file should match the key listed in the "mapstructure" struct tag in the preprocessedConfig struct.
// Keys can also be set as environment variables. Simply attach UTEACH_ to the beginning and capitalize the entire key.
// All keys must be set (even as an empty string) in the config file, even if the key is set as an env variable.
func Load(path string) (*Config, error) {
	if err := loadViper(path); err != nil {
		return nil, err
	}

	preprocessed := new(preprocessedConfig)
	if err := viper.Unmarshal(preprocessed); err != nil {
		return nil, err
	}

	conf := new(Config)
	conf.HTTPAddress = preprocessed.HTTPAddress

	// for file paths, if it is relative we want the path to be relative to the config's path and not the cwd
	dir := filepath.Dir(path)
	conf.DBPath = joinIfNotAbs(dir, preprocessed.DBPath)
	conf.TemplatesPath = joinIfNotAbs(dir, preprocessed.TemplatesPath)
	conf.StaticFilesPath = joinIfNotAbs(dir, preprocessed.StaticFilesPath)

	var err error
	conf.CookieAuthenticationKey, err = base64.StdEncoding.DecodeString(preprocessed.CookieAuthenticationKeyBase64)
	if err != nil {
		return nil, err
	}

	conf.CookieEncryptionKey, err = base64.StdEncoding.DecodeString(preprocessed.CookieEncryptionKeyBase64)
	if err != nil {
		return nil, err
	}

	oauth2Config := &oauth2.Config{
		ClientID:     preprocessed.OAuth2ClientID,
		ClientSecret: preprocessed.OAuth2ClientSecret,
		RedirectURL:  preprocessed.OAuth2RedirectURL,
		Scopes:       []string{"openid", "name", "email", "nickname"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  authURL,
			TokenURL: tokenURL,
		},
	}

	conf.OAuth2 = oauth2Config
	conf.OAuth2UserInfoURL = userInfoURL

	return conf, nil
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
