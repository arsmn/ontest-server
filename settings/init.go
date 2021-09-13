package settings

import (
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

func Init(cfgFile string) error {
	if cfgFile == "" {
		home, err := homedir.Dir()
		if err != nil {
			return err
		}

		viper.AddConfigPath(".")
		viper.AddConfigPath("../")
		viper.AddConfigPath(home)
		viper.AddConfigPath("/etc/ontest/")
		viper.SetConfigName("ontest")
	} else {
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigType("toml")
	viper.SetEnvPrefix("ontest")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	setDefaults()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigParseError); ok {
			panic(err)
		}
	}

	return nil
}

func ConfigFileUsed() string {
	cfg := viper.ConfigFileUsed()
	if cfg == "" {
		cfg = "no config file used"
	}
	return cfg
}
