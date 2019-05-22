package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/j30ng/ghub/profile"

	"github.com/spf13/cobra"
)

var selectProfileCmd = &cobra.Command{
	Use:   "select-profile",
	Short: "Select a ghub profile.",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("profile name must be given for the action select-profile.")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		givenProfileName := args[0]
		if err := profile.SetSelectedProfile(givenProfileName); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	ConfigCmd.AddCommand(selectProfileCmd)
}
