package root

import (
	"github.com/j30ng/ghub/profile"
	"github.com/j30ng/ghub/rest"

	"github.com/mitchellh/mapstructure"
)

func User(profile profile.Profile) (*UserResponse, error) {
	response, err := rest.MakeAPICall(profile, "/user")
	if err != nil {
		return nil, err
	}
	var ret UserResponse
	if err = mapstructure.Decode(response, &ret); err != nil {
		return nil, err
	}
	return &ret, nil
}

type UserResponse struct {
	Login string
}
