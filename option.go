package cnf

import "github.com/mitchellh/mapstructure"

type OptionFunc func(*Options) *Options

type Options struct {
	Providers         []Provider
	MapstructureHooks []mapstructure.DecodeHookFunc
}

func DefaultOptions() *Options {
	return &Options{}
}

func WithProvider(p Provider) OptionFunc {
	return func(opts *Options) *Options {
		opts.Providers = append(opts.Providers, p)
		return opts
	}
}

func WithMapstructureHooks(hooks ...mapstructure.DecodeHookFunc) OptionFunc {
	return func(opts *Options) *Options {
		opts.MapstructureHooks = hooks
		return opts
	}
}
