package cmd

import (
	"context"
	"os"
	"reflect"
	"strings"

	"github.com/go-viper/mapstructure/v2"
	"github.com/mikemrm/masscan-exporter/internal/collector"
	"github.com/mikemrm/masscan-exporter/internal/exporter"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type ctxConfigKey struct{}

var configKey = ctxConfigKey{}

type config struct {
	LogLevel   zerolog.Level      `mapstructure:"loglevel"`
	Collectors []collector.Config `mapstructure:"collectors"`
	Exporter   exporter.Config    `mapstructure:"exporter"`
	Server     struct {
		Listen string `mapstructure:"listen"`
	} `mapstructure:"server"`
}

func getConfig(ctx context.Context) config {
	return ctx.Value(configKey).(config)
}

func initialize(cmd *cobra.Command, _ []string) {
	logger := zerolog.New(os.Stderr)

	v := viper.GetViper()

	v.SetEnvPrefix("MASSCAN_EXPORTER")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	v.BindPFlags(cmd.Root().Flags())
	v.BindPFlags(cmd.Root().PersistentFlags())

	v.AutomaticEnv()

	v.SetConfigFile(v.GetString("config"))

	if err := v.ReadInConfig(); err != nil {
		if !os.IsNotExist(err) {
			logger.Fatal().Err(err).Msgf("failed to load config file: %s", v.GetString("config"))
		}
	}

	var cfg config

	if err := v.Unmarshal(&cfg, decodeUnmarshalText); err != nil {
		logger.Fatal().Err(err).Msg("failed to build config")
	}

	ctx := context.WithValue(cmd.Context(), configKey, cfg)

	ctx = logger.Level(cfg.LogLevel).WithContext(ctx)

	cmd.SetContext(ctx)
}

type textUnmarshaller interface {
	UnmarshalText(text []byte) error
}

func decodeUnmarshalText(config *mapstructure.DecoderConfig) {
	hook := func(from, to reflect.Value) (any, error) {
		toI, ok := reflect.Indirect(to).Addr().Interface().(textUnmarshaller)
		if !ok {
			return from.Interface(), nil
		}

		switch fromI := from.Interface().(type) {
		case []byte:
			toI.UnmarshalText(fromI)
		case string:
			toI.UnmarshalText([]byte(fromI))
		default:
			return fromI, nil
		}

		return nil, nil
	}

	if config.DecodeHook != nil {
		config.DecodeHook = mapstructure.ComposeDecodeHookFunc(config.DecodeHook, hook)
	} else {
		config.DecodeHook = hook
	}
}

func init() {
	RootCmd.PersistentFlags().String("config", "config.yaml", "config file path")
	RootCmd.PersistentFlags().String("loglevel", "info", "set the log level")
}
