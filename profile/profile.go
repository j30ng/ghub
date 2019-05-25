package profile

import (
	"errors"
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

type Profile struct {
	Name     string
	Userid   string
	Token    string
	Endpoint string
}

func Find(predicate func(Profile) bool) (*Profile, error) {
	profiles, err := GetProfiles()
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

func WithSameNameAs(name string) func(Profile) bool {
	return func(p Profile) bool {
		return p.Name == name
	}
}

func WithSameTokenAs(token string) func(Profile) bool {
	return func(p Profile) bool {
		return p.Token == token
	}
}

func GetProfiles() ([]Profile, error) {
	rawProfiles := viper.Get("profiles")

	var profiles []Profile
	err := mapstructure.Decode(rawProfiles, &profiles)
	if err != nil {
		return nil, err
	}

	for i, profile := range profiles {
		if profile.Token == "" {
			return nil, errors.New(fmt.Sprintf("profile[%d].token is empty.", i+1))
		}
		if profile.Endpoint == "" {
			return nil, errors.New(fmt.Sprintf("profile[%d].endpoint is empty.", i+1))
		}
	}
	return profiles, nil
}

func GetSelectedProfile() (*Profile, error) {
	selectedProfileName, ok := viper.Get("selectedProfile").(string)
	if !ok {
		return nil, errors.New("The field selectedProfile seems to have a non-string value.")
	}
	if selectedProfile, err := Find(WithSameNameAs(selectedProfileName)); err != nil {
		return nil, err
	} else if selectedProfile == nil {
		return nil, errors.New(fmt.Sprintf("selectedProfile is %s, and there seems to be no profile by that name.", selectedProfileName))
	} else {
		return selectedProfile, nil
	}
}

func SetProfiles(profiles []Profile) error {
	viper.Set("profiles", profiles)
	return viper.WriteConfig()
}

func SetSelectedProfile(profileName string) error {
	if profilePtr, err := Find(WithSameNameAs(profileName)); err != nil {
		return err
	} else if profilePtr == nil {
		return errors.New(fmt.Sprintf("No profile by the name %s exists.", profileName))
	} else {
		viper.Set("selectedProfile", profileName)
		viper.WriteConfig()

		return nil
	}
}

func CreateProfile(newProfile Profile) (created *Profile, err error) {
	if found, err := Find(WithSameTokenAs(newProfile.Token)); err != nil {
		return nil, err
	} else if found != nil {
		return nil, errors.New(fmt.Sprintf("Profile %s has the same token. You might want to use it instead.", found.Name))
	}

	if found, err := Find(WithSameNameAs(newProfile.Name)); err != nil {
		return nil, err
	} else if found != nil {
		return nil, errors.New(fmt.Sprintf("A profile with the name %s already exists. try another name.", newProfile.Name))
	}

	existingProfiles, err := GetProfiles()
	if err != nil {
		return nil, err
	}

	updatedProfiles := append(existingProfiles, newProfile)
	if err = SetProfiles(updatedProfiles); err != nil {
		return nil, err
	}

	return &newProfile, nil
}

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

