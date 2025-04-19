package cmd

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/mikemrm/masscan-exporter/internal/exporter"
	"github.com/mikemrm/masscan-exporter/internal/masscan"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type ctxConfigKey struct{}

var configKey = ctxConfigKey{}

type config struct {
	Masscan  masscan.Config  `mapstructure:"masscan"`
	Exporter exporter.Config `mapstructure:"exporter"`
	Server   struct {
		Listen string `mapstructure:"listen"`
	} `mapstructure:"server"`
}

func getConfig(ctx context.Context) config {
	return ctx.Value(configKey).(config)
}

func initialize(cmd *cobra.Command, _ []string) {
	v := viper.GetViper()

	v.SetEnvPrefix("MASSCAN_EXPORTER")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	v.BindPFlags(cmd.Root().Flags())
	v.BindPFlags(cmd.Root().PersistentFlags())

	v.AutomaticEnv()

	v.SetConfigFile(v.GetString("config"))

	if err := v.ReadInConfig(); err != nil {
		if !os.IsNotExist(err) {
			log.Fatalf("failed to load config file: %s: %s", v.GetString("config"), err.Error())
		}
	}

	var cfg config

	if err := v.Unmarshal(&cfg); err != nil {
		log.Fatalf("failed to build config: %s", err.Error())
	}

	ctx := context.WithValue(cmd.Context(), configKey, cfg)

	cmd.SetContext(ctx)
}

func init() {
	RootCmd.PersistentFlags().String("config", "config.yaml", "config file path")
}
