package cmd

import (
	"bytes"
	"net/http"
	"time"

	"github.com/mikemrm/masscan-exporter/internal/collector"
	"github.com/mikemrm/masscan-exporter/internal/exporter"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

func runExporter(cmd *cobra.Command, _ []string) {
	ctx := cmd.Context()

	logger := zerolog.Ctx(ctx)

	collectorLogger := logger.With().Str("component", "collector").Logger()
	exporterLogger := logger.With().Str("component", "exporter").Logger()
	serverLogger := logger.With().Str("component", "server").Logger()

	cfg := getConfig(ctx)

	registry := prometheus.NewRegistry()

	cfg.Exporter.Registerer = registry

	if len(cfg.Collectors) == 0 {
		collectorLogger.Warn().Msg("no collectors configured")
	}

	for _, colCfg := range cfg.Collectors {
		collector, err := collector.NewCollector(ctx, collector.WithConfig(colCfg))
		if err != nil {
			collectorLogger.Fatal().
				Err(err).
				Msgf("failed to initialize collector: %s", colCfg.Name)
		}

		cfg.Exporter.Collectors = append(cfg.Exporter.Collectors, collector)
	}

	if _, err := exporter.New(ctx, exporter.WithConfig(cfg.Exporter)); err != nil {
		exporterLogger.Fatal().
			Err(err).
			Msg("failed to initialize exporter")
	}

	mux := http.NewServeMux()

	mux.Handle("/metrics", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s := time.Now()

		serverLogger.Info().Msg("starting request")

		promhttp.HandlerFor(registry, promhttp.HandlerOpts{}).ServeHTTP(w, r)

		serverLogger.Info().Msgf("request completed in %s", time.Since(s))
	}))

	mux.Handle("/livez", http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {}))

	mux.Handle("/readyz", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		healthy := true

		var out bytes.Buffer

		if len(cfg.Collectors) != 0 {
			out.WriteString("[+]config ok\n")
		} else {
			healthy = false

			out.WriteString("[-]config not ok\n")
		}

		if *cfg.Server.UnhealthyFailedScrapes > 0 {
			for _, collector := range cfg.Exporter.Collectors {
				if collector.FailedScrapes() < *cfg.Server.UnhealthyFailedScrapes {
					out.WriteString("[+]collector/" + collector.Name() + " ok\n")
				} else {
					healthy = false

					out.WriteString("[-]collector/" + collector.Name() + " not ok\n")
				}
			}
		}

		if !healthy {
			w.WriteHeader(http.StatusInternalServerError)
		}

		w.Write(out.Bytes())
	}))

	serverLogger.Info().Msgf("Listening on %s", cfg.Server.Listen)

	if err := http.ListenAndServe(cfg.Server.Listen, mux); err != nil {
		serverLogger.Fatal().Err(err).Msg("error starting server")
	}
}

func init() {
	RootCmd.Flags().String("server.listen", ":9187", "listen address for the metrics server")
}
