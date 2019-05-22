package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

var createProfileCmd = &cobra.Command{
	Use:   "create-profile",
	Short: "Create a ghub profile.",
	Args: func(cmd *cobra.Command, args []string) error {
		name, token, userid := cmd.Flags().Lookup("name").Value.String(), cmd.Flags().Lookup("token").Value.String(), cmd.Flags().Lookup("userid").Value.String()
		if token == "" {
			return errors.New("Personal Access Token must be given with --token.")
		}
		fmt.Println(token)
		if name != "" {
			fmt.Println(name)
		}
		if userid != "" {
			fmt.Println(userid)
		}
		return errors.New("sadf")
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("inside createProfile")
	},
}

func init() {
	configCmd.AddCommand(createProfileCmd)

	createProfileCmd.Flags().String("token", "", "Personal Access Token")
	createProfileCmd.MarkFlagRequired("token")

	createProfileCmd.Flags().StringP("api-endpoint", "e", "https://api.github.com", "The REST API Endpoint. Typically the enpoint for a GitHub Enterprise app is 'https://[github_enterprise_url]/api/v3'.")
	createProfileCmd.Flags().String("name", "", "Name of the profile")
	createProfileCmd.Flags().String("userid", "", "Login User ID")
}
