package root

import (
	"github.com/j30ng/ghub/profile"
	"github.com/j30ng/ghub/rest"

	"github.com/mitchellh/mapstructure"
)

// User makes an API call to the path /user.
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

// UserResponse represents the response from the path /user.
type UserResponse struct {
	Login string
}
