package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/j30ng/ghub/profile"

	"github.com/spf13/cobra"
)

var currentProfileCmd = &cobra.Command{
	Use:   "current-profile",
	Short: "Prints out the name of the current profile.",
	Args:  currentProfileArgs,
	RunE:  currentProfileRunE,
}

func init() {
	Cmd.AddCommand(currentProfileCmd)

	currentProfileCmd.Flags().BoolP("detail", "d", false, "Print detailed information about the profile in JSON format.")
}

func currentProfileArgs(cmd *cobra.Command, args []string) error {
	if len(args) != 0 {
		return errors.New("Arguments should not be given for create-profile")
	}
	return nil
}

func currentProfileRunE(cmd *cobra.Command, args []string) error {
	selectedProfile, err := profile.SelectedProfile()
	if err != nil {
		return err
	}
	inDetail, err := cmd.Flags().GetBool("detail")
	if err != nil {
		return err
	}
	if inDetail {
		if err = printlnAsJSON(*selectedProfile); err != nil {
			return err
		}
	} else {
		printlnCurrentProfile(*selectedProfile)
	}
	return nil
}

func printlnAsJSON(in interface{}) error {
	builder := strings.Builder{}
	if err := json.NewEncoder(&builder).Encode(in); err != nil {
		return err
	}
	fmt.Println(builder.String())
	return nil
}

func printlnCurrentProfile(profile profile.Profile) {
	fmt.Println(profile.Name)
}
