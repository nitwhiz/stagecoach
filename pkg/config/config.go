package config

import (
	"errors"
	"github.com/spf13/viper"
	"os"
)

type config struct {
	DestinationDirectory string `mapstructure:"destinationDirectory"`
	AuthorizationToken   string `mapstructure:"authorizationToken"`
}

var C config

func Load() error {
	viper.SetConfigName("stagecoach")
	viper.SetConfigType("yaml")

	viper.AddConfigPath("/etc/stagecoach/")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if err := viper.Unmarshal(&C); err != nil {
		return err
	}

	return C.Validate()
}

func (c *config) Validate() error {
	stat, err := os.Stat(c.DestinationDirectory)

	if err != nil {
		return err
	}

	if !stat.IsDir() {
		return errors.New("'" + c.DestinationDirectory + "' is not a directory")
	}

	return nil
}
