package raw

import (
	"errors"
	"fmt"

	"github.com/j30ng/ghub/profile"
	"github.com/j30ng/ghub/rest"
	"github.com/spf13/cobra"
)

// Cmd represents the 'raw' subcommand.
var Cmd = &cobra.Command{
	Use:   "raw",
	Short: "Send a request to the specified path.",
	Args:  rawArgs,
	RunE:  rawRunE,
}

func rawArgs(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("Got %d arguments (expected %d)", len(args), 1)
	}
	return nil
}

func rawRunE(cmd *cobra.Command, args []string) error {
	path := args[0]
	profile, err := profile.SelectedProfile()
	if err != nil {
		return errors.New("Error loading profile.\n\n" + err.Error())
	}
	json, err := rest.MakeRawAPICallWithProfile(*profile, path)
	if err != nil {
		return errors.New("There was a problem making a rest API call.\n\n" + err.Error())
	}
	fmt.Println(*json)
	return nil
}
