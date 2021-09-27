package settings

import (
	"fmt"
	"time"

	"github.com/arsmn/ontest-server/module/xlog"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

type (
	Provider interface {
		Settings() *Config
	}
	Config struct {
		l *xlog.Logger

		Mode  string
		Serve struct {
			Domain         string
			APIURL         string
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
		Mail struct {
			SMTP struct {
				From       string
				Password   string
				Host, Port string
			}
		}
		Cache struct {
			Redis struct {
				DSN      string
				Password string
				DB       int
			}
		}
		CORS struct {
			AllowedOrigins   []string
			AllowedMethods   []string
			AllowedHeaders   []string
			AllowCredentials bool
		}
		External struct {
			IPGeoLocation struct {
				APIKey string
			}
		}
	}
)

func New(l *xlog.Logger) *Config {
	conf := new(Config)

	// Mode
	conf.Mode = viper.GetString(keyMode)

	// Serve
	conf.Serve.Domain = viper.GetString(keyServeDomain)
	conf.Serve.APIURL = viper.GetString(keyServeAPIURL)
	conf.Serve.StartupMessage = viper.GetBool(keyServeStartupMessageEnabled)
	conf.Serve.Public.Port = viper.GetString(keyServePublicPort)
	conf.Serve.Public.Host = viper.GetString(keyServePublicHost)

	// SQL
	conf.SQL.DSN = viper.GetString(keySQLDSN)
	conf.SQL.Driver = viper.GetString(keySQLDriver)

	// Argon2
	conf.Argon2.Memory = viper.GetUint32(keyHasherArgon2ConfigMemory)
	conf.Argon2.Iterations = viper.GetUint32(keyHasherArgon2ConfigIterations)
	conf.Argon2.Parallelism = uint8(viper.GetUint(keyHasherArgon2ConfigParallelism))
	conf.Argon2.SaltLength = viper.GetUint32(keyHasherArgon2ConfigSaltLength)
	conf.Argon2.KeyLength = viper.GetUint32(keyHasherArgon2ConfigKeyLength)

	// Session
	conf.Session.Cookie = viper.GetString(keySessionCookie)
	conf.Session.Lifespan = viper.GetDuration(keySessionLifespan)

	// OAuth
	conf.OAuth.StateCookie = viper.GetString(keyOAuthStateCookie)
	conf.OAuth.CookieLifespan = viper.GetDuration(keyOAuthCookieLifespan)
	conf.OAuth.Google.ClientID = viper.GetString(keyOAuthGoogleClientID)
	conf.OAuth.Google.ClientSecret = viper.GetString(keyOAuthGoogleClientSecret)
	conf.OAuth.Google.RedirectURL = viper.GetString(keyOAuthGoogleRedirectURL)
	conf.OAuth.Google.Scopes = viper.GetStringSlice(keyOAuthGoogleScopes)
	conf.OAuth.GitHub.ClientID = viper.GetString(keyOAuthGitHubClientID)
	conf.OAuth.GitHub.ClientSecret = viper.GetString(keyOAuthGitHubClientSecret)
	conf.OAuth.GitHub.RedirectURL = viper.GetString(keyOAuthGitHubRedirectURL)
	conf.OAuth.GitHub.Scopes = viper.GetStringSlice(keyOAuthGitHubScopes)
	conf.OAuth.LinkedIn.ClientID = viper.GetString(keyOAuthLinkedInClientID)
	conf.OAuth.LinkedIn.ClientSecret = viper.GetString(keyOAuthLinkedInClientSecret)
	conf.OAuth.LinkedIn.RedirectURL = viper.GetString(keyOAuthLinkedInRedirectURL)
	conf.OAuth.LinkedIn.Scopes = viper.GetStringSlice(keyOAuthLinkedInScopes)

	// Client
	conf.Client.WebURL = viper.GetString(keyClientWebURL)

	// Mail
	conf.Mail.SMTP.From = viper.GetString(keyMailSMTPFrom)
	conf.Mail.SMTP.Password = viper.GetString(keyMailSMTPPassword)
	conf.Mail.SMTP.Host = viper.GetString(keyMailSMTPHost)
	conf.Mail.SMTP.Port = viper.GetString(keyMailSMTPPort)

	// Cache
	conf.Cache.Redis.DSN = viper.GetString(keyCacheRedisDSN)
	conf.Cache.Redis.Password = viper.GetString(keyCacheRedisPassword)
	conf.Cache.Redis.DB = viper.GetInt(keyCacheRedisDB)

	// CORS
	conf.CORS.AllowedOrigins = viper.GetStringSlice(keyCORSAllowedOrigins)
	conf.CORS.AllowedMethods = viper.GetStringSlice(keyCORSAllowedMethods)
	conf.CORS.AllowedHeaders = viper.GetStringSlice(keyCORSAllowedHeaders)
	conf.CORS.AllowCredentials = viper.GetBool(keyCORSAllowCredenials)

	// ExternalAPI
	conf.External.IPGeoLocation.APIKey = viper.GetString(keyExternalIPGeoLocationAPIKey)

	return conf
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

func (c *Config) IsMode(mode string) bool {
	return c.Mode == mode
}

func (c *Config) IsProduction() bool {
	return c.IsMode("production")
}

func (c *Config) IsStaging() bool {
	return c.IsMode("staging")
}

func (c *Config) IsDevelopment() bool {
	return c.IsMode("development")
}

func APIURL() string {
	return viper.GetString(keyServeAPIURL)
}
