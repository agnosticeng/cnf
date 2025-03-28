package cnf

import (
	"dario.cat/mergo"
	mapstructure_hooks "github.com/agnosticeng/mapstructure-hooks"
	"github.com/mitchellh/mapstructure"
)

func LoadStruct[T any](optFuncs ...OptionFunc) (T, error) {
	var (
		v   T
		err = Load(&v, optFuncs...)
	)

	return v, err
}

func Load(i interface{}, optFuncs ...OptionFunc) error {
	var (
		opts = DefaultOptions()
		m    = make(map[string]interface{})
	)

	for _, optFunc := range optFuncs {
		opts = optFunc(opts)
	}

	for _, p := range opts.Providers {
		pm, err := p.ReadMap()

		if err != nil {
			return err
		}

		if err := mergo.Merge(&m, pm, mergo.WithOverride, mergo.WithSliceDeepCopy); err != nil {
			return err
		}
	}

	var hooks []mapstructure.DecodeHookFunc

	hooks = append(hooks, mapstructure.StringToTimeDurationHookFunc())
	hooks = append(hooks, mapstructure_hooks.All()...)

	if len(opts.MapstructureHooks) > 0 {
		hooks = append(hooks, opts.MapstructureHooks...)
	}

	var mdc mapstructure.DecoderConfig

	mdc.Metadata = nil
	mdc.Result = &i
	mdc.WeaklyTypedInput = true
	mdc.Squash = true
	mdc.DecodeHook = mapstructure.ComposeDecodeHookFunc(hooks...)

	d, err := mapstructure.NewDecoder(&mdc)

	if err != nil {
		return err
	}

	return d.Decode(m)
}
