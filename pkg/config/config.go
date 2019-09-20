// Package config handles the registration of environment variables and preferred filesystem paths.
package config

import (
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"path"
)

var v *viper.Viper

// New returns a preconfigured Viper configuration helper
func New() (*viper.Viper, error) {
	if v == nil {
		return setupViper()
	}
	return v, nil
}

// setupViper collects config settings and returns a Viper object
func setupViper() (*viper.Viper, error) {
	v = viper.New()
	v.SetConfigName("config")
	home, err := homedir.Dir()
	if err != nil {
		return nil, err
	}
	v.AddConfigPath(path.Join(home, ".chilioverflow"))
	v.AddConfigPath(".")
	v.WatchConfig()

	v.SetEnvPrefix("chilioverflow")

	v.SetDefault("configPath", path.Join(home, ".chilioverflow"))
	v.SetDefault("apiUrl", "localhost")

	if err = v.ReadInConfig(); err != nil {
		return nil, err
	}
	return v, nil
}
