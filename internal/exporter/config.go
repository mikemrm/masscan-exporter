package exporter

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
	"github.com/spf13/pflag"
)

type Config struct {
	logger     *zerolog.Logger
	Timeout    time.Duration         `mapstructure:"timeout"`
	Registerer prometheus.Registerer `mapstructure:"-"`
	Labels     prometheus.Labels     `mapstructure:"labels"`
}

func newConfig(opts ...Option) Config {
	var cfg Config

	for _, opt := range opts {
		cfg = opt.apply(cfg)
	}

	if cfg.logger == nil {
		l := zerolog.Nop()

		cfg.logger = &l
	}

	if cfg.Timeout <= 0 {
		cfg.Timeout = 10 * time.Minute
	}

	if cfg.Registerer == nil {
		cfg.Registerer = prometheus.DefaultRegisterer
	}

	return cfg
}

type Option interface {
	apply(Config) Config
}

type optionFunc func(Config) Config

func (fn optionFunc) apply(cfg Config) Config {
	return fn(cfg)
}

// WithConfig replaces the existing Config.
func WithConfig(cfg Config) Option {
	return optionFunc(func(_ Config) Config {
		return cfg
	})
}

// WithLogger configures the exporter logger.
func WithLogger(logger zerolog.Logger) Option {
	return optionFunc(func(cfg Config) Config {
		cfg.logger = &logger

		return cfg
	})
}

// WithRegisterer configures the prometheus registerer.
func WithRegisterer(reg prometheus.Registerer) Option {
	return optionFunc(func(cfg Config) Config {
		cfg.Registerer = reg

		return cfg
	})
}

// WithTimeout configures the timeout running masscan.
func WithTimeout(timeout time.Duration) Option {
	return optionFunc(func(cfg Config) Config {
		cfg.Timeout = timeout

		return cfg
	})
}

// WithLabels configures constant labels.
func WithLabels(labels prometheus.Labels) Option {
	return optionFunc(func(cfg Config) Config {
		cfg.Labels = labels

		return cfg
	})
}
func AddFlags(flags *pflag.FlagSet) {
	flags.Duration("exporter.timeout", 10*time.Minute, "sets the timeout for the scan to complete")
	flags.StringToString("exporter.labels", nil, "configure constant prometheus labels")
}
