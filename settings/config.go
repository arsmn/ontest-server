package settings

import (
	"github.com/arsmn/ontest/module/xlog"
	"github.com/spf13/viper"
)

type (
	SQL struct {
		DSN    string
		Driver string
	}
	Argon2 struct {
		Memory      uint32
		Iterations  uint32
		Parallelism uint8
		SaltLength  uint32
		KeyLength   uint32
	}

	Config struct {
		l *xlog.Logger

		sql    *SQL
		argon2 *Argon2
	}
)

func New(l *xlog.Logger) *Config {
	conf := new(Config)

	// sql
	conf.sql.DSN = viper.GetString(keySQLDSN)
	conf.sql.Driver = viper.GetString(keySQLDSN)

	// argon2
	conf.argon2.Memory = viper.GetUint32(keyHasherArgon2ConfigMemory)
	conf.argon2.Iterations = viper.GetUint32(keyHasherArgon2ConfigIterations)
	conf.argon2.Parallelism = uint8(viper.GetUint(keyHasherArgon2ConfigParallelism))
	conf.argon2.SaltLength = viper.GetUint32(keyHasherArgon2ConfigSaltLength)
	conf.argon2.KeyLength = viper.GetUint32(keyHasherArgon2ConfigKeyLength)

	return conf
}

func (c *Config) Argon() *Argon2 {
	return c.argon2
}

func (c *Config) SQL() *SQL {
	return c.sql
}
