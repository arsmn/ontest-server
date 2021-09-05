package settings

import (
	"github.com/spf13/viper"
)

type (
	Argon2 struct {
		Memory      uint32 `json:"memory"`
		Iterations  uint32 `json:"iterations"`
		Parallelism uint8  `json:"parallelism"`
		SaltLength  uint32 `json:"salt_length"`
		KeyLength   uint32 `json:"key_length"`
	}

	Config struct {
		argon2 Argon2
	}
)

func NewConfig() *Config {
	conf := new(Config)

	// Argon2
	conf.argon2.Memory = viper.GetUint32(keyHasherArgon2ConfigMemory)
	conf.argon2.Iterations = viper.GetUint32(keyHasherArgon2ConfigIterations)
	conf.argon2.Parallelism = uint8(viper.GetUint(keyHasherArgon2ConfigParallelism))
	conf.argon2.SaltLength = viper.GetUint32(keyHasherArgon2ConfigSaltLength)
	conf.argon2.KeyLength = viper.GetUint32(keyHasherArgon2ConfigKeyLength)

	return conf
}
