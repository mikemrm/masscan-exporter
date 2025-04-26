package masscan

import (
	"reflect"
	"strings"
	"time"
)

const (
	DefaultBinPath   = "/usr/bin/masscan"
	DefaultTempDir   = "/tmp"
	DefaultWaitDelay = 20 * time.Second
)

type Config struct {
	TempDir   string        `mapstructure:"temp_dir"`
	BinPath   string        `mapstructure:"bin_path"`
	WaitDelay time.Duration `mapstructure:"wait_delay"`
	MaxRate   int           `mapstructure:"max_rate"`

	Ranges []string `mapstructure:"ranges"`
	Ports  []string `mapstructure:"ports"`

	Config     MasscanConfig `mapstructure:"config"`
	ConfigPath string        `mapstructure:"config_path"`
}

func newConfig(opts ...Option) Config {
	var cfg Config

	for _, opt := range opts {
		cfg = opt.apply(cfg)
	}

	if cfg.TempDir == "" {
		cfg.TempDir = DefaultTempDir
	}

	if cfg.BinPath == "" {
		cfg.BinPath = DefaultBinPath
	}

	if cfg.WaitDelay <= 0 {
		cfg.WaitDelay = DefaultWaitDelay
	}

	return cfg
}

type MasscanConfig string

func (c *MasscanConfig) Decode(from reflect.Value) (any, error) {
	if from.Kind() != reflect.String {
		return from.Interface(), nil
	}

	*c = MasscanConfig(strings.TrimSpace(from.Interface().(string)))

	return nil, nil
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

func WithRanges(ranges ...string) Option {
	return optionFunc(func(cfg Config) Config {
		cfg.Ranges = append(cfg.Ranges, ranges...)

		return cfg
	})
}

func WithPorts(ports ...string) Option {
	return optionFunc(func(cfg Config) Config {
		cfg.Ports = append(cfg.Ports, ports...)

		return cfg
	})
}
