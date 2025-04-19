package main

import (
	"log"
	"net/http"

	"github.com/mikemrm/masscan-exporter/cmd"
	"github.com/mikemrm/masscan-exporter/internal/exporter"
	"github.com/mikemrm/masscan-exporter/internal/masscan"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		panic(err)
	}
}

func oldmain() {
	cfg := masscan.Config{
		BinPath: "../../robertdavidgraham/masscan/bin/masscan",
		Ranges: []string{
			"10.7.74.1",
		},
		Ports: []string{
			"80", "443", "22", "8443",
		},
	}

	masscan, err := masscan.New(masscan.WithConfig(cfg))
	if err != nil {
		panic(err)
	}

	registry := prometheus.NewRegistry()

	if _, err := exporter.New(masscan, exporter.WithRegisterer(registry)); err != nil {
		log.Fatalf("failed to initialize exporter: %s", err.Error())
	}

	mux := http.NewServeMux()

	mux.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))

	listenAddr := ":9187"

	log.Printf("Listening on %s", listenAddr)

	http.ListenAndServe(listenAddr, mux)
}
