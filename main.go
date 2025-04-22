package main

import (
	"github.com/mikemrm/masscan-exporter/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		panic(err)
	}
}
