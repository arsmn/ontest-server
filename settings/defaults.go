package settings

import (
	"runtime"
	"time"

	"github.com/spf13/viper"
)

const (
	modeDefault                          = "dev"
	serveDefaultPublicPort               = 8080
	hasherArgon2DefaultMemory     uint32 = 64 * 1024
	hasherArgon2DefaultIterations uint32 = 1
	hasherArgon2DefaultSaltLength uint32 = 16
	hasherArgon2DefaultKeyLength  uint32 = 32
	sessionDefaultCookie                 = "ot_token"
	sessionDefaultLifespan               = time.Hour * 24
	oauthDefaultStateCookie              = "ot_state"
	oauthDefaultCookieLifespan           = time.Hour * 24
)

var hasherArgon2DefaultParallelism = uint8(runtime.NumCPU())

func setDefaults() {
	// Mode
	viper.SetDefault(keyMode, modeDefault)

	// Serve
	viper.SetDefault(keyServePublicPort, serveDefaultPublicPort)

	// Argon2
	viper.SetDefault(keyHasherArgon2ConfigMemory, hasherArgon2DefaultMemory)
	viper.SetDefault(keyHasherArgon2ConfigIterations, hasherArgon2DefaultIterations)
	viper.SetDefault(keyHasherArgon2ConfigKeyLength, hasherArgon2DefaultSaltLength)
	viper.SetDefault(keyHasherArgon2ConfigKeyLength, hasherArgon2DefaultKeyLength)
	viper.SetDefault(keyHasherArgon2ConfigParallelism, hasherArgon2DefaultParallelism)

	// Session
	viper.SetDefault(keySessionCookie, sessionDefaultCookie)
	viper.SetDefault(keySessionLifespan, sessionDefaultLifespan)

	// OAuth
	viper.SetDefault(keyOAuthStateCookie, oauthDefaultStateCookie)
	viper.SetDefault(keyOAuthCookieLifespan, oauthDefaultCookieLifespan)
}
