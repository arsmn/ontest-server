package settings

import (
	"runtime"

	"github.com/spf13/viper"
)

const (
	argon2DefaultMemory     uint32 = 4 * 1024 * 1024
	argon2DefaultIterations uint32 = 4
	argon2DefaultSaltLength uint32 = 16
	argon2DefaultKeyLength  uint32 = 32
)

var argon2DefaultParallelism = uint8(runtime.NumCPU() * 2)

func setDefaults() {
	// Serve
	viper.SetDefault(keyServePublicPort, 8080)

	// Argon2
	viper.SetDefault(keyHasherArgon2ConfigMemory, argon2DefaultMemory)
	viper.SetDefault(keyHasherArgon2ConfigIterations, argon2DefaultIterations)
	viper.SetDefault(keyHasherArgon2ConfigKeyLength, argon2DefaultSaltLength)
	viper.SetDefault(keyHasherArgon2ConfigKeyLength, argon2DefaultKeyLength)
	viper.SetDefault(keyHasherArgon2ConfigParallelism, argon2DefaultParallelism)
}
