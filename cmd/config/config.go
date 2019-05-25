package config

import (
	"github.com/spf13/cobra"
)

// Cmd represents the 'config' subcommand
var Cmd = &cobra.Command{
	Use:   "config",
	Short: "Configure ghub.",
}
