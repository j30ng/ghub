package config

import (
	"errors"
	"fmt"

	"github.com/j30ng/ghub/profile"
	"github.com/j30ng/ghub/rest/root"

	"github.com/spf13/cobra"
)

var createProfileCmd = &cobra.Command{
	Use:   "create-profile --token [token]",
	Short: "Create a ghub profile.",
	RunE:   createProfileRunE,
}

func init() {
	Cmd.AddCommand(createProfileCmd)

	createProfileCmd.Flags().StringVar(&newProfile.Token, "token", "", "Personal Access Token (Required)")
	createProfileCmd.MarkFlagRequired("token")

	createProfileCmd.Flags().StringVarP(&newProfile.Endpoint, "endpoint", "e", "https://api.github.com",
		"The REST API Endpoint.\nA typical enpoint for a GitHub Enterprise is 'https://[github_enterprise_url]/api/v3'.")

	createProfileCmd.Flags().StringVar(&newProfile.Name, "name", "", "A name for the profile (defaults to 'profile-[loginId]'.)")
	createProfileCmd.Flags().StringVar(&newProfile.Userid, "userid", "", "Login user ID (defaults to the login id retrieved with the given token.)")
}

var newProfile profile.Profile

func createProfileRunE(cmd *cobra.Command, args []string) error {
	completeNameAndId()
	created, err := profile.CreateProfile(newProfile)
	if err != nil {
		return errors.New("An error occurred while creating profile.\n\n" + err.Error())
	}
	fmt.Println(fmt.Sprintf("Created profile %s. To use it, select it by executing:\n\n" +
		"  ghub config select-profile %s", created.Name, created.Name))
	return nil
}

func completeNameAndId() error {
	user, err := root.User(newProfile)
	if err != nil {
		return err
	}
	if newProfile.Userid == "" {
		newProfile.Userid = user.Login
	}

	if newProfile.Name != "" {
		return nil
	}

	if availableName, err := availableProfileName(newProfile.Userid); err != nil {
		return err
	} else {
		newProfile.Name = availableName
	}
	return nil
}

func availableProfileName(userid string) (string, error) {
	abort := make(chan struct{})
	defer close(abort)
	for candidateName := range profile.GenerateProfileNames(abort, userid) {
		if found, err := profile.Find(profile.WithSameNameAs(candidateName)); err != nil {
			return "", err
		} else if found == nil {
			return candidateName, nil
		}
	}
	return "", nil
}

