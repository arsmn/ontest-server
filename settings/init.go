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
		viper.AddConfigPath(home)
		viper.AddConfigPath("/etc/ontest/")
		viper.SetConfigName(".ontest")
	} else {
		viper.SetConfigFile(cfgFile)
	}

	viper.SetEnvPrefix("OT")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigParseError); ok {
			panic(err)
		}
	}

	return nil
}

func ConfigFileUsed() string {
	return viper.ConfigFileUsed()
}
