package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var selectProfileCmd = &cobra.Command{
	Use:   "create-profile",
	Short: "Create a ghub profile.",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return errors.New("profile name must be given for action create-profile.")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		givenProfileName := args[1]

		profiles, err := getProfiles()
		if err != nil {
			fmt.Println(err)
			return
		}

		matchesGivenProfileName := func(profile Profile) bool { return profile.Name == givenProfileName }
		if profilePtr := findProfile(profiles, matchesGivenProfileName); profilePtr == nil {
			fmt.Println(fmt.Sprintf("No profile by the name %s exists.", givenProfileName))
		} else {
			viper.Set("selectedProfile", profilePtr.Name)
			viper.WriteConfig()
			fmt.Println(fmt.Sprintf("Using profile %s.", profilePtr.Name))
		}
	},
}

func init() {
	configCmd.AddCommand(selectProfileCmd)
}
