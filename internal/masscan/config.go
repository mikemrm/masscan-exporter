package masscan

import (
	"reflect"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/spf13/pflag"
)

const (
	DefaultBinPath   = "/usr/bin/masscan"
	DefaultTempDir   = "/tmp"
	DefaultWaitDelay = 20 * time.Second
)

type Config struct {
	logger    *zerolog.Logger
	TempDir   string        `mapstructure:"temp_dir"`
	BinPath   string        `mapstructure:"bin_path"`
	WaitDelay time.Duration `mapstructure:"wait_delay"`
	MaxRate   int           `mapstructure:"max_rate"`

	Ranges []string `mapstructure:"ranges"`
	Ports  []string `mapstructure:"ports"`

	Config     MasscanConfig `mapstructure:"config"`
	ConfigPath string        `mapstructure:"config_path"`
}

func (c Config) withDefaults() Config {
	if c.logger == nil {
		l := zerolog.Nop()

		c.logger = &l
	}

	if c.TempDir == "" {
		c.TempDir = DefaultTempDir
	}

	if c.BinPath == "" {
		c.BinPath = DefaultBinPath
	}

	if c.WaitDelay <= 0 {
		c.WaitDelay = DefaultWaitDelay
	}

	return c
}

type MasscanConfig string

func (c *MasscanConfig) Decode(from reflect.Value) (any, error) {
	if from.Kind() != reflect.String {
		return from.Interface(), nil
	}

	*c = MasscanConfig(strings.TrimSpace(from.Interface().(string)))

	return nil, nil
}

type Option func(c Config) Config

func WithConfig(c Config) Option {
	return func(_ Config) Config {
		return c
	}
}

func WithLogger(logger zerolog.Logger) Option {
	return func(c Config) Config {
		c.logger = &logger

		return c
	}
}

func WithRanges(ranges ...string) Option {
	return func(c Config) Config {
		c.Ranges = append(c.Ranges, ranges...)

		return c
	}
}

func WithPorts(ports ...string) Option {
	return func(c Config) Config {
		c.Ports = append(c.Ports, ports...)

		return c
	}
}

func AddFlags(flags *pflag.FlagSet) {
	flags.String("masscan.temp_dir", "/tmp", "configures the temporary directory used by masscan")
	flags.String("masscan.bin_path", "/usr/bin/masscan", "sets the masscan binary path")
	flags.Duration("masscan.wait_delay", 20*time.Second, "sets the delay to wait for the command to exit")
	flags.Int("masscan.max_rate", 0, "sets the max rate to run the scan")

	flags.StringSlice("masscan.ranges", nil, "set the ip ranges to scan")
	flags.StringSlice("masscan.ports", nil, "configures the ports to scan")

	flags.String("masscan.config_path", "", "provide a masscan config file")
	flags.String("masscan.config", "", "provide a masscan config")
}
