package settings

import (
	"fmt"
	"time"

	"github.com/arsmn/ontest-server/module/xlog"
	"github.com/spf13/viper"
)

type (
	Serve struct {
		StartupMessage bool
		Public         struct {
			Port string
			Host string
		}
	}
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
	Session struct {
		Cookie   string
		Domain   string
		Lifespan time.Duration
	}
	Provider interface {
		Settings() *Config
	}
	Config struct {
		l *xlog.Logger

		mode    string
		serve   Serve
		sql     SQL
		argon2  Argon2
		session Session
	}
)

func New(l *xlog.Logger) *Config {
	conf := new(Config)

	// Mode
	conf.mode = viper.GetString(keyMode)

	// Serve
	conf.serve.StartupMessage = viper.GetBool(keyServeStartupMessageEnabled)
	conf.serve.Public.Port = viper.GetString(keyServePublicPort)
	conf.serve.Public.Host = viper.GetString(keyServePublicHost)

	// SQL
	conf.sql.DSN = viper.GetString(keySQLDSN)
	conf.sql.Driver = viper.GetString(keySQLDriver)

	// Argon2
	conf.argon2.Memory = viper.GetUint32(keyHasherArgon2ConfigMemory)
	conf.argon2.Iterations = viper.GetUint32(keyHasherArgon2ConfigIterations)
	conf.argon2.Parallelism = uint8(viper.GetUint(keyHasherArgon2ConfigParallelism))
	conf.argon2.SaltLength = viper.GetUint32(keyHasherArgon2ConfigSaltLength)
	conf.argon2.KeyLength = viper.GetUint32(keyHasherArgon2ConfigKeyLength)

	// Session
	conf.session.Domain = viper.GetString(keySessionDomain)
	conf.session.Cookie = viper.GetString(keySessionCookie)
	conf.session.Lifespan = viper.GetDuration(keySessionLifespan)

	return conf
}

func (c *Config) StartupMessageEnabled() bool {
	return c.serve.StartupMessage
}

func (c *Config) PublicListenOn() string {
	return c.listenOn("public")
}

func (c *Config) listenOn(key string) string {
	port := viper.GetInt("serve." + key + ".port")
	if port < 1 {
		c.l.Fatal(fmt.Sprintf("serve.%s.port can not be zero or negative", key))
	}

	return fmt.Sprintf("%s:%d", viper.GetString("serve."+key+".host"), port)
}

func (c *Config) HasherArgon2() Argon2 {
	return c.argon2
}

func (c *Config) SQL() SQL {
	return c.sql
}

func (c *Config) Mode() string {
	return c.mode
}

func (c *Config) IsProd() bool {
	return c.mode == "prod"
}

func (c *Config) Session() Session {
	return c.session
}
