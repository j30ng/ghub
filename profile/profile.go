package profile

import (
	"errors"
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

type Profile struct {
	Name   string
	Userid string
	Token  string
	Endpoint string
}

func FindProfile(profiles []Profile, predicate func(Profile) bool) *Profile {
	for _, p := range profiles {
		if predicate(p) {
			return &p
		}
	}
	return nil
}

func WithSameNameAs(name string) (func(Profile) bool) {
	return func(p Profile) bool {
		return p.Name == name
	}
}

func WithSameTokenAs(token string) (func(Profile) bool) {
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

func SetProfiles(profiles []Profile) {
	viper.Set("profiles", profiles)
	viper.WriteConfig()
}

func SetSelectedProfile(profileName string) error {
	profiles, err := GetProfiles()
	if err != nil {
		return err
	}

	if profilePtr := FindProfile(profiles, WithSameNameAs(profileName)); profilePtr == nil {
		return errors.New(fmt.Sprintf("No profile by the name %s exists.", profileName))
	} else {
		SetSelectedProfile(profilePtr.Name)

		viper.Set("selectedProfile", profileName)
		viper.WriteConfig()

		return nil
	}
}
