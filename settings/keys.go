package settings

const (
	keyMode                          = "mode"
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
	keySessionDomain                 = "session.domain"
)
