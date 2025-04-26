package exporter

import (
	"github.com/mikemrm/masscan-exporter/internal/collector"
	"github.com/prometheus/client_golang/prometheus"
)

type Config struct {
	Registerer prometheus.Registerer  `mapstructure:"-"`
	Collectors []*collector.Collector `mapstructure:"-"`
}

func newConfig(opts ...Option) Config {
	var cfg Config

	for _, opt := range opts {
		cfg = opt.apply(cfg)
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

// WithRegisterer configures the prometheus registerer.
func WithRegisterer(reg prometheus.Registerer) Option {
	return optionFunc(func(cfg Config) Config {
		cfg.Registerer = reg

		return cfg
	})
}
