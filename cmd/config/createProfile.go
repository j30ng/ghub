package config

import (
	"fmt"
	"os"

	"github.com/j30ng/ghub/profile"
	"github.com/j30ng/ghub/rest"

	"github.com/spf13/cobra"
)

var createProfileCmd = &cobra.Command{
	Use:   "create-profile --token [token]",
	Short: "Create a ghub profile.",
	Run: run,
}

func init() {
	ConfigCmd.AddCommand(createProfileCmd)

	createProfileCmd.Flags().String("token", "", "Personal Access Token (Required)")
	createProfileCmd.MarkFlagRequired("token")

	createProfileCmd.Flags().StringP(
		"endpoint", "e", "https://api.github.com",
		"The REST API Endpoint.\nA typical enpoint for a GitHub Enterprise is 'https://[github_enterprise_url]/api/v3'.")

	createProfileCmd.Flags().String("name", "", "A name for the profile (defaults to 'profile-[loginId]'.)")
	createProfileCmd.Flags().String("userid", "", "Login user ID (defaults to the login id retrieved with the given token.)")
}

func run(cmd *cobra.Command, args []string) {
	if profiles, err := profile.GetProfiles(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		name := cmd.Flags().Lookup("name").Value.String()
		userid := cmd.Flags().Lookup("userid").Value.String()
		token := cmd.Flags().Lookup("token").Value.String()
		endpoint := cmd.Flags().Lookup("endpoint").Value.String()

		newProfile := profile.Profile{Name: name, Userid: userid, Token: token, Endpoint: endpoint}

		createProfile(newProfile, profiles)
	}
}

func createProfile(newProfile profile.Profile, existingProfiles []profile.Profile) {
	if found := profile.FindProfile(existingProfiles, profile.WithSameTokenAs(newProfile.Token)); found != nil {
		fmt.Println(fmt.Sprintf("Profile %s has the same token. You might want to use it instead.", found.Name))
		os.Exit(1)
	} else if newProfile.Name != "" && isProfileNameExists(newProfile.Name, existingProfiles) {
		fmt.Println(fmt.Sprintf("A profile with the name %s already exists. try another name.", newProfile.Name))
		os.Exit(1)
	}
	completeProfileAndSave(newProfile, existingProfiles)
}

func completeProfileAndSave(newProfile profile.Profile, existingProfiles []profile.Profile) {
	m, err := rest.Endpoint(newProfile.Endpoint).Path("/user").Token(newProfile.Token)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if newProfile.Userid == "" {
		newProfile.Userid = m["login"].(string)
	}
	if newProfile.Name == "" {
		completeProfileName(&newProfile, existingProfiles)
	}
	saveProfile(newProfile, existingProfiles)

}

func completeProfileName(newProfile *profile.Profile, existingProfiles []profile.Profile) {
	abort := make(chan struct{})
	nameChannel := generateProfileNames(abort, newProfile.Userid)
	for candidateName := range nameChannel {
		if !isProfileNameExists(candidateName, existingProfiles) {
			newProfile.Name = candidateName
			close(abort)
			break
		}
	}
}

func saveProfile(newProfile profile.Profile, existingProfiles []profile.Profile) {
	updatedProfiles := append(existingProfiles, newProfile)
	profile.SetProfiles(updatedProfiles)
	fmt.Println(fmt.Sprintf(
`Created profile %s. To use it, select it by executing:

  ghub config select-profile %s
`, newProfile.Name, newProfile.Name,
	))
}

func generateProfileNames(abort <-chan struct{}, userid string) <-chan string {
	ch := make(chan string)
	go func () {
		defer close(ch)
		nextName := fmt.Sprintf("profile-%s", userid)
		for i := 1; ; i++ {
			select {
			case ch <- nextName:
			case <-abort:
				return
			}
			nextName = fmt.Sprintf("profile-%s-%d", userid, i)
		}
	}()
	return ch
	
}

func isProfileNameExists(name string, existingProfiles []profile.Profile) bool {
	found := profile.FindProfile(existingProfiles, profile.WithSameNameAs(name))
	return found != nil
}

