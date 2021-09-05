package settings

type Settings interface {
	Argon2() Argon2
}

type Provider interface {
	Settings() Settings
}
