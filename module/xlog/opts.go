package xlog

func newOptions(opts []Option) *options {
	o := new(options)
	for _, f := range opts {
		f(o)
	}
	return o
}

func SetLevel(level string) Option {
	return func(o *options) {
		o.level = getZapLevel(level)
	}
}

func ForceEncoder(enc string) Option {
	return func(o *options) {
		o.encoder = enc
	}
}
