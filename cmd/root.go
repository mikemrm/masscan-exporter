package cmd

import (
	"github.com/spf13/cobra"
)

var RootCmd = cobra.Command{
	Use:              "masscan-exporter",
	PersistentPreRun: initialize,
	Run:              runExporter,
}
