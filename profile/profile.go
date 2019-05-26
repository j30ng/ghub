package profile

import (
	"errors"
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

// Profile represents a user profile, which is stored in & read from the config file
type Profile struct {
	Name       string
	Userid     string
	Token      string
	APIBaseURL string
}

// Find tests the given predicate against each profile, and returns the first profile
// encountered that satisfies the predicate. nil is returned when no profile satisfies
// the predicate, or in case of an error.
func Find(predicate func(Profile) bool) (*Profile, error) {
	profiles, err := Profiles()
	if err != nil {
		return nil, err
	}
	for _, p := range profiles {
		if predicate(p) {
			return &p, nil
		}
	}
	return nil, nil
}

// WithSameNameAs returns a predicate function that tests if a profile has the given name.
func WithSameNameAs(name string) func(Profile) bool {
	return func(p Profile) bool {
		return p.Name == name
	}
}

// WithSameTokenAs returns a predicate function that tests if a profile has the given token.
func WithSameTokenAs(token string) func(Profile) bool {
	return func(p Profile) bool {
		return p.Token == token
	}
}

// Profiles reads in then returns all profiles from the config file. Returns an error,
// should there be a profile in an undesirable state (with an empty token or an empty base url).
func Profiles() ([]Profile, error) {
	rawProfiles := viper.Get("profiles")

	var profiles []Profile
	err := mapstructure.Decode(rawProfiles, &profiles)
	if err != nil {
		return nil, err
	}

	for i, profile := range profiles {
		if profile.Token == "" {
			return nil, fmt.Errorf("profile[%d].token is empty", i+1)
		}
		if profile.APIBaseURL == "" {
			return nil, fmt.Errorf("profile[%d].apibaseurl is empty", i+1)
		}
	}
	return profiles, nil
}

// SelectedProfile reads in then returns the selected profile from the config file.
func SelectedProfile() (*Profile, error) {
	selectedProfileName, ok := viper.Get("selectedprofile").(string)
	if !ok {
		return nil, errors.New("A non-string value found for the field selectedprofile")
	}
	if selectedProfile, err := Find(WithSameNameAs(selectedProfileName)); err != nil {
		return nil, err
	} else if selectedProfile == nil {
		return nil, fmt.Errorf("Profile not found (name: %s)", selectedProfileName)
	} else {
		return selectedProfile, nil
	}
}

// SetProfiles takes a slice of Profile and overwrites it to the config file.
func SetProfiles(profiles []Profile) error {
	viper.Set("profiles", profiles)
	return viper.WriteConfig()
}

// SetSelectedProfile writes to the config file, marking the specified profile as selected.
func SetSelectedProfile(profileName string) error {
	if profilePtr, err := Find(WithSameNameAs(profileName)); err != nil {
		return err
	} else if profilePtr == nil {
		return fmt.Errorf("No profile by the name %s exists", profileName)
	} else {
		viper.Set("selectedProfile", profileName)
		viper.WriteConfig()

		return nil
	}
}

// Create takes a Profile, and attempts to append it to the config file.
func Create(newProfile Profile) (created *Profile, err error) {
	if found, err := Find(WithSameTokenAs(newProfile.Token)); err != nil {
		return nil, err
	} else if found != nil {
		return nil, fmt.Errorf("Duplicate token (%s)", found.Name)
	}

	if found, err := Find(WithSameNameAs(newProfile.Name)); err != nil {
		return nil, err
	} else if found != nil {
		return nil, fmt.Errorf("Duplicate profile name (%s)", newProfile.Name)
	}

	existingProfiles, err := Profiles()
	if err != nil {
		return nil, err
	}

	updatedProfiles := append(existingProfiles, newProfile)
	if err = SetProfiles(updatedProfiles); err != nil {
		return nil, err
	}

	return &newProfile, nil
}

// GenerateProfileNames returns a string channel that emits profile names using the default profile-naming strategy.
func GenerateProfileNames(abort <-chan struct{}, userid string) <-chan string {
	ch := make(chan string)
	go func() {
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
