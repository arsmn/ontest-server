package settings

const (
	keyMode                          = "mode"
	keyServeDomain                   = "serve.domain"
	keyServeStartupMessageEnabled    = "serve.startup_message_enabled"
	keyServePublicPort               = "serve.public.port"
	keyServePublicHost               = "serve.public.host"
	keySQLDSN                        = "sql.dsn"
	keySQLDriver                     = "sql.driver"
	keyHasherArgon2ConfigMemory      = "hashers.argon2.memory"
	keyHasherArgon2ConfigIterations  = "hashers.argon2.iterations"
	keyHasherArgon2ConfigParallelism = "hashers.argon2.parallelism"
	keyHasherArgon2ConfigSaltLength  = "hashers.argon2.salt_length"
	keyHasherArgon2ConfigKeyLength   = "hashers.argon2.key_length"
	keySessionCookie                 = "session.cookie"
	keySessionLifespan               = "session.lifespan"
	keyOAuthStateCookie              = "oauth.state_cookie"
	keyOAuthCookieLifespan           = "oauth.cookie_lifespan"
	keyOAuthGoogleClientID           = "oauth.google.client_id"
	keyOAuthGoogleClientSecret       = "oauth.google.client_secret"
	keyOAuthGoogleRedirectURL        = "oauth.google.redirect_url"
	keyOAuthGoogleScopes             = "oauth.google.scopes"
	keyOAuthGitHubClientID           = "oauth.github.client_id"
	keyOAuthGitHubClientSecret       = "oauth.github.client_secret"
	keyOAuthGitHubRedirectURL        = "oauth.github.redirect_url"
	keyOAuthGitHubScopes             = "oauth.github.scopes"
	keyOAuthLinkedInClientID         = "oauth.linkedin.client_id"
	keyOAuthLinkedInClientSecret     = "oauth.linkedin.client_secret"
	keyOAuthLinkedInRedirectURL      = "oauth.linkedin.redirect_url"
	keyOAuthLinkedInScopes           = "oauth.linkedin.scopes"
	keyClientWebURL                  = "client.web_url"
	keyMailSMTPFrom                  = "mail.smtp.from"
	keyMailSMTPPassword              = "mail.smtp.password"
	keyMailSMTPHost                  = "mail.smtp.host"
	keyMailSMTPPort                  = "mail.smtp.port"
)
