package cmd

import (
	"net/http"
	"sync"
	"time"

	"github.com/mikemrm/masscan-exporter/internal/exporter"
	"github.com/mikemrm/masscan-exporter/internal/masscan"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

func runExporter(cmd *cobra.Command, _ []string) {
	ctx := cmd.Context()

	logger := zerolog.Ctx(ctx)
	cfg := getConfig(ctx)

	masscan, err := masscan.New(ctx, masscan.WithConfig(cfg.Masscan))
	if err != nil {
		logger.Fatal().Err(err).Send()
	}

	registry := prometheus.NewRegistry()

	cfg.Exporter.Registerer = registry

	if _, err := exporter.New(ctx, masscan, exporter.WithConfig(cfg.Exporter)); err != nil {
		logger.Fatal().Err(err).Msg("failed to initialize exporter")
	}

	mux := http.NewServeMux()

	var (
		inFlightMu sync.RWMutex
		inFlight   bool
	)

	isInFlight := func() bool {
		inFlightMu.RLock()
		defer inFlightMu.RUnlock()

		return inFlight
	}

	setInFlight := func(v bool) {
		inFlightMu.Lock()
		defer inFlightMu.Unlock()

		inFlight = v
	}

	mux.Handle("/metrics", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isInFlight() {
			logger.Warn().Msg("request already in flight")

			w.WriteHeader(http.StatusTooManyRequests)

			return
		}

		setInFlight(true)
		defer setInFlight(false)

		s := time.Now()

		logger.Info().Msg("starting request")

		promhttp.HandlerFor(registry, promhttp.HandlerOpts{}).ServeHTTP(w, r)

		logger.Info().Msgf("request completed in %s", time.Since(s))
	}))

	logger.Info().Msgf("Listening on %s", cfg.Server.Listen)

	if err := http.ListenAndServe(cfg.Server.Listen, mux); err != nil {
		logger.Fatal().Err(err).Msg("error starting server")
	}
}

func init() {
	masscan.AddFlags(RootCmd.Flags())
	exporter.AddFlags(RootCmd.Flags())

	RootCmd.Flags().String("server.listen", ":9187", "listen address for the metrics server")
}
