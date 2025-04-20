package cmd

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/mikemrm/masscan-exporter/internal/exporter"
	"github.com/mikemrm/masscan-exporter/internal/masscan"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
)

func runExporter(cmd *cobra.Command, _ []string) {
	cfg := getConfig(cmd.Context())

	masscan, err := masscan.New(masscan.WithConfig(cfg.Masscan))
	if err != nil {
		panic(err)
	}

	registry := prometheus.NewRegistry()

	cfg.Exporter.Registerer = registry

	if _, err := exporter.New(masscan, exporter.WithConfig(cfg.Exporter)); err != nil {
		log.Fatalf("failed to initialize exporter: %s", err.Error())
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
			log.Println("request already in flight")

			w.WriteHeader(http.StatusTooManyRequests)

			return
		}

		setInFlight(true)
		defer setInFlight(false)

		s := time.Now()

		log.Println("starting request")

		promhttp.HandlerFor(registry, promhttp.HandlerOpts{}).ServeHTTP(w, r)

		log.Printf("request completed in %s", time.Since(s))
	}))

	log.Printf("Listening on %s", cfg.Server.Listen)

	if err := http.ListenAndServe(cfg.Server.Listen, mux); err != nil {
		log.Fatalf("error from listen and serve: %s", err.Error())
	}
}

func init() {
	masscan.AddFlags(RootCmd.Flags())
	exporter.AddFlags(RootCmd.Flags())

	RootCmd.Flags().String("server.listen", ":9187", "listen address for the metrics server")
}
