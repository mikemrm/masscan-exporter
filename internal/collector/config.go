package collector

import (
	"errors"
	"time"

	"github.com/adhocore/gronx"
	"github.com/mikemrm/masscan-exporter/internal/masscan"
)

var (
	ErrNameRequired    = errors.New("collector name required")
	ErrInvalidSchedule = errors.New("invalid collector schedule")
)

type Config struct {
	Name        string         `mapstructure:"name"`
	Schedule    string         `mapstructure:"schedule"`
	ScanOnStart bool           `mapstructure:"scan_on_start"`
	StartDelay  time.Duration  `mapstructure:"start_delay"`
	Masscan     masscan.Config `mapstructure:"masscan"`
	Timeout     time.Duration  `mapstructure:"timeout"`
}

func (c Config) Validate() error {
	if c.Name == "" {
		return ErrNameRequired
	}

	if c.Schedule == "" || !gronx.IsValid(c.Schedule) {
		return ErrInvalidSchedule
	}

	return nil
}

func newConfig(opts ...Option) Config {
	var cfg Config

	for _, opt := range opts {
		cfg = opt.apply(cfg)
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
