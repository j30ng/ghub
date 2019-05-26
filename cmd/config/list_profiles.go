package config

import (
	"errors"

	"github.com/j30ng/ghub/profile"

	"github.com/spf13/cobra"
)

var listProfilesCmd = &cobra.Command{
	Use:   "list-profiles",
	Short: "List all available profiles.",
	Args:  listProfilesArgs,
	RunE:  listProfilesRunE,
}

func init() {
	Cmd.AddCommand(listProfilesCmd)

	listProfilesCmd.Flags().BoolP("detail", "d", false, "Print all profiles in detail in JSON format.")
}

func listProfilesArgs(cmd *cobra.Command, args []string) error {
	if len(args) != 0 {
		return errors.New("Arguments should not be given for list-profiles")
	}
	return nil
}

func listProfilesRunE(cmd *cobra.Command, args []string) error {
	profiles, err := profile.Profiles()
	if err != nil {
		return err
	}
	inDetail, err := cmd.Flags().GetBool("detail")
	if err != nil {
		return err
	}
	if inDetail {
		if err = printlnAsJSON(profiles); err != nil {
			return err
		}
		return nil
	}
	for _, p := range profiles {
		printlnCurrentProfile(p)
	}
	return nil
}
