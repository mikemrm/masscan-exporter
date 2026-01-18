package masscan

import (
	"context"
	"io"
	"net/http"
	"reflect"
	"slices"
	"strings"
	"time"

	"github.com/go-viper/mapstructure/v2"
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

	Ranges DynamicValue[[]string] `mapstructure:"ranges"`
	Ports  DynamicValue[[]string] `mapstructure:"ports"`

	Config     DynamicValue[string] `mapstructure:"config"`
	ConfigPath string               `mapstructure:"config_path"`
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
		cfg.Ranges.Value = append(slices.Clone(cfg.Ranges.Value), ranges...)

		return cfg
	})
}

func WithPorts(ports ...string) Option {
	return optionFunc(func(cfg Config) Config {
		cfg.Ports.Value = append(slices.Clone(cfg.Ports.Value), ports...)

		return cfg
	})
}

// DynamicValue allows for a value to be dynamically loaded.
//
// When loaded from configuration no matter the DynamicValue T type,
// if provided a string prefixed with env://, file://, http://, or https://.
// The appropriate field is configured.
//
// Any other string value is considered static and will be used as is,
// if the field supports that data type.
type DynamicValue[T any] struct {
	Value     T         `mapstructure:"value"`
	Env       string    `mapstructure:"env"`
	File      string    `mapstructure:"file"`
	URL       string    `mapstructure:"url"`
	URLConfig URLConfig `mapstructure:"url_config"`

	empty T
}

func (v DynamicValue[T]) valueEmpty() bool {
	value := reflect.ValueOf(v.Value)

	valueZero := value.IsZero()

	if !valueZero {
		if value.Kind() == reflect.Slice && value.Len() == 0 {
			return true
		}

		return false
	}

	return true
}

// Configured returns true if any field (except URLConfig) is configured.
func (v DynamicValue[T]) Configured() bool {
	return !v.valueEmpty() || v.Env != "" || v.File != "" || v.URL != ""
}

// GetValue will return the static value if it is not empty.
// Otherwise it will dynamically load the value from the other configuration.
func (v DynamicValue[T]) GetValue(ctx context.Context) (T, error) {
	if !v.valueEmpty() {
		return v.Value, nil
	}

	switch {
	case v.Env != "":
		return loadEnv[T](ctx, v.Env)
	case v.File != "":
		return loadFile[T](ctx, v.File)
	case v.URL != "":
		return loadURL[T](ctx, v.URL, v.URLConfig)
	}

	return v.empty, nil
}

// UnmarshalMapstructure implement mapstructures unmarshaller to decode the value.
func (v *DynamicValue[T]) UnmarshalMapstructure(input any) error {
	switch value := input.(type) {
	case string:
		scheme, remain, found := strings.Cut(value, "://")
		if found {
			switch scheme {
			case "env":
				v.Env = remain

				return nil
			case "file":
				v.File = remain

				return nil
			case "http", "https":
				v.URL = value

				return nil
			}
		}
	}

	type dynamicValue DynamicValue[T]

	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		ErrorUnused: true,
		ZeroFields:  true,
		Result:      (*dynamicValue)(v),
	})
	if err != nil {
		return err
	}

	// Attempt to decode the dynamic structure.
	// If unsuccessful, assume the input is a static value and decode again.
	if err := decoder.Decode(input); err != nil {
		in := map[string]any{
			"value": input,
		}

		if err := decoder.Decode(in); err != nil {
			return err
		}
	}

	return nil
}

// URLConfig allows for the http request to be configured.
// All string values may use the env:// or file:// prefix to load its value dynamically.
type URLConfig struct {
	Method  string            `mapstructure:"method"`
	Auth    URLAuthConfig     `mapstructure:"auth"`
	Headers map[string]string `mapstructure:"headers"`
	Body    string            `mapstructure:"body"`
}

func (c URLConfig) getMethod() string {
	if c.Method != "" {
		return strings.ToUpper(c.Method)
	}

	return http.MethodGet
}

func (c URLConfig) getBody() io.Reader {
	if body := loadValue(c.Body); body != "" {
		return strings.NewReader(body)
	}

	return nil
}

func (c URLConfig) getHeaders() map[string]string {
	ret := make(map[string]string, len(c.Headers))

	for k, v := range c.Headers {
		ret[k] = loadValue(v)
	}

	return ret
}

// URLAuthConfig provides basic an bearer options for authorization.
// Any value may be prefixed with env:// or file:// to dynamically load the value.
type URLAuthConfig struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Bearer   string `mapstructure:"bearer"`
}
