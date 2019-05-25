package config

import (
	"errors"
	"fmt"

	"github.com/j30ng/ghub/profile"

	"github.com/spf13/cobra"
)

var selectProfileCmd = &cobra.Command{
	Use:   "select-profile",
	Short: "Select a ghub profile.",
	Args:  selectProfileArgs,
	RunE:  selectProfileRunE,
}

func init() {
	Cmd.AddCommand(selectProfileCmd)
}

func selectProfileArgs(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("No profile name is given for the command select-profile")
	}
	return nil
}

func selectProfileRunE(cmd *cobra.Command, args []string) error {
	givenProfileName := args[0]
	if err := profile.SetSelectedProfile(givenProfileName); err != nil {
		return errors.New("An Error occurred trying to select profile.\n\n" + err.Error())
	}
	fmt.Println("Selected profile " + givenProfileName + ".")
	return nil
}
