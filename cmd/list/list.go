package list

import (
	"github.com/spf13/cobra"
)

// Cmd represnets the 'list' subcommand.
var Cmd = &cobra.Command{
	Use:   "list",
	Short: "List some useful stuffs.",
}
