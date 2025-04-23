package main

import (
	"github.com/mikemrm/masscan-exporter/cmd"
	"github.com/rs/zerolog/log"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		log.Fatal().Err(err).Send()
	}
}
