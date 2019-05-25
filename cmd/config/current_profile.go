package config

import (
	"errors"
	"fmt"
	"encoding/json"
	"strings"

	"github.com/j30ng/ghub/profile"

	"github.com/spf13/cobra"
)

var currentProfileCmd = &cobra.Command{
	Use:   "current-profile",
	Short: "Prints out the name of the current profile.",
	Args: currentProfileArgs,
	RunE: currentProfileRunE,
}

func init() {
	Cmd.AddCommand(currentProfileCmd)

	currentProfileCmd.Flags().BoolP("verbose", "v", false, "Print all about the profile in JSON format.");
}

func currentProfileArgs(cmd *cobra.Command, args []string) error {
	if len(args) != 0 {
		return errors.New("create-profile action takes no arguments.")
	}
	return nil
}

func currentProfileRunE(cmd *cobra.Command, args []string) error {
	selectedProfile, err := profile.GetSelectedProfile()
	if err != nil {
		return err
	}
	verbose, err := cmd.Flags().GetBool("verbose")
	if err != nil {
		return err
	}
	if verbose {
		if err = printCurrentProfileVerbose(*selectedProfile); err != nil {
			return err
		}
	} else {
		printCurrentProfile(*selectedProfile)
	}
	return nil
}

func printCurrentProfileVerbose(profile profile.Profile) error {
	builder := strings.Builder{}
	if err := json.NewEncoder(&builder).Encode(profile); err != nil {
		return err
	}
	fmt.Println(builder.String())
	return nil
}

func printCurrentProfile(profile profile.Profile) {
	fmt.Println(profile.Name)
}

