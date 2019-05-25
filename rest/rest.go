package rest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/j30ng/ghub/profile"
)

// MakeRawAPICallWithProfile is identical to MakeRawAPICall, except it takes a profile.Profile
// as an argument instead.
func MakeRawAPICallWithProfile(profile profile.Profile, path string) (*string, error) {
	return MakeRawAPICall(profile.Endpoint, profile.Token, path)
}

// MakeRawAPICall makes an API call to the URL determined by the arguments passed to it, then
// returns the response body as a string.
func MakeRawAPICall(endpoint string, token string, path string) (*string, error) {
	apiBaseURL := strings.TrimRight(endpoint, "/")
	fullURL := apiBaseURL + "/" + strings.TrimLeft(path, "/")

	req, err := generateRequest(fullURL, token)
	if err != nil {
		return nil, err
	}

	res, err := (&http.Client{}).Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	str := string(b)
	return &str, nil
}

// MakeAPICall makes an API call via MakeRawAPICallWithProfile, tries to convert the response
// into a map[string]interface{}, then returns the converted response if the conversion was successful.
func MakeAPICall(profile profile.Profile, path string) (interface{}, error) {
	jsonString, err := MakeRawAPICallWithProfile(profile, path)
	if err != nil {
		return nil, err
	}

	var m map[string]interface{}
	reader := strings.NewReader(*jsonString)
	if err = json.NewDecoder(reader).Decode(&m); err != nil {
		return nil, err
	}
	return m, nil
}

func generateRequest(url string, token string) (*http.Request, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	return req, nil
}
