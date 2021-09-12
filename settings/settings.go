package settings

type Provider interface {
	Settings() *Config
}
