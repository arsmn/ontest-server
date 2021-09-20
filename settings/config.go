package settings

import (
	"fmt"
	"time"

	"github.com/arsmn/ontest-server/module/xlog"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

type (
	Serve struct {
		Domain         string
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
		Lifespan time.Duration
	}
	OAuth struct {
		StateCookie    string
		CookieLifespan time.Duration
		Google         oauth2.Config
		GitHub         oauth2.Config
		LinkedIn       oauth2.Config
	}
	Client struct {
		WebURL string
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
		oauth   OAuth
		client  Client
	}
)

func New(l *xlog.Logger) *Config {
	conf := new(Config)

	// Mode
	conf.mode = viper.GetString(keyMode)

	// Serve
	conf.serve.Domain = viper.GetString(keyServeDomain)
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
	conf.session.Cookie = viper.GetString(keySessionCookie)
	conf.session.Lifespan = viper.GetDuration(keySessionLifespan)

	// OAuth
	conf.oauth.StateCookie = viper.GetString(keyOAuthStateCookie)
	conf.oauth.CookieLifespan = viper.GetDuration(keyOAuthCookieLifespan)
	conf.oauth.Google.ClientID = viper.GetString(keyOAuthGoogleClientID)
	conf.oauth.Google.ClientSecret = viper.GetString(keyOAuthGoogleClientSecret)
	conf.oauth.Google.RedirectURL = viper.GetString(keyOAuthGoogleRedirectURL)
	conf.oauth.Google.Scopes = viper.GetStringSlice(keyOAuthGoogleScopes)
	conf.oauth.GitHub.ClientID = viper.GetString(keyOAuthGitHubClientID)
	conf.oauth.GitHub.ClientSecret = viper.GetString(keyOAuthGitHubClientSecret)
	conf.oauth.GitHub.RedirectURL = viper.GetString(keyOAuthGitHubRedirectURL)
	conf.oauth.GitHub.Scopes = viper.GetStringSlice(keyOAuthGitHubScopes)
	conf.oauth.LinkedIn.ClientID = viper.GetString(keyOAuthLinkedInClientID)
	conf.oauth.LinkedIn.ClientSecret = viper.GetString(keyOAuthLinkedInClientSecret)
	conf.oauth.LinkedIn.RedirectURL = viper.GetString(keyOAuthLinkedInRedirectURL)
	conf.oauth.LinkedIn.Scopes = viper.GetStringSlice(keyOAuthLinkedInScopes)

	// Client
	conf.client.WebURL = viper.GetString(keyClientWebURL)

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

func (c *Config) Domain() string {
	return c.serve.Domain
}

func (c *Config) IsProd() bool {
	return c.mode == "prod"
}

func (c *Config) Session() Session {
	return c.session
}

func (c *Config) OAuth() OAuth {
	return c.oauth
}

func (c *Config) Client() Client {
	return c.client
}
