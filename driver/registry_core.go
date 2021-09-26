package driver

import (
	"context"

	"github.com/arsmn/ontest-server/app"
	"github.com/arsmn/ontest-server/app/service"
	"github.com/arsmn/ontest-server/module/cache"
	"github.com/arsmn/ontest-server/module/hash"
	"github.com/arsmn/ontest-server/module/mail"
	"github.com/arsmn/ontest-server/module/oauth"
	"github.com/arsmn/ontest-server/module/xlog"
	"github.com/arsmn/ontest-server/persistence"
	"github.com/arsmn/ontest-server/persistence/sql"
	"github.com/arsmn/ontest-server/settings"
)

type RegistryCore struct {
	l *xlog.Logger
	c *settings.Config

	app            app.App
	persister      persistence.Persister
	cacher         cache.Cacher
	passwordHasher hash.Hasher
	mailer         mail.Mailer

	googleOAuth   oauth.OAuther
	githubOAuth   oauth.OAuther
	linkedinOAuth oauth.OAuther
}

func NewRegistryCore() *RegistryCore {
	return &RegistryCore{}
}

func (r *RegistryCore) Init(ctx context.Context) (err error) {
	r.persister, err = sql.NewPersister(r)
	if err != nil {
		return err
	}

	r.mailer, err = mail.NewMailerSMTP(r)
	if err != nil {
		return err
	}

	r.app = service.NewAppService(r)
	r.cacher = cache.NewCacherRedis(r)
	r.passwordHasher = hash.NewHasherArgon2(r)

	r.googleOAuth = oauth.NewOAutherGoogle(r)
	r.githubOAuth = oauth.NewOAutherGitHub(r)
	r.linkedinOAuth = oauth.NewOAutherLinkedIn(r)

	return nil
}

func (r *RegistryCore) WithLogger(l *xlog.Logger) Registry {
	r.l = l
	return r
}

func (r *RegistryCore) WithConfig(c *settings.Config) Registry {
	r.c = c
	return r
}

func (r *RegistryCore) Logger() *xlog.Logger {
	return r.l
}

func (r *RegistryCore) Settings() *settings.Config {
	return r.c
}

func (r *RegistryCore) App() app.App {
	return r.app
}

func (r *RegistryCore) Persister() persistence.Persister {
	return r.persister
}

func (r *RegistryCore) Cacher() cache.Cacher {
	return r.cacher
}

func (r *RegistryCore) Hasher() hash.Hasher {
	return r.passwordHasher
}

func (r *RegistryCore) Mailer() mail.Mailer {
	return r.mailer
}

func (r *RegistryCore) OAuther(typ oauth.OAuthProviderType) oauth.OAuther {
	switch typ {
	case oauth.GoogleType:
		return r.googleOAuth
	case oauth.GitHubType:
		return r.githubOAuth
	case oauth.LinkedInType:
		return r.linkedinOAuth
	default:
		return oauth.NewOAutherNop()
	}
}
