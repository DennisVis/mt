// Copyright (c) 2021 Dennis Vis
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package mt

type config struct {
	SkipValidation bool
	Lax            bool
	StopOnError    bool
}

type option = func(cfg config) config

var defaultConfig = config{
	SkipValidation: false,
	Lax:            false,
	StopOnError:    false,
}

// SkipValidation will skip message validation and return messages as-is. The difference with Lax is that with this
// option validation is skipped entirely and there won't be any gathering of validation errors. This option is
// therefore only useful when performance is more important than validity.
//
// Default: false
func SkipValidation(skip bool) option {
	return func(cfg config) config {
		cfg.SkipValidation = skip
		return cfg
	}
}

// Lax will cause messages that fail validation to be returned anyway. Otherwise invalid messages are discarded. In
// both cases any validation errors will be returned as parse errors. This is the main difference with
// SkipValidation as validation is still performed and its errors are sent to the client.
//
// Default: false
func Lax(lax bool) option {
	return func(cfg config) config {
		cfg.Lax = lax
		return cfg
	}
}

// StopOnError will make the parsing process stop on the first error.
//
// Default: false
func StopOnError(stop bool) option {
	return func(cfg config) config {
		cfg.StopOnError = stop
		return cfg
	}
}

func optionsToConfig(option []option) config {
	cfg := defaultConfig

	for _, opt := range option {
		cfg = opt(cfg)
	}

	return cfg
}
