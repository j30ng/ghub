package cmd

import (
	cmd "github.com/j30ng/ghub/cmd"

	"errors"
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure ghub.",
}

func init() {
	cmd.RootCmd.AddCommand(configCmd)
	configCmd.AddCommand(createProfileCmd)
	configCmd.AddCommand(currentProfileCmd)
}

type Profile struct {
	Name   string
	Userid string
	Token  string
}

func findProfile(profiles []Profile, predicate func(Profile) bool) *Profile {
	for _, p := range profiles {
		if predicate(p) {
			return &p
		}
	}
	return nil
}

func getProfiles() ([]Profile, error) {
	rawProfiles := viper.Get("profiles")

	var profiles []Profile
	err := mapstructure.Decode(rawProfiles, &profiles)
	if err != nil {
		return nil, err
	}

	for i, profile := range profiles {
		if profile.Userid == "" {
			return nil, errors.New(fmt.Sprintf("profile[%d].userid is empty.", i+1))
		}
		if profile.Token == "" {
			return nil, errors.New(fmt.Sprintf("profile[%d].token is empty.", i+1))
		}
	}
	return profiles, nil
}
